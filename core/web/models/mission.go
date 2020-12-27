package models

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

// 任务和检查点是一对多的关系
type MissionCheckpoints struct {
	gorm.Model
	Mission         uint   `gorm:"uniqueIndex:mission_checkpoint"`
	CheckPoint      uint   `gorm:"uniqueIndex:mission_checkpoint"`
	Percent         uint   // 该检查点占任务总分的百分比
	Priority        int    // 自定义排序
	TargetContainer string // 目标容器
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
		m.ExecContainer = execC
		m.Command = cmd
		m.WhiteListContainer = strings.Join(whiteListC, ";")
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
		db = builder.build(db)
	}
	err = db.Find(&ms).Error
	return
}
