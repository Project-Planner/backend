package model

import (
	"encoding/xml"
)

type User struct {
	XMLName xml.Name  `xml:"user"`
	Name    Attribute `xml:"name"`
	Items   Items     `xml:"items"`
}

type Items struct {
	XMLName   xml.Name            `xml:"items"`
	Calendars []CalendarReference `xml:"calendar"`
}

type CalendarReference struct {
	XMLName xml.Name `xml:"calendar"`
	Link    string   `xml:"href,attr"`
	Perm 	Permission `xml:"perm,attr"`
}

func NewUser(name string) User {
	return User{
		Name:  Attribute{Val: name},
		Items: Items{},
	}
}

func (user User) String() string {
	var parsed, _ = xml.MarshalIndent(user, "", "\t")
	return string(parsed)
}
