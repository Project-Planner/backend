package model

type JWT struct {
	Text       string `xml:",chardata"`
	Authorized struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"authorized"`
	TokenID struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"token_id"`
	UserID struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"user_id"`
	Expiry struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"expiry"`
}
