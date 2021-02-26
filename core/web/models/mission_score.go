package models

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func init() {
	migrateQueue = append(migrateQueue, new(MissionScore))
}

// 实验检查点成绩
type MissionScore struct {
	gorm.Model
	Account    uint
	Mission    uint
	Checkpoint uint
	Container  string
}

// 创建新的任务成绩
func NewMissionScoreCallback(ac, mission, checkpoint uint, container string) func(ctx context.Context) (err error) {
	return func(ctx context.Context) (err error) {
		logrus.WithFields(logrus.Fields{
			"account":    ac,
			"mission":    mission,
			"checkpoint": checkpoint,
			"container":  container,
		}).Debug("任务检查点完成")
		return GetGlobalDB().WithContext(ctx).Create(&MissionScore{
			Account:    ac,
			Mission:    mission,
			Checkpoint: checkpoint,
			Container:  container,
		}).Error
	}
}

// 获取用户所有已经完成的检查点
func FindAllAccountFinishMissionScore(ctx context.Context, account, mission uint, containers ...string) (cpIDs []uint, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(MissionScore)).Where(&MissionScore{
		Account: account,
		Mission: mission,
	})
	if len(containers) > 0 {
		db = db.Where("container IN ?", containers)
	}
	err = db.Pluck("checkpoint", &cpIDs).Error
	return
}

/*
	实现 Score 接口
*/
var _ Score = (*MissionScore)(nil)

func (m *MissionScore) ListScores(ctx context.Context, builder *PageBuilder,
	filters ...ScoreFilter) (res []*ScoreListResult, err error) {
	// 获取成绩的数据
	db := GetGlobalDB().WithContext(ctx)
	if builder != nil {
		db = builder.build(db)
	}
	for _, filter := range filters {
		db = filter(db)
	}
	var data = make([]*MissionScore, 0)
	err = db.Find(&data).Error

	var (
		accountIDs    = make([]uint, 0, len(data))
		missionIDs    = make([]uint, 0, len(data))
		checkpointIDs = make([]uint, 0, len(data))
	)
	for _, v := range data {
		accountIDs = append(accountIDs, v.Account)
		missionIDs = append(missionIDs, v.Mission)
		checkpointIDs = append(checkpointIDs, v.Checkpoint)
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

	res = make([]*ScoreListResult, 0, len(data))
	for _, v := range data {
		acName, _ := acMapper[v.Account]
		mcName, _ := mcMapper[v.Mission]
		cpName, _ := cpMapper[v.Checkpoint]
		res = append(res, &ScoreListResult{
			ID:             int(v.ID),
			Account:        acName,
			Exam:           "",
			Mission:        mcName,
			CheckpointName: cpName,
			Container:      v.Container,
		})
	}
	return
}

func (m *MissionScore) DeleteScore(ctx context.Context) (err error) {
	if m.ID == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Unscoped().Delete(new(MissionScore), m.ID).Error
	return
}
