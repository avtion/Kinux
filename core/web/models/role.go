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

// 角色翻译器
func RoleTranslator(identify RoleIdentify) string {
	switch identify {
	case RoleNormalAccount:
		return "学生"
	case RoleManager:
		return "教师"
	case RoleAdmin:
		return "系统管理员"
	default:
		return "游客"
	}
}
