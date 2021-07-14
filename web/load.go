package web

import "io/ioutil"

type loadedXSL struct {
	calendar string
	editItem string
}

var loaded loadedXSL

func load() {
	loaded = loadedXSL{}

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
