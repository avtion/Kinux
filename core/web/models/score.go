package models

import (
	"context"
	"gorm.io/gorm"
)

func init() {
	migrateQueue = append(migrateQueue, new(Score))
}

// Score 检查点成绩
type Score struct {
	gorm.Model
	Account    uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Lesson     uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Exam       uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Mission    uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Checkpoint uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Container  string `gorm:"not null;uniqueIndex:score_unique_index"`
}

// FindAllAccountFinishScoreCpIDs 获取用户所有已经完成考点的检查点ID
func FindAllAccountFinishScoreCpIDs(ctx context.Context, account, exam, mission uint, containers ...string) (
	cpIDs []uint, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Score)).Where(&Score{
		Account: account,
		Mission: mission,
		Exam:    exam,
	})
	if len(containers) > 0 {
		db = db.Where("container IN ?", containers)
	}
	err = db.Pluck("checkpoint", &cpIDs).Error
	return
}

// FindAllAccountFinishScores 获取用户所有已经完成考点的成绩点
func FindAllAccountFinishScores(ctx context.Context, account, lesson, exam, mission uint, containers ...string) (
	score []*Score, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Score)).Where(
		"account = ? AND lesson = ? AND exam = ? AND mission = ?", account, lesson, exam, mission)
	if len(containers) > 0 {
		db = db.Where("container IN ?", containers)
	}
	err = db.Find(&score).Error
	return
}
