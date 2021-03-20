package models

import (
	"context"
	"gorm.io/gorm"
)

func init() {
	migrateQueue = append(migrateQueue, new(Lesson), new(LessonMission), new(LessonExam))
}

// 课程
type Lesson struct {
	gorm.Model
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// 课程实验
type LessonMission struct {
	gorm.Model
	Lesson  uint
	Mission uint
}

// 课程考试
type LessonExam struct {
	gorm.Model
	Lesson uint
	Exam   uint
}

// 课程班级
type LessonDepartment struct {
	gorm.Model
	Department uint
	Lesson     uint
}

func GetLesson(ctx context.Context, id uint) (res *Lesson, err error) {
	res = new(Lesson)
	err = GetGlobalDB().WithContext(ctx).First(&res, id).Error
	return
}

// 根据课程ID查询实验ID
func GetMissionIDsByLessons(ctx context.Context, fns ...func(db *gorm.DB) *gorm.DB) (res []uint, err error) {
	err = GetGlobalDB().WithContext(ctx).Scopes(fns...).Pluck("mission", &res).Error
	return
}
