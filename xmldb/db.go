package xmldb

import "github.com/Project-Planner/backend/model"

// The struct implementing the web.Database interface
type database struct {
	conf DBConfig
}

func (d database) GetLogin(userid string) (model.Login, error) {
	panic("implement me")
	return model.Login{}, model.ErrNotFound // DO THIS IF ENTITY NOT FOUND, return another error if something else went wrong
}

// Example implementation of a method
// func (db database) GetUser(id string) (model.User, error) {}
