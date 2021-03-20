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

// 班级
type Department struct {
	gorm.Model
	Name string `gorm:"unique"` // 名称
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
func ListDepartments(ctx context.Context, name string, page *PageBuilder) (
	ds []*Department, err error) {
	db := GetGlobalDB().WithContext(ctx)

	// 分页
	if page != nil {
		db = page.Build(db)
	}

	// 模糊搜索
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	err = db.Find(&ds).Error
	return
}

// 班级设置名称
func DepartmentNameOpt(name string) DepartmentOpt {
	return func(d *Department) error {
		if name = strings.TrimSpace(name); strings.EqualFold(name, "") {
			return errors.New("班级名不能为空")
		}

		d.Name = name
		return nil
	}
}

// 查询班级的总量
func CountDepartments(ctx context.Context, name string) (res int64, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Department))
	// 模糊搜索
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	err = db.Count(&res).Error
	return
}

// 更新班级
func UpdateDepartment(ctx context.Context, id int, opts ...DepartmentOpt) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	var d = new(Department)
	for _, fn := range opts {
		if err = fn(d); err != nil {
			return
		}
	}
	return GetGlobalDB().WithContext(ctx).Model(&Department{
		Model: gorm.Model{
			ID: uint(id),
		},
	}).UpdateColumns(d).Error
}

// 删除班级
func DeleteDepartment(ctx context.Context, id int) (err error) {
	// TODO 级联删除
	return GetGlobalDB().Unscoped().WithContext(ctx).Delete(new(Department), id).Error
}

// 快速返回选项数据结果
type QuickListDepartmentRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// 用于快速返回班级相关的选项数据
func QuickListDepartment(ctx context.Context, filters ...func(db *gorm.DB) *gorm.DB) (res []*QuickListDepartmentRes, err error) {
	err = GetGlobalDB().Model(new(Department)).WithContext(ctx).Select(
		"id, name").Scopes(filters...).Find(&res).Error
	return
}

// 根据ID获取班级
func GetDepartmentByID(ctx context.Context, id uint) (d *Department, err error) {
	d = new(Department)
	err = GetGlobalDB().WithContext(ctx).First(d, id).Error
	return
}

// 获取班级课程的ID
func (d *Department) GetLessonIDs(ctx context.Context) (res []uint, err error) {
	if err = GetGlobalDB().WithContext(ctx).Model(new(LessonDepartment)).Where(
		&LessonDepartment{Department: d.ID}).Pluck("lesson", &res).Error; err != nil {
		return
	}
	return
}

// 获取班级的课程
func (d *Department) GetLessons(ctx context.Context) (res []*Lesson, err error) {
	// 先获取课程拥有的LessonID
	lessonIDs, err := d.GetLessonIDs(ctx)
	if err != nil {
		return
	}
	if len(lessonIDs) == 0 {
		return
	}
	err = GetGlobalDB().WithContext(ctx).Model(new(Lesson)).Where(
		"lesson IN ?", lessonIDs).Find(&res).Error
	return
}
