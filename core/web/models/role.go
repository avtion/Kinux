package models

type RoleIdentify = uint

const (
	_ RoleIdentify = iota
	RoleAnonymous
	RoleNormalAccount
	RoleManager
	RoleAdmin
)

// 角色数组，用于遍历
var RoleArray = [...]RoleIdentify{RoleAnonymous, RoleNormalAccount, RoleManager, RoleAdmin}
