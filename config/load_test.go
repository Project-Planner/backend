package config

import "testing"

func TestLoad(t *testing.T) {
	c, err := load("../config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if c.Port != 80 {
		t.Fatal("port not parsed correctly")
	}

	want := "./web/statics"
	if c.StaticDir != want {
		t.Fatal("static dir not parsed correctly, want: " + want + " got: " + c.StaticDir)
	}
}
