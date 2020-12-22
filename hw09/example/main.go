package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/kiankw/resthttp"
)

func main() {

	api := resthttp.NewApi()

	api.SetRouter(
		resthttp.GET("/countries", GetAllCountries),
		resthttp.POST("/countries/", PostCountry),
		resthttp.GET("/countries/:code", GetCountry),
		resthttp.DELETE("/countries/:code", DeleteCountry),
	)
	log.Fatal(http.ListenAndServe(":9090", api.MakeHandler()))
}

type Country struct {
	Code string
	Name string
}

var store = map[string]*Country{}

var lock = sync.RWMutex{}

func GetCountry(w http.ResponseWriter, r *http.Request, ps resthttp.Params) {
	code := ps.ByName("code")

	lock.RLock()
	var country *Country
	if store[code] != nil {
		country = &Country{}
		*country = *store[code]
	}
	lock.RUnlock()

	countryjson, _ := json.Marshal(country)

	fmt.Fprint(w, string(countryjson))
}

func GetAllCountries(w http.ResponseWriter, r *http.Request, _ resthttp.Params) {
	lock.RLock()
	countries := make([]Country, len(store))
	i := 0
	for _, country := range store {
		countries[i] = *country
		i++
	}
	lock.RUnlock()

	countriesjson, _ := json.Marshal(countries)

	fmt.Fprint(w, string(countriesjson))
}

func PostCountry(w http.ResponseWriter, r *http.Request, _ resthttp.Params) {
	country := Country{}

	s, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(s, &country)
	lock.Lock()
	store[country.Code] = &country
	lock.Unlock()

	countryjson, _ := json.Marshal(country)

	fmt.Fprint(w, string(countryjson))
}

func DeleteCountry(w http.ResponseWriter, r *http.Request, ps resthttp.Params) {
	code := ps.ByName("code")
	lock.Lock()
	delete(store, code)
	lock.Unlock()
	fmt.Fprintf(w, string(http.StatusOK))
}
