package models

import "gorm.io/gorm"

func init() {
	migrateQueue = append(migrateQueue, new(Deployment))
}

// K8S Deployment部署文件
type Deployment struct {
	gorm.Model
	Name string `gorm:"unique"` // 名称
	Raw  []byte // Yaml源码
}
