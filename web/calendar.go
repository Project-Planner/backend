package web

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Project-Planner/backend/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func loadedXSLHandler(xsl string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendXSL(w, r, xsl)
	})
}

func getCalendarHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Read)
	if err != nil {
		return
	}

	m := r.URL.Query().Get("mode")
	var xslLink string
	switch m {
	case "edit":
		xslLink = "/editItem.xsl?"
	case "project":
		xslLink = "/projectView.xsl?"
	case "calendar":
		xslLink = "/calendar.xsl?"
	default:
		xslLink = "/calendar.xsl?"
	}

	xmlRaw, _ := xml.Marshal(c)
	xmlStr := addStylesheet(string(xmlRaw), conf.AuthedPathName+xslLink+r.URL.RawQuery)

	w.Write([]byte(xmlStr))
}

func deleteCalendarHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Owner)
	if err != nil {
		return
	}

	if c.Name.Val == c.Owner.Val {
		writeError(w, "you must not delete default calendar", http.StatusMethodNotAllowed)
		return
	}

	err = db.DeleteCalendar(c.GetID())
	if err == model.ErrNotFound {
		writeError(w, "calendar not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/html/mainPage.html", http.StatusSeeOther)
}

func putCalendarHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Edit)
	if err != nil {
		return
	}

	o, err := model.NewCalendar(r, c.Owner.Val)

	c.Update(o)

	err = db.SetCalendar(c.ID.Val, c)
	if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/html/mainPage.html", http.StatusSeeOther)
}

func postCalendarHandler(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		writeError(w, "", http.StatusUnauthorized)
		return
	}

	c, err := model.NewCalendar(r, authedUser)
	if err == model.ErrReqFieldMissing {
		aXML, _ := xml.Marshal(c)
		writeError(w, "required field was missing, got:\n"+string(aXML), http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		writeError(w, "could not parse sent data", http.StatusBadRequest)
		return
	}

	if !legalName(c.Name.Val) {
		writeError(w, "illegal name", http.StatusUnprocessableEntity)
		return
	}

	_, err = db.GetCalendar(c.GetID())
	if err != nil && err != model.ErrNotFound {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	} else if err == nil {
		writeError(w, "calendar already exists", http.StatusConflict)
		return
	}

	err = db.AddCalendar(c.Owner.Val, c.Name.Val)
	if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/html/mainPage.html", http.StatusSeeOther)
}

func getUserCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		writeError(w, "", http.StatusUnauthorized)
		return
	}

	u, err := db.GetUser(userid)
	if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	b, _ := xml.Marshal(u)
	xmlStr := addStylesheet(string(b), conf.AuthedPathName+"/showCalendars.xsl?"+r.URL.RawQuery)

	w.Write([]byte(xmlStr))
}

// sendXSL sends the given xsl with the URL query params from r. No further call is required
func sendXSL(w http.ResponseWriter, r *http.Request, xsl string) {
	vars := allFromURL(r.URL.Query())

	newXSL := varsIntoXSL(xsl, vars...)

	w.Write([]byte(newXSL))
}

func legalName(n string) bool {
	regex := "^[-_A-Za-z0-9]+$"
	b, _ := regexp.Match(regex, []byte(n))
	return b
}

//getCalendarIfPermission returns the requested calendar, after it has checked whether the minPerm are met by the
// requesting account. If err != nil is returned, then this error has already been dealt with via writeError and
// is just returned to indicate a guard statement early return.
func getCalendarIfPermission(w http.ResponseWriter, r *http.Request, minPerm model.Permission) (model.Calendar, error) {
	retErr := errors.New("error already reported")

	v := mux.Vars(r)
	authedUser, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		writeError(w, "", http.StatusUnauthorized)
		return model.Calendar{}, retErr
	}

	uID, ok := v[userIDStr]
	if !ok {
		uID = authedUser
	}
	cID, ok := v[calendarIDStr]
	if !ok {
		cID = authedUser
	}

	c, err := db.GetCalendar(uID + "/" + cID)
	if err == model.ErrNotFound {
		writeError(w, "", http.StatusNotFound)
		return model.Calendar{}, retErr
	} else if err != nil {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return model.Calendar{}, retErr
	}

	perm := model.CalendarPermissions(c, authedUser)

	if perm < minPerm {
		writeError(w, "no permissions to view/edit/create this item", http.StatusForbidden)
		return model.Calendar{}, retErr
	}

	return c, nil
}

func addStylesheet(xml, href string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/xsl" href="%s"?>`, href) + "\n" + xml
}

func varsIntoXSL(xsl string, vars ...varXLS) string {
	idx := strings.Index(xsl, `">`) + 2 // index after xsl:stylesheet tag

	newXSL := xsl[:idx] + "\n"
	for _, v := range vars {
		newXSL += v.String()
	}
	newXSL += xsl[idx:]

	return newXSL
}
