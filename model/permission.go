package model

type Permission int

const (
	None Permission = iota
	Read
	Edit
	Owner
)
