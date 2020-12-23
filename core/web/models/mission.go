package models

import "gorm.io/gorm"

func init() {
	migrateQueue = append(migrateQueue, new(Mission), new(MissionCheckpoints))
}

// 任务
type Mission struct {
	gorm.Model
	Name            string // 名称
	Desc            string // 描述
	Namespace       string // 命名空间(默认为default)
	ExecContainer   string // 默认执行的容器（默认为空）
	Command         string // WebShell执行的命令
	ForbidContainer string // 允许访问的容器（空的情况下允许全部容器）
	Deployment      uint   // k8s部署文件
	Total           uint   // 任务总分（默认值为100）
	VNCEnable       bool   // 启用VNC桌面访问
	VNCContainer    string // VNC目标容器
	VNCPort         string // VNC目标接口
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
