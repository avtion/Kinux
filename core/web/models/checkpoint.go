package models

import "gorm.io/gorm"

func init() {
	migrateQueue = append(migrateQueue, new(Checkpoint))
}

// 检查点方式
type CheckpointMethod = uint

const (
	_ CheckpointMethod = iota
	MethodExec
	MethodStdout
	MethodTargetPort
)

var _ = [...]CheckpointMethod{
	MethodExec,
	MethodStdout,
	MethodTargetPort,
}

// 检查点
type Checkpoint struct {
	gorm.Model
	Name   string // 名称
	Desc   string // 描述
	Method uint   // 对应 CheckpointMethod
	In     string // 输入 - 可以是指令，也可以是目标端口
	Out    string // 输出 - 监听的目标输出，可以是指令执行结果，也可以是目标端口的响应
}
