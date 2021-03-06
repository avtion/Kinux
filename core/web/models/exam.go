package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

func init() {
	migrateQueue = append(migrateQueue, new(Exam), new(ExamMissions))
}

const defaultExamNamespace = "default"

// 考试 - 对 Mission 的封装
type Exam struct {
	gorm.Model
	Name       string        // 名称
	Desc       string        // 描述
	Namespace  string        // 命名空间
	Total      uint          // 任务总分（默认值为100）
	BeginAt    time.Time     // 开始时间
	EndAt      time.Time     // 结束时间
	ForceOrder bool          // 强制按照顺序完成任务
	TimeLimit  time.Duration // 考试限制时间 （分钟作为单位）
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

func listExams(namespace string, builder *PageBuilder) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if namespace != "" {
			db = db.Where("namespace = ?", namespace)
		}
		if builder != nil {
			db = db.Scopes(builder.build)
		}
		return db
	}
}

// 获取考试列表
func ListExams(ctx context.Context, namespace string, builder *PageBuilder) (res []*Exam, err error) {
	err = GetGlobalDB().WithContext(ctx).Scopes(listExams(namespace, builder)).Find(&res).Error
	return
}

// 获取考试列表
func CountExams(ctx context.Context, namespace string) (res int64, err error) {
	err = GetGlobalDB().WithContext(ctx).Scopes(listExams(namespace, nil)).Count(&res).Error
	return
}

// 获取考试的实验
func ListExamMissions(ctx context.Context, examsIDs ...uint) (res []*ExamMissions, err error) {
	err = GetGlobalDB().WithContext(ctx).Where("exam IN ?", examsIDs).Find(&res).Error
	return
}

// 保存考试
func (e *Exam) Save(ctx context.Context) (err error) {
	if e.ID == 0 {
		return errors.New("ID为空")
	}
	return GetGlobalDB().WithContext(ctx).Save(e).Error
}

// 创建考试
func (e *Exam) Create(ctx context.Context) (err error) {
	if e.Name = strings.TrimSpace(e.Name); e.Name == "" {
		return errors.New("考试名为空")
	}

	if e.Namespace = strings.TrimSpace(e.Namespace); e.Namespace == "" {
		e.Namespace = defaultExamNamespace
	}
	if e.BeginAt.After(e.EndAt) {
		return errors.New("考试开始时间不能超过结束时间")
	}
	return GetGlobalDB().WithContext(ctx).Create(e).Error
}

// 删除考试
func DeleteExam(ctx context.Context, id uint) (err error) {
	return GetGlobalDB().WithContext(ctx).Delete(new(Exam), id).Error
}

// 获取已有考试的命名空间列表
func GetExamsNamespaces(ctx context.Context) (res []string, err error) {
	err = GetGlobalDB().WithContext(ctx).Model(new(Exam)).Distinct().Pluck("namespace", &res).Error
	return
}
