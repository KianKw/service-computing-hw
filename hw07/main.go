package main

import (
	"os"

	"github.com/KianKw/cloudgo/service"
	flag "github.com/spf13/pflag"
)

const (
	PORT string = "9090"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}

	server := service.NewServer()
	server.Run(":" + port)
}
