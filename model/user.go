package model

import "encoding/xml"

type User struct {
	XMLName xml.Name  `xml:"user"`
	Name    Attribute `xml:"name"`
	Items   Items     `xml:"items"`
}

type Items struct {
	XMLName   xml.Name   `xml:"items"`
	Calendars []Calendar `xml:"calendar"`
}

type Calendar struct {
	XMLName xml.Name `xml:"calendar"`
	Link    string   `xml:"href,attr"`
}

func NewUser(name string) User {
	return User{
		Name:  NewAttribute("name", name),
		Items: Items{},
	}
}

func (user User) ToString() string {
	var parsed, _ = xml.MarshalIndent(user, "", "\t")
	return string(parsed)
}
