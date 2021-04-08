package models

import (
	"Kinux/core/k8s"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
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
	Name  string `gorm:"unique"` // 名称
	Desc  string // 描述
	Total uint   // 任务总分（默认值为100）
	Guide string `gorm:"type:text"` // 以markdown为格式的说明文档

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
func ListMissions(ctx context.Context, name string, builder *PageBuilder) (ms []*Mission, err error) {
	db := GetGlobalDB().WithContext(ctx)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if builder != nil {
		db = builder.Build(db)
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
func CountMissions(ctx context.Context, name string) (res int64, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Mission))
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	err = db.Count(&res).Error
	return
}

// 删除任务
func DeleteMission(ctx context.Context, id uint) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Unscoped().Delete(&Mission{
		Model: gorm.Model{ID: id},
	}, id).Error
	return
}

// 创建任务
func AddMission(ctx context.Context, name string, dp uint, opts ...MissionBuildOpt) (err error) {
	deployment, err := GetDeployment(ctx, dp)
	if err != nil {
		return
	}
	m, err := MissionBuilder(ctx, name, deployment, opts...)
	if err != nil {
		return
	}
	return GetGlobalDB().WithContext(ctx).Create(m).Error
}

// 修改任务
func EditMission(ctx context.Context, id uint, name string, dp uint, opts ...MissionBuildOpt) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	deployment, err := GetDeployment(ctx, dp)
	if err != nil {
		return
	}
	m, err := MissionBuilder(ctx, name, deployment, opts...)
	if err != nil {
		return
	}
	return GetGlobalDB().WithContext(ctx).Model(&Mission{
		Model: gorm.Model{
			ID: id,
		},
	}).Updates(m).Error
}

// 修改任务的文档
func UpdateMissionGuide(ctx context.Context, id uint, text string) (err error) {
	return GetGlobalDB().WithContext(ctx).Model(new(Mission)).Where("id = ?", id).Update("guide", text).Error
}

// 批量获取实验
func GetMissions(ctx context.Context, fns ...func(db *gorm.DB) *gorm.DB) (res []*Mission, err error) {
	err = GetGlobalDB().WithContext(ctx).Scopes(fns...).Find(&res).Error
	return
}

// 删除钩子
func (m *Mission) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Unscoped().Delete(new(MissionCheckpoints), "mission = ?", m.ID).Error
}

var _ callbacks.BeforeDeleteInterface = (*Mission)(nil)
