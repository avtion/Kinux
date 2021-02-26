package models

import (
	"context"
	"gorm.io/gorm"
)

type ScoreType = int

const (
	_ ScoreType = iota
	MissionScoreType
	ExamScoreType
)

type ScoreFilter func(db *gorm.DB) *gorm.DB

// 获取成绩数据结果
type ScoreListResult struct {
	ID             int    `json:"id"`
	Account        string `json:"account"`
	Exam           string `json:"exam"`
	Mission        string `json:"mission"`
	CheckpointName string `json:"checkpoint_name"`
	Container      string `json:"container"`
}

// 成绩接口
type Score interface {
	ListScores(ctx context.Context, builder *PageBuilder,
		filters ...ScoreFilter) (res []*ScoreListResult, err error)
	DeleteScore(ctx context.Context, id uint) (err error)
}

// 根据 ScoreType 创建对应的成绩对象
func NewScore(t ScoreType) Score {
	switch t {
	case MissionScoreType:
		return new(MissionScore)
	case ExamScoreType:
		return new(ExamScore)
	}
	return nil
}
