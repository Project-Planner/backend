package model

import "errors"

var (
	// ErrNotFound should be returned if an entity that was searched for could not be found
	ErrNotFound = errors.New("error: entity not found")
)
