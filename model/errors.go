package model

import "errors"

var (
	// ErrNotFound should be returned if an entity that was searched for could not be found
	ErrNotFound = errors.New("error: entity not found")
	// ErrReqFieldMissing should be returned when parsing a struct from a web request and if important fields are missing
	ErrReqFieldMissing = errors.New("error: required field for entity parsing is missing")
	// ErrAlreadyExists should be returned when an item already exists.
	ErrAlreadyExists = errors.New("error: item already exists")
)
