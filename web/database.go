package web

import "github.com/Project-Planner/backend/model"

// Database represents the interface for the web web to use for persistent storage.
type Database interface {
	// Here go all methods required by the web web
	// example: GetUser(id string) (model.User, error)

	//GetLogin returns the login data (username, hashed pw) for the specified user, an model.ErrNotFound if the user
	// could not be found, or an error if something went wrong internally
	GetLogin(userid string) (model.Login, error)

	// AddUser adds a user with the provided username and hashedPW to the persistence layer. It returns an error if
	// the user ist not added. (Keep in mind to add a login to the auth file and to the user file).
	AddUser(userid, hashedPW string) error

	// DeleteUser deletes the user AND LOGIN with the given ID. Returns model.ErrNotFound if user was not found
	DeleteUser(userid string) error

	// GetUser returns the user, model.ErrNotFound, or another internal server error
	GetUser(userid string) (model.User, error)

	// GetCalendar returns the calendar for the specified user and calendar name. Return model.ErrNotFound if calendar
	// not found.
	GetCalendar(calendarid string) (model.Calendar, error)

	// SetCalendar sets the given calendar to the given ID. This overrides any existing calendar or creates a new one.
	// DO NOT forget to add the calendar to the user file
	SetCalendar(calendarid string, c model.Calendar) error

	// DeleteCalendar deletes the calendar with the given ID. Returns model.ErrNotFound if calendar was not found.
	// DO NOT forget to remove the calendar from the user file
	DeleteCalendar(calendarid string) error
}
