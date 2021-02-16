package models

import (
	"context"
	"gorm.io/gorm"
)

func init() {
	migrateQueue = append(migrateQueue, new(ExamScore))
}

// 考试检查点成绩
type ExamScore struct {
	gorm.Model
	Account    uint
	Exam       uint
	Mission    uint
	Checkpoint uint
}

// 创建新的实验成绩
func NewExamScoreCallback(ac, exam, mission, checkpoint uint) func(ctx context.Context) (err error) {
	return func(ctx context.Context) (err error) {
		db := GetGlobalDB().WithContext(ctx)
		return db.Create(&ExamScore{
			Account:    ac,
			Exam:       exam,
			Mission:    mission,
			Checkpoint: checkpoint,
		}).Error
	}
}
