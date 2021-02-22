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
