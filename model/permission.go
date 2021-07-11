package model

type permission int

const (
	VIEW permission = iota
	EDIT
	OWNER
)
