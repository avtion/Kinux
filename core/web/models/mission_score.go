package models

import (
	"context"
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
func FindAllAccountFinishMissionScore(ctx context.Context, account, mission uint, container string) (cpIDs []uint, err error) {
	err = GetGlobalDB().WithContext(ctx).Model(new(MissionScore)).Where(&MissionScore{
		Account:   account,
		Mission:   mission,
		Container: container,
	}).Pluck("checkpoint", &cpIDs).Error
	return
}
