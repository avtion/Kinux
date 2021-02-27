package models

import (
	"Kinux/tools"
	"bytes"
	"context"
	"errors"
	"gorm.io/gorm"
	appV1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
)

func init() {
	migrateQueue = append(migrateQueue, new(Deployment))
}

// K8S Deployment部署文件
type Deployment struct {
	gorm.Model
	Name string `gorm:"unique"` // 名称
	Raw  []byte // 配置字节流
}

// 存储新的K8S Deployment
func CrateOrUpdateDeployment(ctx context.Context, name string, raw []byte) (id uint, err error) {
	// 校验K8S配置文件
	if err = yaml.UnmarshalStrict(raw, new(appV1.Deployment)); err != nil {
		return
	}
	if name == "" {
		name = tools.GetRandomString(12)
	}

	var dp = new(Deployment)

	// 根据Name查询目标Deployment是否已经存在
	if err = GetGlobalDB().WithContext(ctx).Where(&Deployment{Name: name}).First(dp).Error; err != nil {
		// 创建目标Deployment
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dp = &Deployment{
				Name: name,
				Raw:  raw,
			}
			err = GetGlobalDB().WithContext(ctx).Create(dp).Error
			return dp.ID, err
		}
		return
	}

	// 如果存储的Deployment配置相同，则直接返回
	if bytes.Equal(dp.Raw, raw) {
		return dp.ID, nil
	}

	// 更新
	dp.Raw = raw
	if err = GetGlobalDB().WithContext(ctx).Save(dp).Error; err != nil {
		return
	}
	return dp.ID, nil
}

// 批量获取Deployment
func ListDeployment(ctx context.Context, name string, page *PageBuilder) (res []*Deployment, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Deployment))
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if page != nil {
		db = page.build(db)
	}
	err = db.Find(&res).Error
	return
}

// 指定ID获取
func GetDeployment(ctx context.Context, id uint) (res *Deployment, err error) {
	res = new(Deployment)
	err = GetGlobalDB().WithContext(ctx).First(res, id).Error
	return
}

// 修改Deployment
func EditDeployment(ctx context.Context, id uint, raw []byte) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Model(new(Deployment)).Where(
		"id = ?", id).Update("raw", raw).Error
	return
}

// 删除Deployment
func DeleteDeployment(ctx context.Context, id uint) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Unscoped().Delete(new(Deployment), id).Error
	return
}

// 新增Deployment
func AddDeployment(ctx context.Context, name string, raw []byte) (err error) {
	err = GetGlobalDB().WithContext(ctx).Create(&Deployment{
		Name: name,
		Raw:  raw,
	}).Error
	return
}

// 快速获取Deployment
func QuickListDeployment(ctx context.Context, name string) (res []*struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Deployment)).Select("id, name")
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	err = db.Find(&res).Error
	return
}

// 统计Deployment
func CountDeployment(ctx context.Context, name string) (res int64, err error) {
	db := GetGlobalDB().WithContext(ctx).Model(new(Deployment))
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	err = db.Count(&res).Error
	return
}
