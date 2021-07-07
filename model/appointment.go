package model

import (
	"encoding/xml"
	"github.com/google/uuid"
	"net/http"
)

type Appointment struct {
	Text      string    `xml:",chardata"`
	ID        string    `xml:"id,attr"`
	Name      Attribute `xml:"name"`
	StartDate Attribute `xml:"startDate"`
	StartTime Attribute `xml:"startTime"`
	EndTime   Attribute `xml:"endTime"`
	EndDate   Attribute `xml:"endDate"`
	Desc      string    `xml:"desc"`
}

// NewAppointment parses appointment from the request. Returns ErrReqFieldMissing if it could not fully be parsed,
// the result might be still useful, however. Returns other errors in case of a bad request
func NewAppointment(r *http.Request) (Appointment, error) {
	// Parse HTML form from body
	if err := r.ParseForm(); err != nil {
		return Appointment{}, err
	}

	var a Appointment
	var retErr error

	if vs, ok := r.Form["name"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		a.Name = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["startDate"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		a.StartDate = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["endDate"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		a.EndDate = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["startTime"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		a.StartTime = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["endTime"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		a.EndTime = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["desc"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		a.Desc = vs[0]
	}

	id, _ := uuid.NewRandom()
	a.ID = id.String()

	return a, retErr
}

func (a Appointment) String() string {
	aXML, _ := xml.Marshal(a)
	return string(aXML)
}

func (a Appointment) GetID() string {
	return a.ID
}
