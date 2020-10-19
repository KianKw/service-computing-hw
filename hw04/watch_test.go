package tinyini

import "testing"


func Testlisten (t *testing.T) {
	filename := "example.ini"
	var listener ListenFunc
	listener = loadFile
	conf1, _ := listener.listen(filename)
	excepted, _ := loadFile(filename)
	if !isSame(conf1, excepted) {
		t.Errorf("expected '%q' but got not the same config", excepted.filePath)
	}
}

