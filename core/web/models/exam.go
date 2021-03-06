package models

import (
	"context"
	"errors"
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

// 获取考试名字的映射
func GetExamsNameMapper(ctx context.Context, id ...uint) (res map[uint]string, err error) {
	type api struct {
		ID   uint
		Name string
	}
	if len(id) == 0 {
		return nil, errors.New("没有考试ID参数")
	}

	var data = make([]*api, 0)

	if err = GetGlobalDB().WithContext(ctx).Model(new(Exam)).Where(
		"id IN ?", id).Find(&data).Error; err != nil {
		return
	}

	res = make(map[uint]string, len(data))
	for _, v := range data {
		res[v.ID] = v.Name
	}
	return
}

// 获取考试列表
func ListExams(ctx context.Context, namespace string, builder *PageBuilder) (res []*Exam, err error) {
	err = GetGlobalDB().WithContext(ctx).Where(
		"namespace = ?", namespace).Scopes(builder.build).Find(&res).Error
	return
}

// 获取考试的实验
func ListExamMissions(ctx context.Context, examsIDs ...uint) (res []*ExamMissions, err error) {
	err = GetGlobalDB().WithContext(ctx).Where("exam IN ?", examsIDs).Find(&res).Error
	return
}
