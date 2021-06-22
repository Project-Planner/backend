package web

import "github.com/Project-Planner/backend/model"

// Database represents the interface for the web web to use for persistent storage.
type Database interface {
	// Here go all methods required by the web web
	// example: GetUser(id string) (model.User, error)

	//GetLogin returns the login data (username, hashed pw) for the specified user, an model.ErrNotFound if the user
	// could not be found, or an error if something went wrong internally
	GetLogin(userid string) (model.Login, error)
}
