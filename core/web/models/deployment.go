package models

import "gorm.io/gorm"

// K8S Deployment部署文件
type Deployment struct {
	gorm.Model
	Name string // 名称
	Raw  []byte // Yaml源码
}
