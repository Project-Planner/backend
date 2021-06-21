package model

import "encoding/xml"

type Auth struct {
	XMLName xml.Name `xml:"auth"`
	Text    string   `xml:",chardata"`
	Logins  []Login  `xml:"login"`
}
