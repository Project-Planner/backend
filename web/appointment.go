package web

import (
	"encoding/xml"
	"errors"
	"github.com/Project-Planner/backend/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func postAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	a, err := model.NewAppointment(r)
	if err == model.ErrReqFieldMissing {
		aXML, _ := xml.Marshal(a)
		http.Error(w, "required field was missing, got:\n"+string(aXML), http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		http.Error(w, "could not parse sent data", http.StatusBadRequest)
		return
	}

	c, err := getCalendarIfPermission(w, r, model.Edit)
	if err != nil {
		// err reporting already done by method call
		return
	}

	c.Items.Appointments.Appointment = append(c.Items.Appointments.Appointment, a)

	err = db.SetCalendar(c.ID.Val, c)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(a.String()))
}

func putAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data for put
	a, err := model.NewAppointment(r)
	if err != nil && err != model.ErrReqFieldMissing {
		http.Error(w, "could not parse sent data", http.StatusBadRequest)
		return
	}

	// get calendar, must be able to edit
	c, err := getCalendarIfPermission(w, r, model.Edit)
	if err != nil {
		// err reporting already done by method call
		return
	}

	items := c.Items.Appointments.Appointment

	// find idx of item to be edited
	ids := make([]model.Identifier, len(items))
	for i, v := range items {
		ids[i] = v
	}
	idx, err := itemIdx(w, r, ids...)
	if err != nil {
		return // err reporting already done by method call
	}

	items[idx].Update(a)

	err = db.SetCalendar(c.ID.Val, c)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(items[idx].String()))
}

func deleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Edit)
	if err != nil {
		// err reporting already done by method call
		return
	}

	items := c.Items.Appointments.Appointment

	ids := make([]model.Identifier, len(items))
	for i, v := range items {
		ids[i] = v
	}
	idx, err := itemIdx(w, r, ids...)
	if err != nil {
		return // err reporting already done by method call
	}

	items[idx] = items[len(items)-1]
	c.Items.Appointments.Appointment = items[:len(items)-1]

	err = db.SetCalendar(c.ID.Val, c)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// itemIdx returns the index of the requested id in arr (by ID) and handles error reporting. Return -1 means,
// not found.
// In case of non-nil error just return in the calling function.
func itemIdx(w http.ResponseWriter, r *http.Request, arr ...model.Identifier) (int, error) {
	id, ok := mux.Vars(r)[itemIDStr]
	if !ok {
		http.Error(w, "id of item missing", http.StatusBadRequest)
		return -1, errors.New("bad request")
	}

	idx := -1
	for i, v := range arr {
		if v.GetID() == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		http.Error(w, "item with given id not found", http.StatusNotFound)
		return idx, model.ErrNotFound
	}

	return idx, nil
}
