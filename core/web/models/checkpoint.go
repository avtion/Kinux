package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

func init() {
	migrateQueue = append(migrateQueue, new(Checkpoint))
}

// 检查点方式
type CheckpointMethod = uint

const (
	_                CheckpointMethod = iota
	MethodExec                        // 用户执行的命令
	MethodStdout                      // 终端输出的结果
	MethodTargetPort                  // 指定端口
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
	In     string // 输入 - 可以是指令，也可以是目标端口
	Out    string // 输出 - 监听的目标输出，可以是指令执行结果，也可以是目标端口的响应
	Method uint   // 对应 CheckpointMethod
}

// 创建检查点
func (c *Checkpoint) Create(ctx context.Context) (err error) {
	db := GetGlobalDB().WithContext(ctx)
	switch c.Method {
	case MethodExec, MethodTargetPort:
		err = db.Select("CreatedAt", "Name", "Desc", "Method", "In").Create(c).Error
	case MethodStdout:
		err = db.Select("CreatedAt", "Name", "Desc", "Method", "Out").Create(c).Error
	default:
		err = errors.New("未知检查方式")
	}
	return
}

// 根据id查询所有的检查点
func FindCheckpoints(ctx context.Context, ids ...uint) (cps []*Checkpoint, err error) {
	if len(ids) == 0 {
		return
	}
	err = GetGlobalDB().WithContext(ctx).Model(new(Checkpoint)).Find(&cps, ids).Error
	return
}

// 获取检查点名字的映射
func GetCheckpointsNameMapper(ctx context.Context, id ...uint) (res map[uint]string, err error) {
	type api struct {
		ID   uint
		Name string
	}
	if len(id) == 0 {
		return nil, errors.New("没有检查点ID参数")
	}

	var data = make([]*api, 0)

	if err = GetGlobalDB().WithContext(ctx).Model(new(Checkpoint)).Where(
		"id IN ?", id).Find(&data).Error; err != nil {
		return
	}

	res = make(map[uint]string, len(data))
	for _, v := range data {
		res[v.ID] = v.Name
	}
	return
}
