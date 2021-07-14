package web

import "io/ioutil"

var loaded struct {
	calendar string
	editItem string
}

func load() {
	loaded = struct {
		calendar string
		editItem string
	}{}

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
}
