package model

import "encoding/xml"

type Attribute struct {
	XMLName xml.Name
	Value   string `xml:"val,attr"`
}

func NewAttribute(name, value string) Attribute {
	return Attribute{
		XMLName: xml.Name{"", name},
		Value:   value,
	}
}
