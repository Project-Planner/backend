package model

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Calendar struct {
	XMLName     xml.Name  `xml:"calendar"`
	Text        string    `xml:",chardata"`
	Name        Attribute `xml:"name"`
	Owner       Attribute `xml:"owner"`
	ID          Attribute `xml:"id"`
	Desc        string    `xml:"desc"`
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
			Task []Task `xml:"task"`
		} `xml:"tasks"`
	} `xml:"items"`
}

// NewCalendar parses calendar from the request. Returns ErrReqFieldMissing if it could not fully be parsed,
// the result might be still useful, however. Returns other errors in case of a bad request
func NewCalendar(r *http.Request, owner string) (Calendar, error) {
	// Parse HTML form from body
	if err := r.ParseForm(); err != nil {
		return Calendar{}, err
	}

	var c Calendar
	var retErr error

	c.Owner.Val = owner

	if vs, ok := r.Form["name"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		c.Name = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["desc"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		if vs[0] == "" {
			c.Desc = " "
		} else {
			c.Desc = vs[0]
		}
	}

	c.Permissions.View.User = []Attribute{{Val: owner}}
	c.Permissions.Edit.User = []Attribute{{Val: owner}}

	c.ID = Attribute{Val: fmt.Sprintf("%s/%s", c.Owner.Val, c.Name.Val)}

	return c, retErr
}

// Update all non initial fields of o in the receiver.
func (c *Calendar) Update(o Calendar) {
	// you may only update description, because name is part of the ID
	if o.Desc != "" {
		c.Desc = o.Desc
	}
}

func (c Calendar) String() string {
	aXML, _ := xml.Marshal(c)
	return string(aXML)
}

func (c Calendar) GetID() string {
	return c.ID.Val
}
