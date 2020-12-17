package models

type RoleIdentify = uint

const (
	RoleAnonymous RoleIdentify = iota
	RoleNormalAccount
	RoleManager
	RoleAdmin
)
