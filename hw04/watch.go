package tinyini

import "runtime"

var annotationFlag string

// Depending on the operating system, determine the symbol of the annotation
func init() {
	sysType := runtime.GOOS
	if sysType == "windows" {
		annotationFlag = ";"
	} else if sysType == "linux" {
		annotationFlag = "#"
	}
}

type Listener interface {
	listen(filename string)
}

type ListenFunc func(filename string) (*config, error)

func (f ListenFunc) listen(filename string) (*config, error) {
	return f(filename)
}

func Watch(filename string, listener ListenFunc) (*config, error) {
	listener = keepLoadListener
	return listener.listen(filename)
}