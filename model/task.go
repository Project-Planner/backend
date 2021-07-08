package model

import (
	"encoding/xml"
	"github.com/google/uuid"
	"net/http"
)

type Task struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name Attribute `xml:"name"`
	Milestone struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"milestone"`
	Duedate Attribute `xml:"duedate"`
	Duetime Attribute `xml:"duetime"`
	Desc     string `xml:"desc"`
	Subtasks struct {
		Text    string `xml:",chardata"`
		Subtask []Subtask `xml:"subtask"`
	} `xml:"subtasks"`
}

// NewTask parses task from the request. Returns ErrReqFieldMissing if it could not fully be parsed,
// the result might be still useful, however. Returns other errors in case of a bad request
func NewTask(r *http.Request) (Task, error) {
	// Parse HTML form from body
	if err := r.ParseForm(); err != nil {
		return Task{}, err
	}

	var t Task
	var retErr error

	if vs, ok := r.Form["name"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		t.Name = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["duedate"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		t.Duedate = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["duetime"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		t.Duetime = Attribute{Val: vs[0]}
	}

	if vs, ok := r.Form["milestone-id"]; !ok || len(vs) != 1 {
		// don't do anything for the optional field
	} else {
		t.Milestone.ID = vs[0]
	}

	if vs, ok := r.Form["desc"]; !ok || len(vs) != 1 {
		retErr = ErrReqFieldMissing
	} else {
		if vs[0] == "" {
			t.Desc = " "
		} else {
			t.Desc = vs[0]
		}
	}

	// insert subtasks here, possibly

	id, _ := uuid.NewRandom()
	t.ID = id.String()

	return t, retErr
}

// Update all non initial fields of o in the receiver.
func (t *Task) Update(o Task) {
	if o.Name.Val != "" {
		t.Name.Val = o.Name.Val
	}

	if o.Duedate.Val != "" {
		t.Duedate.Val = o.Duedate.Val
	}

	if o.Duetime.Val != "" {
		t.Duetime.Val = o.Duetime.Val
	}

	if o.Milestone.ID != "" {
		t.Milestone.ID = o.Milestone.ID
	}

	if o.Desc != "" {
		t.Desc = o.Desc
	}
}

func (t Task) String() string {
	aXML, _ := xml.Marshal(t)
	return string(aXML)
}

func (t Task) GetID() string {
	return t.ID
}

type Subtask struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name Attribute `xml:"name"`
	Duedate Attribute `xml:"duedate"`
	Duetime Attribute `xml:"duetime"`
	Desc string `xml:"desc"`
}
