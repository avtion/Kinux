package models

import "gorm.io/gorm"

// 任务
type Mission struct {
	gorm.Model
	Name       string // 名称
	Desc       string // 描述
	Namespace  string // 命名空间(默认为default)
	Deployment uint   // k8s部署文件
	Total      uint   // 任务总分（默认值为100）
}

// 任务和检查点是一对多的关系
type MissionCheckpoints struct {
	gorm.Model
	Mission    uint `gorm:"uniqueIndex:unique_index"`
	CheckPoint uint `gorm:"uniqueIndex:unique_index"`
	Percent    uint // 该检查点占任务总分的百分比
	Priority   int  // 自定义排序
}
