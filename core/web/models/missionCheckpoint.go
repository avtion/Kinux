package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"strings"
)

// 结构体
type MissionCheckpoints struct {
	gorm.Model
	Mission         uint   `gorm:"uniqueIndex:mission_checkpoint"`
	CheckPoint      uint   `gorm:"uniqueIndex:mission_checkpoint"`
	Percent         uint   // 该检查点占任务总分的百分比
	Priority        int    // 自定义排序
	TargetContainer string `gorm:"uniqueIndex:mission_checkpoint"` // 目标容器
}

func (mc *MissionCheckpoints) Validate() (err error) {
	if mc.Mission == 0 {
		return errors.New("实验为空")
	}
	if mc.CheckPoint == 0 {
		return errors.New("检查点为空")
	}
	if mc.TargetContainer == "" {
		return errors.New("容器为空")
	}
	if mc.Percent == 0 {
		return errors.New("所占成绩比例为空")
	}
	return nil
}

func (mc *MissionCheckpoints) BeforeSave(tx *gorm.DB) (err error) {
	rest, err := countMissionCheckpointPercent(tx, mc.Mission)
	if err != nil {
		return
	}
	if mc.Percent > rest {
		return errors.New("所占成绩比例超过100%")
	}
	return
}

var _ callbacks.BeforeSaveInterface = (*MissionCheckpoints)(nil)

// 获取任务相关的全部检查点
func FindAllMissionCheckpoints(ctx context.Context, mission uint, containers ...string) (mcs []*MissionCheckpoints, err error) {
	db := GetGlobalDB().WithContext(ctx).Where(&MissionCheckpoints{Mission: mission})
	if len(containers) > 0 {
		db = db.Where("target_container IN ?", containers)
	}
	err = db.Find(&mcs).Error
	return
}

// 编辑任务的检查点（用于任务全部检查点更新）
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

// 获取任务相关的检查点（内部实现）
func listMissionCheckpoints(missionID uint, containers []string, builder *PageBuilder) (fn func(db *gorm.DB) *gorm.DB, err error) {
	if missionID == 0 {
		err = errors.New("实验ID为空")
		return
	}
	return func(db *gorm.DB) *gorm.DB {
		db = db.Model(new(MissionCheckpoints)).Where(&MissionCheckpoints{Mission: missionID})
		if builder != nil {
			db = db.Scopes(builder.build)
		}
		if len(containers) > 0 {
			db = db.Where("target_container IN ?", containers)
		}
		return db
	}, nil
}

// 获取任务相关的检查点
func ListMissionCheckpoints(ctx context.Context, missionID uint, containers []string, builder *PageBuilder) (
	mcs []*MissionCheckpoints, err error) {
	fn, err := listMissionCheckpoints(missionID, containers, builder)
	if err != nil {
		return
	}
	err = GetGlobalDB().WithContext(ctx).Scopes(fn).Find(&mcs).Error
	return
}

// 统计任务相关的检查点
func CountMissionCheckpoints(ctx context.Context, missionID uint, containers []string) (res int64, err error) {
	fn, err := listMissionCheckpoints(missionID, containers, nil)
	if err != nil {
		return
	}
	err = GetGlobalDB().WithContext(ctx).Scopes(fn).Count(&res).Error
	return
}

// 统计任务现有的检查点已经占用的比例
func CountMissionCheckpointPercent(ctx context.Context, missionID uint) (res uint, err error) {
	return countMissionCheckpointPercent(GetGlobalDB().WithContext(ctx), missionID)
}

// 统计任务现有的检查点已经占用的比例（内部实现）
func countMissionCheckpointPercent(db *gorm.DB, missionID uint) (res uint, err error) {
	if db == nil {
		db = GetGlobalDB().WithContext(context.Background())
	}
	var data = make([]uint, 0)
	err = db.Model(new(MissionCheckpoints)).Where(
		"mission = ?", missionID).Pluck("percent", &data).Error
	if err != nil {
		return
	}
	for _, v := range data {
		res += v
	}
	return
}

// 新增任务检查点
func AddMissionCheckpoint(ctx context.Context, mc *MissionCheckpoints) (err error) {
	if err = mc.Validate(); err != nil {
		return
	}
	return GetGlobalDB().WithContext(ctx).Create(mc).Error
}

// 更新任务检查点
func EditMissionCheckpoint(ctx context.Context, mc *MissionCheckpoints) (err error) {
	if mc.ID == 0 {
		return errors.New("ID为空")
	}
	return GetGlobalDB().WithContext(ctx).Save(mc).Error
}

// 删除任务检查点
func DeleteMissionCheckpoint(ctx context.Context, id uint) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	return GetGlobalDB().WithContext(ctx).Unscoped().Delete(new(MissionCheckpoints), id).Error
}
