package models

import (
	"context"
	"errors"
	"fmt"
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

// 获取检查点
func ListCheckpoints(ctx context.Context, name string, method CheckpointMethod, builder *PageBuilder) (
	res []*Checkpoint, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Checkpoint))
	if builder != nil {
		db = builder.build(db)
	}
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if method != 0 {
		db = db.Where("method = ?", method)
	}
	err = db.Find(&res).Error
	return
}

// 统计检查点
func CountCheckpoints(ctx context.Context, name string, method CheckpointMethod) (res int64, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Checkpoint))
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if method != 0 {
		db = db.Where("method = ?", method)
	}
	err = db.Count(&res).Error
	return
}

// 删除检查点
func DeleteCheckpoint(ctx context.Context, id uint) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	return GetGlobalDB().Unscoped().WithContext(ctx).Delete(new(Checkpoint), id).Error
}

// 检查点选项结果
type checkpointQuickResult struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 快速获取检查点选项
func QuickListCheckpoint(ctx context.Context, name string, method CheckpointMethod) (
	res []*checkpointQuickResult, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Checkpoint))
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if method != 0 {
		db = db.Where("method = ?", method)
	}
	err = db.Find(&res).Error
	return
}

// 更新检查点
func (c *Checkpoint) Edit(ctx context.Context) (err error) {
	if c.ID == 0 {
		return errors.New("id为空")
	}
	switch c.Method {
	case MethodTargetPort, MethodStdout, MethodExec:
	default:
		return errors.New("检查点方式为空")
	}
	return GetGlobalDB().WithContext(ctx).Save(c).Error
}
