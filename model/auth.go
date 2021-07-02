package model

import "encoding/xml"

type Auth struct {
	XMLName xml.Name `xml:"auth"`
	Logins  []Login  `xml:"login"`
}

func NewAuth() Auth {
	return Auth{}
}

func (auth Auth) ToString() string {
	var parsed, _ = xml.MarshalIndent(auth, "", "\t")
	return string(parsed)
}
