package model

type database interface {
	SetCalendar(calID string, cal Calendar) error
	SetUser(userID string, user User) error
	GetCalendar(calID string) (Calendar, error)
}
