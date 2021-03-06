package models

import (
	"Kinux/tools/cfg"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// 允许接入两种不同的关系型数据库
const (
	dbSelectorMySQL  = "mysql"
	dbSelectorSQLite = "sqlite"
)

var (
	globalDb     *gorm.DB
	migrateQueue []interface{}
)

// 获取全局数据库操作对象
func GetGlobalDB() *gorm.DB {
	if globalDb == nil {
		if err := InitDatabaseConn(context.Background(),
			cfg.DefaultConfig.Database.Name, cfg.DefaultConfig.Database.Dsn); err != nil {
			logrus.Error(err)
			panic(err)
		}
	}
	return globalDb
}

// 初始化数据库连接
func InitDatabaseConn(ctx context.Context, dbSelector, dsn string) (err error) {
	// 根据数据库类型选择不同的Dialector
	var dialector gorm.Dialector
	switch dbSelector {
	case dbSelectorMySQL:
		dialector = mysql.Open(dsn)
	case dbSelectorSQLite:
		dialector = sqlite.Open(dsn)
	default:
		return fmt.Errorf("数据库类型无法匹配: %s", dbSelector)
	}

	// TODO 适配logrus日志
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      false,
			LogLevel:      logger.Info,
		}),
	})
	if err != nil {
		return
	}

	// 数据库同步
	if err = db.AutoMigrate(migrateQueue...); err != nil {
		return
	}
	globalDb = db.WithContext(ctx)
	logrus.Trace("数据库初始化成功")
	return
}

// 分页构造器
type PageBuilder struct {
	Page, Size int
}

func (p *PageBuilder) build(db *gorm.DB) *gorm.DB {
	if p.Page == 0 || p.Size == 0 {
		return db
	}
	return db.Limit(p.Size).Offset((p.Page - 1) * p.Size)
}

func NewPageBuilder(page, size int) *PageBuilder {
	return &PageBuilder{
		Page: page,
		Size: size,
	}
}
