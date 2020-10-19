package main

import (
	"github.com/KianKw/tinyini"
)

func main() {
	var listen tinyini.ListenFunc
	tinyini.Watch("example.ini", listen)
}
