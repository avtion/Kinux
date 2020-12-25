package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
)

func init() {
	migrateQueue = append(migrateQueue, new(Department))
}

const defaultDepartmentNamespace = "default"

// 班级
type Department struct {
	gorm.Model
	Name      string `gorm:"unique"` // 名称
	Namespace string // 可见的命名空间，以分号为间隔
}

type DepartmentOpt func(d *Department) error

// 创建新的Department
func newDepartment(name string, opts ...DepartmentOpt) (d *Department, err error) {
	d = &Department{Name: name}
	for _, opt := range opts {
		if err = opt(d); err != nil {
			return
		}
	}
	return
}

// 创建或更新
func (d *Department) CreateOrUpdate(ctx context.Context) (err error) {
	if d.Name == "" {
		return errors.New("department的名字为空")
	}
	db := GetGlobalDB().WithContext(ctx)
	if d.ID == 0 {
		if err = db.Create(d).Error; err != nil {
			return
		}
	}
	if err = db.Save(d).Error; err != nil {
		return
	}
	return
}

// 创建新的班级
func NewDepartment(ctx context.Context, name string, opts ...DepartmentOpt) (d *Department, err error) {
	d, err = newDepartment(name, opts...)
	if err != nil {
		return
	}
	if err = d.CreateOrUpdate(ctx); err != nil {
		return
	}
	return
}

// 批量获取班级
func ListDepartments(ctx context.Context, name string, page *PageBuilder, _ ...DepartmentOpt) (
	ds []*Department, err error) {
	db := GetGlobalDB().WithContext(ctx)

	// 分页
	if page != nil {
		db = page.build(db)
	}

	// 模糊搜索
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	err = db.Find(&ds).Error
	return
}

// 班级设置命名空间
func DepartmentNsOpt(ns ...string) DepartmentOpt {
	return func(d *Department) error {
		if len(ns) == 0 {
			return nil
		}
		for k, v := range ns {
			if strings.EqualFold(v, defaultDepartmentNamespace) {
				ns = append(ns[:k], ns[k+1:]...)
			}
			if strings.ContainsRune(v, ';') {
				return errors.New("命名空间包括分号")
			}
		}
		d.Namespace = strings.Join(append([]string{defaultDepartmentNamespace}, ns...), ";")
		return nil
	}
}
