package model

type Attribute struct {
	Text string `xml:",chardata"`
	Val  string `xml:"val,attr"`
}
