package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var user UserInfo

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	//unrolled/render
	formatter := render.New(render.Options{
		Directory:  "templates",
		Extensions: []string{".html"},
		IndentJSON: true,
	})

	n := negroni.Classic()
	//gorilla/mux
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
			//fmt.Println(root)
		}
	}

	//js
	mx.HandleFunc("/json", apiTestHandler(formatter)).Methods("GET")
	//template
	mx.HandleFunc("/", homeHandler(formatter)).Methods("GET")
	//form
	mx.HandleFunc("/", submitPorcess).Methods("POST")
	//static file
	mx.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/templates/"))))

}

func homeHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.HTML(w, http.StatusOK, "index", nil)
	}
}

func submitPorcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	user.Username = r.Form["username"][0]
	user.Password = r.Form["password"][0]
	if len(user.Username) == 0 || len(user.Password) == 0 {
		log.Fatal("err")
		return
	}
	t := template.Must(template.ParseFiles("templates/table.html"))
	data := map[string]string{
		"Username": user.Username,
		"Password": user.Password,
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}
