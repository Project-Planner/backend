package model

import "encoding/xml"

// In model go all the structs that are required by both web and xmldb (or any db implementation)

type User struct {
	XMLName xml.Name `xml:"user"`
	Text    string   `xml:",chardata"`
	Name    struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"name"`
	Items struct {
		Text      string `xml:",chardata"`
		Calendars []struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
		} `xml:"calendar"`
	} `xml:"items"`
}
