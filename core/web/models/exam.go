package models

import (
	"gorm.io/gorm"
	"time"
)

func init() {
	migrateQueue = append(migrateQueue, new(Exam), new(ExamMissions))
}

// 考试 - 对 Mission 的封装
type Exam struct {
	gorm.Model
	Name       string    // 名称
	Desc       string    // 描述
	Namespace  string    // 命名空间
	Total      uint      // 任务总分（默认值为100）
	BeginAt    time.Time // 开始时间
	EndAt      time.Time // 结束时间
	ForceOrder bool      // 强制按照顺序完成任务
}

// 实验与任务点为一对多关系
type ExamMissions struct {
	gorm.Model
	Exam     uint `gorm:"uniqueIndex:ex_missions"`
	Mission  uint `gorm:"uniqueIndex:ex_missions"`
	Percent  uint // 任务占考试成绩比例
	Priority int  // 自定义排序
}
