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
	i, err := model.NewAppointment(r)
	c, err := preparePostItem(w, r, i, err)
	if err != nil {
		return
	}

	c.Items.Appointments.Appointment = append(c.Items.Appointments.Appointment, i)

	finishItem(w, r, c)
}

func putAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data for put
	a, err := model.NewAppointment(r)
	c, err := preparePutItem(w, r, err)
	if err != nil {
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

	finishItem(w, r, c)
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

	finishItem(w, r, c)
}

//preparePostItem handles error reporting and just returns an error to indicate to return early.
func preparePostItem(w http.ResponseWriter, r *http.Request, a model.Identifier, err error) (model.Calendar, error) {
	if err == model.ErrReqFieldMissing {
		aXML, _ := xml.Marshal(a)
		writeError(w, "required field was missing, got:\n"+string(aXML), http.StatusUnprocessableEntity)
		return model.Calendar{}, err
	} else if err != nil {
		writeError(w, "could not parse sent data", http.StatusBadRequest)
		return model.Calendar{}, err
	}

	c, err := getCalendarIfPermission(w, r, model.Edit)
	// err reporting already done by method call
	return c, err
}

//preparePutItem handles error reporting and just returns an error to indicate to return early.
func preparePutItem(w http.ResponseWriter, r *http.Request, err error) (model.Calendar, error) {
	if err != nil && err != model.ErrReqFieldMissing {
		writeError(w, "could not parse sent data", http.StatusBadRequest)
		return model.Calendar{}, err
	}

	// get calendar, must be able to edit
	return getCalendarIfPermission(w, r, model.Edit)
}

func finishItem(w http.ResponseWriter, r *http.Request, c model.Calendar) {
	err := db.SetCalendar(c.ID.Val, c)
	if err == model.ErrNotFound {
		writeError(w, "calendar " + c.ID.Val + " does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/html/mainPage.html", http.StatusSeeOther)
}

// itemIdx returns the index of the requested id in arr (by ID) and handles error reporting. Return -1 means,
// not found.
// In case of non-nil error just return in the calling function.
func itemIdx(w http.ResponseWriter, r *http.Request, arr ...model.Identifier) (int, error) {
	id, ok := mux.Vars(r)[itemIDStr]
	if !ok {
		writeError(w, "id of item missing", http.StatusBadRequest)
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
		writeError(w, "item with given id not found", http.StatusNotFound)
		return idx, model.ErrNotFound
	}

	return idx, nil
}
