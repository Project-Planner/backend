package web

import (
	"encoding/xml"
	"fmt"
	"github.com/Project-Planner/backend/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	v := mux.Vars(r)
	authedUser := r.Context().Value(userIDStr).(string)

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
		return
	} else if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	perm := calendarPermissions(c, authedUser)

	if perm == model.None {
		http.Error(w, "no permissions to view this calendar", http.StatusForbidden)
		return
	}

	xmlRaw, _ := xml.Marshal(c)
	xmlStr := addStylesheet(string(xmlRaw), conf.AuthedPathName+"/calendar.xsl?"+r.URL.RawQuery)

	w.Write([]byte(xmlStr))
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
