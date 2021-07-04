package web

import (
	"github.com/Project-Planner/backend/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func getCalendarHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	uID := v[userIDStr]
	cID := v[calendarIDStr]
	authedUser := r.Context().Value(userIDStr).(string)

	c := model.Calendar{} // REMOVE THIS
	/*
		c, err := db.GetCalendar(uID + "/" + cID)
		if err == model.ErrNotFound {
			http.Error(w, "", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	*/

	perm := calendarPermissions(c, authedUser)

	if perm > model.None {
		// TODO: serve calendar
	}
}
