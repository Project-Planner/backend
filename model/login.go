package model

type Login struct {
	Text string `xml:",chardata"`
	Name struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"name"`
	Hash struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"hash"`
}
