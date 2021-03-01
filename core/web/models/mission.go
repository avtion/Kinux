package models

import (
	"Kinux/core/k8s"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	v1 "k8s.io/api/core/v1"
	"strings"
)

func init() {
	migrateQueue = append(migrateQueue, new(Mission), new(MissionCheckpoints))
}

var (
	// Deployment不存在
	ErrMissionDeploymentNotExist = errors.New("mission deployment not exist")
)

// 任务
type Mission struct {
	gorm.Model

	// 任务本身相关
	Name      string `gorm:"uniqueIndex:mission_name_ns"` // 名称
	Desc      string // 描述
	Namespace string `gorm:"uniqueIndex:mission_name_ns"` // 命名空间(默认为default)
	Total     uint   // 任务总分（默认值为100）
	Guide     string `gorm:"type:text"` // 以markdown为格式的说明文档

	// K8S Deployment相关
	Deployment         uint   // k8s部署文件
	ExecContainer      string // 默认执行的容器（默认为空并访问首个容器）
	Command            string // WebShell执行的命令 TODO 支持shell脚本
	WhiteListContainer string // 白名单容器（若为空则放行首个容器）

	// VNC
	VNCEnable    bool   // 启用VNC桌面访问
	VNCContainer string // VNC目标容器
	VNCPort      string // VNC目标接口
}

// 任务构造选项
type MissionBuildOpt func(m *Mission) (err error)

// 任务构造器
func MissionBuilder(_ context.Context, name string, dp *Deployment, opts ...MissionBuildOpt) (m *Mission, err error) {
	if dp == nil || dp.ID == 0 {
		return nil, ErrMissionDeploymentNotExist
	}
	m = &Mission{
		Name:       name,
		Namespace:  "default",
		Total:      100,
		Deployment: dp.ID,
	}
	for _, opt := range opts {
		if err = opt(m); err != nil {
			return
		}
	}
	if m.Deployment == 0 {
		return nil, ErrMissionDeploymentNotExist
	}
	return
}

// 启用VNC远程桌面
func MissionOptVnc(container, port string) MissionBuildOpt {
	return func(m *Mission) (err error) {
		m.VNCEnable = true
		m.VNCContainer = container
		m.VNCPort = port
		return
	}
}

// 任务描述
func MissionOptDesc(desc string) MissionBuildOpt {
	return func(m *Mission) (err error) {
		m.Desc = desc
		return
	}
}

// 任务命名空间
func MissionOptNs(ns string) MissionBuildOpt {
	return func(m *Mission) (err error) {
		if ns == "" {
			ns = defaultDepartmentNamespace
		}
		m.Namespace = ns
		return
	}
}

// 任务Deployment
func MissionOptDeployment(cmd, execC string, whiteListC []string) MissionBuildOpt {
	return func(m *Mission) (err error) {
		if len(whiteListC) > 0 {
			for _, v := range whiteListC {
				if strings.ContainsRune(v, ';') {
					return fmt.Errorf("设置白名单容器不允许包含';'字符: %s", v)
				}
			}
		}
		if cmd == "" {
			cmd = "bash"
		}
		m.ExecContainer = execC
		m.Command = cmd
		m.WhiteListContainer = strings.Join(whiteListC, ";")
		return
	}
}

// 任务分数
func MissionOptTotal(total uint) MissionBuildOpt {
	return func(m *Mission) (err error) {
		m.Total = total
		return
	}
}

// 任务创建或更新内部实现
func (m *Mission) CreateOrUpdate(ctx context.Context) (err error) {
	db := GetGlobalDB().WithContext(ctx)
	if err = db.Transaction(func(tx *gorm.DB) error {
		// 检查Deployment是否存在
		if _err := tx.First(new(Deployment), m.Deployment).Error; _err != nil {
			return ErrMissionDeploymentNotExist
		}

		// 创建
		if m.ID == 0 {
			if _err := tx.Create(m).Error; _err != nil {
				return _err
			}
			return nil
		}

		// 更新
		_err := tx.Save(m).Error

		return _err
	}); err != nil {
		return
	}
	return
}

// 创建并更新任务
func CrateOrUpdateMission(ctx context.Context, name string, dp *Deployment, opts ...MissionBuildOpt) (m *Mission, err error) {
	m, err = MissionBuilder(ctx, name, dp, opts...)
	if err != nil {
		return
	}
	if err = m.CreateOrUpdate(ctx); err != nil {
		return
	}
	return
}

// 批量查询任务
func ListMissions(ctx context.Context, name string, ns []string, builder *PageBuilder) (ms []*Mission, err error) {
	db := GetGlobalDB().WithContext(ctx)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if len(ns) > 0 {
		db = db.Where("namespace IN ?", ns)
	}
	if builder != nil {
		db = builder.build(db)
	}
	err = db.Find(&ms).Error
	return
}

// 根据ID获取任务
func GetMission(ctx context.Context, id uint) (ms *Mission, err error) {
	ms = new(Mission)
	err = GetGlobalDB().WithContext(ctx).First(ms, id).Error
	return
}

// 校验容器是否在白名单内
func (m *Mission) IsContainerAllowed(container string) (res bool) {
	if strings.ContainsRune(container, ';') {
		return false
	}
	WhiteListC := strings.Split(m.WhiteListContainer, ";")
	for _, v := range WhiteListC {
		if strings.EqualFold(container, v) {
			return true
		}
	}
	return false
}

// TODO CMD的分割
func (m *Mission) GetCommand() (cmd []string) {
	return []string{m.Command}
}

// 根据mission的deployment获取可用的容器
func (m *Mission) ListAllowedContainers(ctx context.Context) (res []v1.Container, err error) {
	if m.Deployment == 0 {
		err = errors.New("任务不存在deployment")
		return
	}
	dp, err := GetDeployment(ctx, m.Deployment)
	if err != nil {
		return
	}
	dpCfg, err := k8s.ParseDeploymentConfig(dp.Raw, false)
	if err != nil {
		return
	}
	containers := dpCfg.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		err = errors.New("该任务不存在可用容器")
		return
	}
	for _, c := range containers {
		if m.IsContainerAllowed(c.Name) {
			res = append(res, c)
		}
	}
	return
}

// 获取任务名字的映射
func GetMissionsNameMapper(ctx context.Context, id ...uint) (res map[uint]string, err error) {
	type api struct {
		ID   uint
		Name string
	}
	if len(id) == 0 {
		return nil, errors.New("没有任务ID参数")
	}

	var data = make([]*api, 0)
	if err = GetGlobalDB().WithContext(ctx).Model(new(Mission)).Where(
		"id IN ?", id).Find(&data).Error; err != nil {
		return
	}

	res = make(map[uint]string, len(data))
	for _, v := range data {
		res[v.ID] = v.Name
	}
	return
}

// 统计实验数量
func CountMissions(ctx context.Context, name string, ns []string) (res int64, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Mission))
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if len(ns) > 0 {
		db = db.Where("namespace IN ?", ns)
	}
	err = db.Count(&res).Error
	return
}

// 删除任务 TODO 删除正在运行的Deployment
func DeleteMission(ctx context.Context, id uint) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Unscoped().Delete(new(Mission), id).Error
	return
}

/*
	任务和检查点
*/

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

// 获取所有任务的命名空间
func ListMissionNamespaces(ctx context.Context) (res []string, err error) {
	err = GetGlobalDB().WithContext(ctx).Model(new(Mission)).Pluck("namespace", &res).Error
	return
}
