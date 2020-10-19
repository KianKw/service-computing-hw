package tinyini

import (
	"bufio"
	"testing"
)

func TestPathExists(t *testing.T) {
	flag := PathExists("example.ini")
	expected := true
	if flag != expected {
		t.Errorf("expected '%t' but got '%t'", expected, flag)
	}
}

func TestreadFile(t *testing.T) {
	file, _ := readFile("example.ini")


	buf := bufio.NewReader(file)
	str, _ := buf.ReadString('\n')
	str, _ = buf.ReadString('\n')
	str, _ = buf.ReadString('\n')
	str, _ = buf.ReadString('\n')
	expected := "[paths]"
	if str != expected {
		t.Errorf("expected '%q' but got '%q'", expected, str)
	}
}

func TestloadFile(t *testing.T) {
	conf, err := loadFile("example.ini")
	if err != nil {
		return
	}
	pathname := conf.filePath
	expected := "example.ini"
	if expected != pathname {
		t.Errorf("expected '%q' but got '%q'", expected, pathname)
	}
}
