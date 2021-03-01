package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
)

// 结构体
type MissionCheckpoints struct {
	gorm.Model
	Mission         uint   `gorm:"uniqueIndex:mission_checkpoint"`
	CheckPoint      uint   `gorm:"uniqueIndex:mission_checkpoint"`
	Percent         uint   // 该检查点占任务总分的百分比
	Priority        int    // 自定义排序
	TargetContainer string // 目标容器
}

// 获取任务相关的全部检查点
func FindAllMissionCheckpoints(ctx context.Context, mission uint, containers ...string) (mcs []*MissionCheckpoints, err error) {
	db := GetGlobalDB().WithContext(ctx).Where(&MissionCheckpoints{Mission: mission})
	if len(containers) > 0 {
		db = db.Where("target_container IN ?", containers)
	}
	err = db.Find(&mcs).Error
	return
}

// 编辑任务的检查点
func EditMissionCheckpoints(ctx context.Context, missionID uint, checkpoints ...struct {
	CheckpointID    uint
	Percent         uint
	Priority        int
	TargetContainer string
}) (err error) {
	db := GetGlobalDB().WithContext(ctx)

	// 查询对应的任务
	var mission = new(Mission)
	if err = db.First(mission, missionID).Error; err != nil {
		return
	}

	// 分数余量
	var (
		allowance      int = 100
		waitAllocation     = 0
	)

	// 检查任务点是否存在
	cpIDs := make([]uint, 0, len(checkpoints))
	for _, cp := range checkpoints {
		cpIDs = append(cpIDs, cp.CheckpointID)

		// 顺便计算分数比例
		allowance -= int(cp.Percent)
		if cp.Percent == 0 {
			waitAllocation += 1
		}

		if !strings.EqualFold(cp.TargetContainer, "") && !mission.IsContainerAllowed(cp.TargetContainer) {
			return errors.New("非法容器" + cp.TargetContainer)
		}
	}
	if allowance < 0 {
		return errors.New("任务点占比总和不能超过100")
	}
	var cpCounter int64
	if err = db.Model(&Checkpoint{}).Where("id IN ?", cpIDs).Count(&cpCounter).Error; err != nil {
		return
	}
	if cpCounter != int64(len(checkpoints)) {
		err = errors.New("任务点不存在，请刷新后尝试")
		return
	}

	// 检查点分数计算
	var average int
	if waitAllocation != 0 {
		average = allowance / waitAllocation
	}

	// 创建队列
	insertQueue := make([]*MissionCheckpoints, 0, len(checkpoints))
	for _, cp := range checkpoints {
		// 计算检查点分数比例
		var p = cp.Percent
		if p == 0 {
			if waitAllocation == 1 {
				p = uint(allowance)
			} else {
				p = uint(average)
				allowance -= average
			}
			waitAllocation--
		}

		// 追加创建队列
		insertQueue = append(insertQueue, &MissionCheckpoints{
			Mission:         missionID,
			CheckPoint:      cp.CheckpointID,
			Percent:         p,
			Priority:        cp.Priority,
			TargetContainer: cp.TargetContainer,
		})
	}

	// 事务
	err = db.Transaction(func(tx *gorm.DB) error {
		// 先删除原本所有的任务的检查点
		if _err := tx.Where(&MissionCheckpoints{Mission: missionID}).Delete(&MissionCheckpoints{}).Error; _err != nil {
			return _err
		}

		// 批量插入
		if _err := tx.Create(&insertQueue).Error; _err != nil {
			return _err
		}
		return nil
	})
	return
}

// 获取用户需要完成的检查点
func FindAllTodoCheckpoints(ctx context.Context, account, mission uint, containers ...string) (cps []*Checkpoint, err error) {
	if account == 0 || mission == 0 {
		return nil, errors.New("缺乏参数，无法获取用户需要完成的检查点")
	}
	// 首先获取全部的检查点
	mcs, err := FindAllMissionCheckpoints(ctx, mission, containers...)
	if err != nil {
		return
	}

	// 查找已经完成的检查点
	finishedCheckpointsIDs, err := FindAllAccountFinishMissionScore(ctx, account, mission, containers...)
	if err != nil {
		return
	}

	// 过滤已经完成的检查点
	var idMapper = make(map[uint]struct{}, len(finishedCheckpointsIDs))
	for _, id := range finishedCheckpointsIDs {
		idMapper[id] = struct{}{}
	}
	var todoCheckpointIDs = make([]uint, 0)
	for _, mc := range mcs {
		if _, isExist := idMapper[mc.CheckPoint]; !isExist {
			todoCheckpointIDs = append(todoCheckpointIDs, mc.CheckPoint)
		}
	}

	return FindCheckpoints(ctx, todoCheckpointIDs...)
}
