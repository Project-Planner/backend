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
		writeError(w, "", http.StatusBadRequest)
		return
	}

	// parse share data form url encoded form
	vs, ok := r.Form["calendarName"]
	if !ok || len(vs) != 1 {
		writeError(w, "calendar name missing", http.StatusUnprocessableEntity)
		return
	}
	calendarName := vs[0]

	vs, ok = r.Form["userName"]
	if !ok || len(vs) != 1 {
		writeError(w, "calendar name missing", http.StatusUnprocessableEntity)
		return
	}
	userName := vs[0]

	vs, ok = r.Form["perm"]
	if !ok || len(vs) != 1 {
		writeError(w, "perm missing", http.StatusUnprocessableEntity)
		return
	}
	perm := vs[0]

	// get person initiating the share
	owner, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		writeError(w, "", http.StatusUnauthorized)
		return
	}

	// creating calendar id from owner and calendar to be shared
	id := fmt.Sprintf("%s/%s", owner, calendarName)

	c, err := db.GetCalendar(id)
	if err == model.ErrNotFound {
		writeError(w, "calendar "+id+" not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	// only owners can share calendars
	if c.Owner.Val != owner {
		writeError(w, "not owner of the calendar", http.StatusForbidden)
		return
	}

	user, err := db.GetUser(userName)
	if err == model.ErrNotFound {
		writeError(w, "specified user name not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	deleteFrom := func(us []model.Attribute, toDel model.Attribute) []model.Attribute {
		var idxs []int
		for i, v := range us {
			if v.Val == toDel.Val {
				idxs = append(idxs, i)
			}
		}
		if len(idxs) == 0 {
			return us
		}

		for i := len(idxs) - 1; i >= 0; i-- {
			us[idxs[i]] = us[len(us) - 1]
			us = us[:len(us) - 1]
		}

		return us
	}

	addUserReq := false

	// give user the permission to either view or edit
	userAttr := model.Attribute{Val: userName}
	if perm == "view" {
		c.Permissions.Edit.User = deleteFrom(c.Permissions.Edit.User, userAttr)
		c.Permissions.View.User = append(c.Permissions.View.User, userAttr)
	} else if perm == "edit" {
		c.Permissions.Edit.User = append(c.Permissions.Edit.User, userAttr)
	} else if perm == "none" {
		c.Permissions.Edit.User = deleteFrom(c.Permissions.Edit.User, userAttr)
		c.Permissions.View.User = deleteFrom(c.Permissions.View.User, userAttr)
		// Deletes the calendar from the user file
		idx := -1
		for i, v := range user.Items.Calendars {
			if v.Link == id {
				idx = i
				break
			}
		}
		if idx != -1 {
			user.Items.Calendars = append(user.Items.Calendars[:idx], user.Items.Calendars[idx+1:]...)
		}
		addUserReq = true
	} else {
		writeError(w, "permission not understood", http.StatusBadRequest)
		return
	}

	// check whether the user has already access to the calendar (view), and based on that add the calendar
	found := false
	for _, v := range user.Items.Calendars {
		if v.Link == id {
			found = true
			break
		}
	}
	// add the calendar if it isn't already added and if it hasn't been removed by permissions none before
	if !found && perm != "none" {
		user.Items.Calendars = append(user.Items.Calendars, model.CalendarReference{Link: id, Perm: perm})
		addUserReq = true
	}

	if addUserReq {
		if err = db.SetUser(userName, user); err != nil {
			log.Println(err)
			writeError(w, "", http.StatusInternalServerError)
			return
		}
	}

	if err = db.SetCalendar(id, c); err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/html/mainPage.html", http.StatusSeeOther)
}
