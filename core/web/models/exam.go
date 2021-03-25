package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"strings"
	"time"
)

func init() {
	migrateQueue = append(migrateQueue, new(Exam), new(ExamMissions))
}

// 考试 - 对 Mission 的封装
type Exam struct {
	gorm.Model
	Name       string        // 名称
	Desc       string        // 描述
	Total      uint          // 任务总分（默认值为100）
	BeginAt    time.Time     // 开始时间
	EndAt      time.Time     // 结束时间
	ForceOrder bool          // 强制按照顺序完成任务
	TimeLimit  time.Duration // 考试限制时间 （分钟作为单位）
	Lesson     uint          // 课程 考试与课程是N:1的关系
}

// 实验与任务点为一对多关系
type ExamMissions struct {
	gorm.Model
	Exam     uint `gorm:"uniqueIndex:ex_missions"`
	Mission  uint `gorm:"uniqueIndex:ex_missions"`
	Percent  uint // 任务占考试成绩比例
	Priority int  // 自定义排序
}

// 查询考试
func GetExam(ctx context.Context, id uint) (res *Exam, err error) {
	res = new(Exam)
	err = GetGlobalDB().WithContext(ctx).First(&res, id).Error
	return
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
func ListExams(ctx context.Context, fns ...func(db *gorm.DB) *gorm.DB) (res []*Exam, err error) {
	err = GetGlobalDB().WithContext(ctx).Scopes(fns...).Find(&res).Error
	return
}

// 获取考试列表
func CountExams(ctx context.Context, fns ...func(db *gorm.DB) *gorm.DB) (res int64, err error) {
	err = GetGlobalDB().WithContext(ctx).Model(new(Exam)).Scopes(fns...).Count(&res).Error
	return
}

// 获取考试的实验
func GetExamMissions(ctx context.Context, examsIDs ...uint) (res []*ExamMissions, err error) {
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

	if e.BeginAt.After(e.EndAt) {
		return errors.New("考试开始时间不能超过结束时间")
	}
	return GetGlobalDB().WithContext(ctx).Create(e).Error
}

// 删除考试
func DeleteExam(ctx context.Context, id uint) (err error) {
	return GetGlobalDB().WithContext(ctx).Delete(new(Exam), id).Error
}

// 获取考试已经使用的实验占比
func GetExamMissionsUsedPercent(ctx context.Context, id uint) (res uint, err error) {
	return getExamMissionsUsedPercent(GetGlobalDB().WithContext(ctx), id)
}

// 获取考试已经使用的实验占比（内部实现）
func getExamMissionsUsedPercent(db *gorm.DB, id uint) (res uint, err error) {
	var data []uint
	err = db.Model(new(ExamMissions)).Where(
		"exam = ?", id).Pluck("percent", &data).Error
	for _, v := range data {
		res += v
	}
	return
}

// 钩子函数
func (em *ExamMissions) BeforeCreate(tx *gorm.DB) (err error) {
	used, err := getExamMissionsUsedPercent(tx, em.Exam)
	if err != nil {
		return
	}
	if em.Percent > (100 - used) {
		return errors.New("所占成绩比例超过100%")
	}
	return
}

var _ callbacks.BeforeCreateInterface = (*ExamMissions)(nil)

// 创建
func (em *ExamMissions) Create(ctx context.Context) (err error) {
	if em.Exam == 0 || em.Mission == 0 {
		return errors.New("必须指定考试和实验")
	}
	if em.Percent == 0 {
		return errors.New("成绩比例不能为空")
	}
	return GetGlobalDB().WithContext(ctx).Create(em).Error
}

// 保存
func (em *ExamMissions) Save(ctx context.Context) (err error) {
	if em.Exam == 0 || em.Mission == 0 {
		return errors.New("必须指定考试和实验")
	}
	if em.Percent == 0 {
		return errors.New("成绩比例不能为空")
	}
	return GetGlobalDB().WithContext(ctx).Save(em).Error
}

// 删除
func DeleteExamMission(ctx context.Context, id uint) (err error) {
	if id == 0 {
		return errors.New("id不能为空")
	}
	return GetGlobalDB().WithContext(ctx).Delete(new(ExamMissions), id).Error
}

func listExamMission(exam, mission uint, builder *PageBuilder) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if exam != 0 {
			db = db.Where("exam = ?", exam)
		}
		if mission != 0 {
			db = db.Where("mission = ?", mission)
		}
		if builder != nil {
			db = db.Scopes(builder.Build)
		}
		return db
	}
}

// 获取考试的实验数据
func ListExamMissions(ctx context.Context, exam, mission uint, builder *PageBuilder) (
	res []*ExamMissions, err error) {
	err = GetGlobalDB().WithContext(ctx).Scopes(listExamMission(exam, mission, builder)).Find(&res).Error
	return
}

// 统计考试的实验数据
func CountExamMissions(ctx context.Context, exam, mission uint) (
	res int64, err error) {
	err = GetGlobalDB().WithContext(ctx).Scopes(listExamMission(exam, mission, nil)).Count(&res).Error
	return
}
