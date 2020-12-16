package models

import "gorm.io/gorm"

// 部门
type Department struct {
	gorm.Model
	Name string
}
