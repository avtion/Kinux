package models

import (
	"context"
	"errors"
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
	Container  string
}

// 创建新的实验成绩
func NewExamScoreCallback(ac, exam, mission, checkpoint uint, container string) func(ctx context.Context) (err error) {
	return func(ctx context.Context) (err error) {
		db := GetGlobalDB().WithContext(ctx)
		return db.Create(&ExamScore{
			Account:    ac,
			Exam:       exam,
			Mission:    mission,
			Checkpoint: checkpoint,
			Container:  container,
		}).Error
	}
}

/*
	实现 Score 接口
*/
var _ Score = (*ExamScore)(nil)

func (m *ExamScore) ListScores(ctx context.Context, builder *PageBuilder,
	filters ...ScoreFilter) (res []*ScoreListResult, err error) {
	// 获取成绩的数据
	db := GetGlobalDB().WithContext(ctx)
	if builder != nil {
		db = builder.build(db)
	}
	for _, filter := range filters {
		db = filter(db)
	}
	var data = make([]*ExamScore, 0)
	err = db.Find(&data).Error

	var (
		accountIDs    = make([]uint, 0, len(data))
		missionIDs    = make([]uint, 0, len(data))
		checkpointIDs = make([]uint, 0, len(data))
		examIDs       = make([]uint, 0, len(data))
	)
	for _, v := range data {
		accountIDs = append(accountIDs, v.Account)
		missionIDs = append(missionIDs, v.Mission)
		checkpointIDs = append(checkpointIDs, v.Checkpoint)
		examIDs = append(examIDs, v.Exam)
	}

	acMapper, err := GetAccountsUsernameMapper(ctx, accountIDs...)
	if err != nil {
		return
	}
	mcMapper, err := GetMissionsNameMapper(ctx, missionIDs...)
	if err != nil {
		return
	}
	cpMapper, err := GetCheckpointsNameMapper(ctx, checkpointIDs...)
	if err != nil {
		return
	}
	exMapper, err := GetExamsNameMapper(ctx, examIDs...)
	if err != nil {
		return
	}

	res = make([]*ScoreListResult, 0, len(data))
	for _, v := range data {
		acName, _ := acMapper[v.Account]
		mcName, _ := mcMapper[v.Mission]
		cpName, _ := cpMapper[v.Checkpoint]
		exName, _ := exMapper[v.Exam]
		res = append(res, &ScoreListResult{
			ID:             int(v.ID),
			Account:        acName,
			Exam:           exName,
			Mission:        mcName,
			CheckpointName: cpName,
			Container:      v.Container,
		})
	}
	return
}

func (m *ExamScore) DeleteScore(ctx context.Context, id uint) (err error) {
	if m.ID == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Unscoped().Delete(new(ExamScore), id).Error
	return
}
