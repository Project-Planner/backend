package model

import "encoding/xml"

// THIS WILL BE REFACTORED
type Calendar struct {
	XMLName xml.Name `xml:"calendar"`
	Text    string   `xml:",chardata"`
	Name    struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"name"`
	Owner struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"owner"`
	ID struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	} `xml:"id"`
	Permissions struct {
		Text string `xml:",chardata"`
		View struct {
			Text string      `xml:",chardata"`
			User []Attribute `xml:"user"`
		} `xml:"view"`
		Edit struct {
			Text string      `xml:",chardata"`
			User []Attribute `xml:"user"`
		} `xml:"edit"`
	} `xml:"permissions"`
	Items struct {
		Text         string `xml:",chardata"`
		Appointments struct {
			Text        string        `xml:",chardata"`
			Appointment []Appointment `xml:"appointment"`
		} `xml:"appointments"`
		Milestones struct {
			Text      string      `xml:",chardata"`
			Milestone []Milestone `xml:"milestone"`
		} `xml:"milestones"`
		Tasks struct {
			Text string `xml:",chardata"`
			Task []struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
				Name struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"name"`
				Milestone struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"milestone"`
				Duedate struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"duedate"`
				Duetime struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				} `xml:"duetime"`
				Desc     string `xml:"desc"`
				Subtasks struct {
					Text    string `xml:",chardata"`
					Subtask []struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
						Name struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						} `xml:"name"`
						Duedate struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						} `xml:"duedate"`
						Duetime struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						} `xml:"duetime"`
						Desc string `xml:"desc"`
					} `xml:"subtask"`
				} `xml:"subtasks"`
			} `xml:"task"`
		} `xml:"tasks"`
	} `xml:"items"`
}

type Attribute struct {
	Text string `xml:",chardata"`
	Val  string `xml:"val,attr"`
}
