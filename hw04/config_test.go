package tinyini

import "testing"

func TestisSame(t *testing.T) {
	con1 := new(config)
	con2 := new(config)
	con1.info = make(map[string]map[string]string)
	con2.info = make(map[string]map[string]string)

	section := "S1"
	con1.info[section] = make(map[string]string)
	con2.info[section] = make(map[string]string)

	key := "path"
	con1.info[section][key] = "abc"
	con1.info[section][key] = "abd"

	flag := isSame(con1, con2)
	expected := false
	if flag != expected {
		t.Errorf("expected '%t' but got '%t'", expected, flag)
	}
}

