
// AUTO-GENERATED CODE; DO NOT EDIT

package web

import (
	"fmt"
	"github.com/gorilla/mux"
)

func attachEndpoints(r *mux.Router) {
	appointmentsRouter := r.PathPrefix("/api/appointments").Subrouter()

	appointmentsRouter.HandleFunc(fmt.Sprintf("/post/{%s}/{%s}", userIDStr, calendarIDStr), postAppointmentHandler).Methods("POST")
	appointmentsRouter.HandleFunc(fmt.Sprintf("/post/{%s}", calendarIDStr), postAppointmentHandler).Methods("POST")
	appointmentsRouter.HandleFunc("/post", postAppointmentHandler).Methods("POST")

	appointmentsRouter.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), methodHandler(nil, putAppointmentHandler, deleteAppointmentHandler)).Methods("POST")	
	appointmentsRouter.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}", calendarIDStr, itemIDStr), methodHandler(nil, putAppointmentHandler, deleteAppointmentHandler)).Methods("POST")
	appointmentsRouter.HandleFunc(fmt.Sprintf("/other/{%s}", itemIDStr), methodHandler(nil, putAppointmentHandler, deleteAppointmentHandler)).Methods("POST")


	milestonesRouter := r.PathPrefix("/api/milestones").Subrouter()

	milestonesRouter.HandleFunc(fmt.Sprintf("/post/{%s}/{%s}", userIDStr, calendarIDStr), postMilestoneHandler).Methods("POST")
	milestonesRouter.HandleFunc(fmt.Sprintf("/post/{%s}", calendarIDStr), postMilestoneHandler).Methods("POST")
	milestonesRouter.HandleFunc("/post", postMilestoneHandler).Methods("POST")

	milestonesRouter.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), methodHandler(nil, putMilestoneHandler, deleteMilestoneHandler)).Methods("POST")	
	milestonesRouter.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}", calendarIDStr, itemIDStr), methodHandler(nil, putMilestoneHandler, deleteMilestoneHandler)).Methods("POST")
	milestonesRouter.HandleFunc(fmt.Sprintf("/other/{%s}", itemIDStr), methodHandler(nil, putMilestoneHandler, deleteMilestoneHandler)).Methods("POST")


	tasksRouter := r.PathPrefix("/api/tasks").Subrouter()

	tasksRouter.HandleFunc(fmt.Sprintf("/post/{%s}/{%s}", userIDStr, calendarIDStr), postTaskHandler).Methods("POST")
	tasksRouter.HandleFunc(fmt.Sprintf("/post/{%s}", calendarIDStr), postTaskHandler).Methods("POST")
	tasksRouter.HandleFunc("/post", postTaskHandler).Methods("POST")

	tasksRouter.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), methodHandler(nil, putTaskHandler, deleteTaskHandler)).Methods("POST")	
	tasksRouter.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}", calendarIDStr, itemIDStr), methodHandler(nil, putTaskHandler, deleteTaskHandler)).Methods("POST")
	tasksRouter.HandleFunc(fmt.Sprintf("/other/{%s}", itemIDStr), methodHandler(nil, putTaskHandler, deleteTaskHandler)).Methods("POST")


}
