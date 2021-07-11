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

func getCalendarXSLHandler(w http.ResponseWriter, r *http.Request) {
	vars := make([]varXLS, len(r.URL.Query()))
	i := 0
	for k, v := range r.URL.Query() {
		vars[i] = varXLS{
			name:  k,
			value: v[0],
		}
		i++
	}

	xsl := varsIntoXSL(calendarXSL, vars...)

	w.Write([]byte(xsl))
}

func getCalendarHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Read)
	if err != nil {
		return
	}

	xmlRaw, _ := xml.Marshal(c)
	xmlStr := addStylesheet(string(xmlRaw), conf.AuthedPathName+"/calendar.xsl?"+r.URL.RawQuery)

	w.Write([]byte(xmlStr))
}

func deleteCalendarHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Owner)
	if err != nil {
		return
	}

	if c.Name.Val == c.Owner.Val {
		http.Error(w, "you must not delete default calendar", http.StatusMethodNotAllowed)
		return
	}

	err = db.DeleteCalendar(c.GetID())
	if err == model.ErrNotFound {
		http.Error(w, "calendar not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(o.String()))
}

func postCalendarHandler(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	c, err := model.NewCalendar(r, authedUser)
	if err == model.ErrReqFieldMissing {
		aXML, _ := xml.Marshal(c)
		http.Error(w, "required field was missing, got:\n"+string(aXML), http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		http.Error(w, "could not parse sent data", http.StatusBadRequest)
		return
	}

	if !legalName(c.Name.Val) {
		http.Error(w, "illegal name", http.StatusUnprocessableEntity)
		return
	}

	_, err = db.GetCalendar(c.GetID())
	if err != nil && err != model.ErrNotFound {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	} else if err == nil {
		http.Error(w, "calendar already exists", http.StatusConflict)
		return
	}

	err = db.AddCalendar(c.Owner.Val, c.GetID())
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(c.String()))
}

func getUserCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	u, err := db.GetUser(userid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	b, _ := xml.Marshal(u)
	w.Write(b)
}

func legalName(n string) bool {
	regex := "^[-+_A-Za-z0-9]+$"
	b, _ := regexp.Match(regex, []byte(n))
	return b
}

//getCalendarIfPermission returns the requested calendar, after it has checked whether the minPerm are met by the
// requesting account. If err != nil is returned, then this error has already been dealt with via http.Error and
// is just returned to indicate a guard statement early return.
func getCalendarIfPermission(w http.ResponseWriter, r *http.Request, minPerm model.Permission) (model.Calendar, error) {
	retErr := errors.New("error already reported")

	v := mux.Vars(r)
	authedUser, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
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
		http.Error(w, "", http.StatusNotFound)
		return model.Calendar{}, retErr
	} else if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println(err)
		return model.Calendar{}, retErr
	}

	perm := calendarPermissions(c, authedUser)

	if perm < minPerm {
		http.Error(w, "no permissions to view/edit/create this item", http.StatusForbidden)
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