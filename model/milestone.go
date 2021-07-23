package model

import (
	"encoding/xml"
	"github.com/google/uuid"
	"net/http"
)

type Milestone struct {
	Text    string    `xml:",chardata"`
	ID      string    `xml:"id,attr"`
	Name    Attribute `xml:"name"`
	Duedate Attribute `xml:"duedate"`
	Duetime Attribute `xml:"duetime"`
	Desc    string    `xml:"desc"`
}

// NewMilestone parses milestone from the request. Returns ErrReqFieldMissing if it could not fully be parsed,
// the result might be still useful, however. Returns other errors in case of a bad request
func NewMilestone(r *http.Request) (Milestone, error) {
	// Parse HTML form from body
	if err := r.ParseForm(); err != nil {
		return Milestone{}, err
	}

	var m Milestone
	var retErr error

	if vs, ok := r.Form["name"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		m.Name = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["endDate"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		m.Duedate = Attribute{Val: transformDate(vs[0])}
	}

	if vs, ok := r.Form["endTime"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		m.Duetime = Attribute{Val: transformDate(vs[0])}
	}

	if vs, ok := r.Form["desc"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		if vs[0] == "" {
			m.Desc = " "
		} else {
			m.Desc = vs[0]
		}
	}

	id, _ := uuid.NewRandom()
	m.ID = id.String()

	return m, retErr
}

// Update all non initial fields of o in the receiver.
func (m *Milestone) Update(o Milestone) {
	if o.Name.Val != "" {
		m.Name.Val = o.Name.Val
	}

	if o.Duedate.Val != "" {
		m.Duedate.Val = o.Duedate.Val
	}

	if o.Duetime.Val != "" {
		m.Duetime.Val = o.Duetime.Val
	}

	if o.Desc != "" {
		m.Desc = o.Desc
	}
}

func (m Milestone) String() string {
	aXML, _ := xml.Marshal(m)
	return string(aXML)
}

func (m Milestone) GetID() string {
	return m.ID
}
