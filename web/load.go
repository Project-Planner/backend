package web

import "io/ioutil"

var calendarXSL string

func load() {
	c, err := ioutil.ReadFile(conf.FrontendDir + "/data/calendar.xsl")
	if err != nil {
		panic(err)
	}
	calendarXSL = string(c)


}
