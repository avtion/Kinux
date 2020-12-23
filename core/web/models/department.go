package models

import "gorm.io/gorm"

func init() {
	migrateQueue = append(migrateQueue, new(Department))
}

// 班级
type Department struct {
	gorm.Model
	Name      string `gorm:"unique"` // 名称
	Namespace string // 可见的命名空间，以分号为间隔
}
