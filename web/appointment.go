package web

import (
	"encoding/xml"
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

}

func deleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Edit)
	if err != nil {
		// err reporting already done by method call
		return
	}

	id, ok := mux.Vars(r)[itemIDStr]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	as := c.Items.Appointments.Appointment

	idx := -1
	for i, v := range as {
		if v.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	as[idx] = as[len(as)-1]
	c.Items.Appointments.Appointment = as[:len(as)-1]

	err = db.SetCalendar(c.ID.Val, c)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
