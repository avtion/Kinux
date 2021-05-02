package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

func init() {
	migrateQueue = append(migrateQueue, new(Score), new(ScoresSaver))
}

// Score 检查点成绩
type Score struct {
	gorm.Model
	Account    uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Lesson     uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Exam       uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Mission    uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Checkpoint uint   `gorm:"not null;uniqueIndex:score_unique_index"`
	Container  string `gorm:"not null;uniqueIndex:score_unique_index;size:256"`
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

type ScoreSaverType = uint

// 存档类型
const (
	_                ScoreSaverType = iota
	ScoreTypeMission                // 实验存档
	ScoreTypeExam                   // 考试存档
)

var _ = []ScoreSaverType{ScoreTypeExam, ScoreTypeMission}

// ScoresSaver 成绩存档
type ScoresSaver struct {
	gorm.Model
	ScoreType    ScoreSaverType `gorm:"index:ssr"`
	RawID        uint           `gorm:"index:ssr"` // 实验或考试的原ID
	RawName      string         // 实验或者考试原名
	RawCreatedAt time.Time      // 实验或者考试的创建时间
	Data         []byte
}

// NewScoreSave 创建新的存档
func NewScoreSave(ctx context.Context, saverType ScoreSaverType, rawID uint, data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("数据为空")
	}
	var (
		rawName      string
		rawCreatedAt time.Time
	)
	switch saverType {
	case ScoreTypeExam:
		temp, err := GetExam(ctx, rawID)
		if err != nil {
			return err
		}
		rawName, rawCreatedAt = temp.Name, temp.CreatedAt
	case ScoreTypeMission:
		temp, err := GetMission(ctx, rawID)
		if err != nil {
			return err
		}
		rawName, rawCreatedAt = temp.Name, temp.CreatedAt
	}
	return GetGlobalDB().WithContext(ctx).Create(&ScoresSaver{
		Model:        gorm.Model{},
		ScoreType:    saverType,
		RawID:        rawID,
		RawName:      rawName,
		RawCreatedAt: rawCreatedAt,
		Data:         data,
	}).Error
}

// ListScoreSave 获取实验存档
func ListScoreSave(ctx context.Context, saverType ScoreSaverType, fns ...func(db *gorm.DB) *gorm.DB) (
	res []*ScoresSaver, err error) {
	err = GetGlobalDB().WithContext(ctx).Where("score_type = ?", saverType).Scopes(fns...).Find(&res).Error
	return
}
