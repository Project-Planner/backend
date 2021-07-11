package model

import "encoding/xml"

type Login struct {
	XMLName xml.Name  `xml:"login"`
	Name    Attribute `xml:"name"`
	Hash    Attribute `xml:"hash"`
}

func NewLogin(name, hash string) Login {
	return Login{
		Name: NewAttribute("name", name),
		Hash: NewAttribute("hash", hash),
	}
}

func (l Login) String() string {
	var parsed, _ = xml.MarshalIndent(l, "", "\t")
	return string(parsed)
}
