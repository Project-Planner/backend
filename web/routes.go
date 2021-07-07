package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var db Database
var conf ServerConfig
var calendarXSL string

// ListenAndServe starts the webserver with the given database implementation and config file
func ListenAndServe(database Database, configuration ServerConfig) {
	db = database
	conf = configuration
	c, err := ioutil.ReadFile(conf.HTMLDir + "/data/calendar.xsl")
	if err != nil {
		panic(err)
	}
	calendarXSL = string(c)

	// create a new router to attach routes to. Redirect to proper routes without trailing slash
	r := mux.NewRouter().StrictSlash(true)

	registerRoutes(r)

	// start the web web with the port from the config and the router
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Port), r))
}

// registerRoutes attaches request handlers and middleware to routes via the router
func registerRoutes(r *mux.Router) {
	// authed router that enforces authentication
	authed := r.PathPrefix(conf.AuthedPathName).Subrouter()

	// attach middleware for all routes
	authed.Use(auth)

	// Example of registering a function to the route "domain.tld/me/calendars", if conf.AuthedPathName = "/me":
	// authed.HandleFunc("/calendars", calandarsHandler)
	authed.HandleFunc(fmt.Sprintf("/c/{%s}/{%s}", userIDStr, calendarIDStr), getCalendarHandler).Methods("GET")
	authed.HandleFunc(fmt.Sprintf("/c/{%s}", calendarIDStr), getCalendarHandler).Methods("GET")
	authed.HandleFunc("/c", getCalendarHandler).Methods("GET")
	authed.HandleFunc("/calendar.xsl", getCalendarXSLHandler)

	// attach auto generated endpoint routes
	attachEndpoints(authed)

	authed.HandleFunc("/logout", logoutHandler).Methods("GET")

	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/register", registerHandler).Methods("POST")

	// serve static files (index, impressum, login, register ...). Note that this has to be registered last.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.StaticDir)))
}
