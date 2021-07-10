package web

import (
	"fmt"
	"github.com/Project-Planner/backend/model"
	"log"
	"net/http"
)

// sharingHandler handles the request to share a calendar with another user
func sharingHandler(w http.ResponseWriter, r *http.Request) {
	// Parse HTML form from body
	if err := r.ParseForm(); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// parse share data form url encoded form
	vs, ok := r.Form["calendarName"]
	if !ok || len(vs) != 1 {
		http.Error(w, "calendar name missing", http.StatusUnprocessableEntity)
		return
	}
	calendarName := vs[0]

	vs, ok = r.Form["userName"]
	if !ok || len(vs) != 1 {
		http.Error(w, "calendar name missing", http.StatusUnprocessableEntity)
		return
	}
	userName := vs[0]

	vs, ok = r.Form["perm"]
	if !ok || len(vs) != 1 {
		http.Error(w, "perm missing", http.StatusUnprocessableEntity)
		return
	}
	perm := vs[0]

	// get person initiating the share
	owner, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	// creating calendar id from owner and calendar to be shared
	id := fmt.Sprintf("%s/%s", owner, calendarName)

	c, err := db.GetCalendar(id)
	if err == model.ErrNotFound {
		http.Error(w, "calendar " + id + " not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// only owners can share calendars
	if c.Owner.Val != owner {
		http.Error(w, "not owner of the calendar", http.StatusForbidden)
		return
	}

	user, err := db.GetUser(userName)
	if err == model.ErrNotFound {
		http.Error(w, "specified user name not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// give user the permission to either view or edit
	if perm == "view" {
		c.Permissions.View.User = append(c.Permissions.View.User, model.Attribute{Val: userName})
	} else if perm == "edit" {
		c.Permissions.Edit.User = append(c.Permissions.Edit.User, model.Attribute{Val: userName})
	} else {
		http.Error(w, "permission not understood", http.StatusBadRequest)
		return
	}

	// check whether the user has already access to the calendar (view), and based on that add the calendar
	found := false
	for _, v := range user.Items.Calendars {
		if v.Href == id {
			found = true
		}
	}
	if !found {
		user.Items.Calendars = append(user.Items.Calendars, struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
		}{Href: id})

		if err = db.SetUser(userName, user); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	if err = db.SetCalendar(id, c); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
