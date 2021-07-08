package xmldb

import "github.com/Project-Planner/backend/model"

// The struct implementing the web.Database interface
type database struct {
	conf DBConfig
}

func (d database) DeleteUser(userid string) error {
	panic("implement me")
}

func (d database) GetUser(userid string) (model.User, error) {
	panic("implement me")
}

func (d database) DeleteCalendar(calendarid string) error {
	panic("implement me")
}

func (d database) SetCalendar(calendarid string, c model.Calendar) error {
	panic("implement me")
}

func (d database) GetCalendar(calendarid string) (model.Calendar, error) {
	panic("implement me")
}

func (d database) AddUser(userid, hashedPW string) error {
	panic("implement me")
}

func (d database) GetLogin(userid string) (model.Login, error) {
	panic("implement me")
	return model.Login{}, model.ErrNotFound // DO THIS IF ENTITY NOT FOUND, return another error if something else went wrong
}

// Example implementation of a method
// func (db database) GetUser(id string) (model.User, error) {}
