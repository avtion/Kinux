package models

import (
	"Kinux/tools"
	"context"
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
func CrateDeployment(ctx context.Context, name string, raw []byte) (id uint, err error) {
	// 校验K8S配置文件
	if err = yaml.UnmarshalStrict(raw, new(appV1.Deployment)); err != nil {
		return
	}
	if name == "" {
		name = tools.GetRandomString(12)
	}
	dp := &Deployment{
		Name: name,
		Raw:  raw,
	}
	if err = GetGlobalDB().WithContext(ctx).Create(dp).Error; err != nil {
		return
	}
	return dp.ID, nil
}
