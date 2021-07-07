// AUTO-GENERATED CODE; DO NOT EDIT

package web

import (
	"fmt"
	"github.com/gorilla/mux"
)

func attachEndpoints(r *mux.Router) {
	appointmentsRouter := r.PathPrefix("/api/appointments").Subrouter()

	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", userIDStr, calendarIDStr), postAppointmentHandler).Methods("POST")
	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}", calendarIDStr), postAppointmentHandler).Methods("POST")
	appointmentsRouter.HandleFunc("/", postAppointmentHandler).Methods("POST")

	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), putAppointmentHandler).Methods("PUT")
	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", calendarIDStr, itemIDStr), putAppointmentHandler).Methods("PUT")
	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}", itemIDStr), putAppointmentHandler).Methods("PUT")

	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), deleteAppointmentHandler).Methods("DELETE")
	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", calendarIDStr, itemIDStr), deleteAppointmentHandler).Methods("DELETE")
	appointmentsRouter.HandleFunc(fmt.Sprintf("/{%s}", itemIDStr), deleteAppointmentHandler).Methods("DELETE")

	milestonesRouter := r.PathPrefix("/api/milestones").Subrouter()

	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", userIDStr, calendarIDStr), postMilestoneHandler).Methods("POST")
	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}", calendarIDStr), postMilestoneHandler).Methods("POST")
	milestonesRouter.HandleFunc("/", postMilestoneHandler).Methods("POST")

	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), putMilestoneHandler).Methods("PUT")
	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", calendarIDStr, itemIDStr), putMilestoneHandler).Methods("PUT")
	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}", itemIDStr), putMilestoneHandler).Methods("PUT")

	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), deleteMilestoneHandler).Methods("DELETE")
	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", calendarIDStr, itemIDStr), deleteMilestoneHandler).Methods("DELETE")
	milestonesRouter.HandleFunc(fmt.Sprintf("/{%s}", itemIDStr), deleteMilestoneHandler).Methods("DELETE")

	tasksRouter := r.PathPrefix("/api/tasks").Subrouter()

	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", userIDStr, calendarIDStr), postTaskHandler).Methods("POST")
	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}", calendarIDStr), postTaskHandler).Methods("POST")
	tasksRouter.HandleFunc("/", postTaskHandler).Methods("POST")

	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), putTaskHandler).Methods("PUT")
	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", calendarIDStr, itemIDStr), putTaskHandler).Methods("PUT")
	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}", itemIDStr), putTaskHandler).Methods("PUT")

	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), deleteTaskHandler).Methods("DELETE")
	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}/{%s}", calendarIDStr, itemIDStr), deleteTaskHandler).Methods("DELETE")
	tasksRouter.HandleFunc(fmt.Sprintf("/{%s}", itemIDStr), deleteTaskHandler).Methods("DELETE")

}
