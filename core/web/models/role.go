package models

type RoleIdentify = uint

const (
	_ RoleIdentify = iota
	RoleAnonymous
	RoleNormalAccount
	RoleManager
	RoleAdmin
)
