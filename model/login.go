package model

import "encoding/xml"

type Login struct {
	XMLName xml.Name  `xml:"login"`
	Name    Attribute `xml:"name"`
	Hash    Attribute `xml:"hash"`
}

func NewLogin(name, hash string) Login {
	return Login{
		Name: Attribute{Val: name},
		Hash: Attribute{Val: hash},
	}
}

func (l Login) String() string {
	var parsed, _ = xml.MarshalIndent(l, "", "\t")
	return string(parsed)
}
