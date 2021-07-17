package web

import (
	"html/template"
	"io/ioutil"
)

var loaded struct {
	calendar string
	editItem string
}

func load() {
	loaded = struct {
		calendar string
		editItem string
	}{}

	// loading the xsl into memory as they are queried very often
	c, err := ioutil.ReadFile(conf.FrontendDir + "/data/calendar.xsl")
	if err != nil {
		panic(err)
	}
	loaded.calendar = string(c)

	c, err = ioutil.ReadFile(conf.FrontendDir + "/data/editItem.xsl")
	if err != nil {
		panic(err)
	}
	loaded.editItem = string(c)

	// Loading the template for fancy error reporting
	tmpl, err := ioutil.ReadFile(conf.FrontendDir + "/html/error.html")
	if err != nil {
		panic(err)
	}
	errTemplate = template.Must(template.New("error").Parse(string(tmpl)))
}
