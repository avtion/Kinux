package models

import (
	"context"
	"gorm.io/gorm"
)

func init() {
	migrateQueue = append(migrateQueue,
		new(Lesson),
		new(LessonMission),
		new(LessonExam),
		new(LessonDepartment),
	)
}

// 课程
type Lesson struct {
	gorm.Model
	Name string
	Desc string
}

// 课程实验
type LessonMission struct {
	gorm.Model
	Lesson   uint `gorm:"uniqueIndex:lesson_mission_unique_index;not null"`
	Mission  uint `gorm:"uniqueIndex:lesson_mission_unique_index;not null"`
	Priority uint `gorm:"not null"`
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
	Department uint `gorm:"uniqueIndex:department_lesson_unique_index;not null"`
	Lesson     uint `gorm:"uniqueIndex:department_lesson_unique_index;not null"`
}

func GetLesson(ctx context.Context, id uint) (res *Lesson, err error) {
	res = new(Lesson)
	err = GetGlobalDB().WithContext(ctx).First(&res, id).Error
	return
}

// 根据课程ID查询实验ID
func GetMissionIDsByLessons(ctx context.Context, fns ...func(db *gorm.DB) *gorm.DB) (res []uint, err error) {
	err = GetGlobalDB().WithContext(ctx).Model(new(LessonMission)).Scopes(fns...).Pluck("mission", &res).Error
	return
}

// 按照优先级排序
func ScopeLessonMissionOrderByOrder(isAsc ...bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(isAsc) > 0 && isAsc[0] == true {
			return db.Order("priority asc")
		}
		return db.Order("priority desc")
	}
}
