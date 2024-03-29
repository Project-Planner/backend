package web

import (
	"fmt"
	"github.com/Project-Planner/backend/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var db model.Database
var conf ServerConfig

// ListenAndServe starts the webserver with the given database implementation and config file
func ListenAndServe(database model.Database, configuration ServerConfig) {
	db = database
	conf = configuration

	// loads templates
	load()

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

	//Get Calendar
	authed.HandleFunc(fmt.Sprintf("/c/{%s}/{%s}", userIDStr, calendarIDStr), getCalendarHandler).Methods("GET")
	authed.HandleFunc(fmt.Sprintf("/c/{%s}", calendarIDStr), getCalendarHandler).Methods("GET")
	authed.HandleFunc("/c", getCalendarHandler).Methods("GET")
	authed.Handle("/calendar.xsl", loadedXSLHandler(loaded.calendar)).Methods("GET")
	authed.Handle("/projectView.xsl", loadedXSLHandler(loaded.project)).Methods("GET")
	authed.Handle("/editItem.xsl", loadedXSLHandler(loaded.editItem)).Methods("GET")

	//Get all Calendars of User
	authed.HandleFunc("/calendars", getUserCalendarsHandler).Methods("GET")
	authed.Handle("/showCalendars.xsl", loadedXSLHandler(loaded.showCalendars)).Methods("GET")

	// Modify Calendar
	authed.HandleFunc("/c", methodHandler(postCalendarHandler, putCalendarHandler, nil)).Methods("POST")
	authed.HandleFunc(fmt.Sprintf("/c/{%s}/{%s}", userIDStr, calendarIDStr),
		methodHandler(nil, putCalendarHandler, deleteCalendarHandler)).Methods("POST")
	authed.HandleFunc(fmt.Sprintf("/c/{%s}", calendarIDStr),
		methodHandler(nil, putCalendarHandler, deleteCalendarHandler)).Methods("POST")

	// Delete User
	authed.HandleFunc("/api/user", methodHandler(nil, nil, deleteUserHandler)).Methods("POST")

	authed.HandleFunc("/api/sharing", sharingHandler).Methods("POST")

	// attach auto generated endpoint routes
	attachEndpoints(authed)

	authed.HandleFunc("/logout", logoutHandler).Methods("GET")

	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/register", registerHandler).Methods("POST")

	// serve static files (index, impressum, login, register ...). Note that this has to be registered last.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.FrontendDir)))
}

func methodHandler(post, put, delete http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse HTML form from body
		if err := r.ParseForm(); err != nil {
			writeError(w, "couldn't parse form", http.StatusBadRequest)
			return
		}

		vs, ok := r.Form["_method"]
		if !ok || len(vs) != 1 {
			post(w, r) // default
			return
		}

		method := strings.ToUpper(vs[0])
		switch method {
		case "PUT":
			put(w,r)
			return
		case "DELETE":
			delete(w,r)
			return
		case "POST":
			fallthrough
		default:
			post(w,r)
			return
		}
	})
}
