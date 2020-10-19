package tinyini

import "testing"

func TestkeepLoadListener(t *testing.T) {
	filename := "example.ini"
	con1 := new(config)
	con2 := new(config)
	con1.info = make(map[string]map[string]string)
	con2.info = make(map[string]map[string]string)

	section := "paths"
	con1.info[section] = make(map[string]string)
	con2.info[section] = make(map[string]string)

	key := "data"
	con1.info[section][key] = "/src/main.go"
	con1.info[section][key] = filename

	flag := isSame(con1, con2)
	expected := false
	if flag != expected {
		t.Errorf("expected '%t' but got '%t'", expected, flag)
	}
}
