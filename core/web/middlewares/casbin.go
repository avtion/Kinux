package middlewares

import (
	"Kinux/core/web/models"
	"Kinux/tools/cfg"
	"errors"
	"github.com/casbin/casbin/v2"
	casbinLog "github.com/casbin/casbin/v2/log"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
	"sync"
)

// 模型数据
const modelStr = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`

type _enforcer struct {
	*casbin.Enforcer
	sync.Once
}

var (
	// 全局处理器
	globalEnforcer = new(_enforcer)
	ErrEnforcerNil = errors.New("enforcer is not ready")
)

// 全局获取处理器
// WARING: 无法确保处理器能够被正常初始化，务必校验处理器对象是否为空值
func GetGlobalEnforcer() *casbin.Enforcer {
	// 初始化Casbin处理器
	globalEnforcer.Do(func() {
		var _err error
		globalEnforcer.Enforcer, _err = newEnforcer(&gormAdapter{DB: models.GetGlobalDB()})
		if _err != nil {
			logrus.Error(_err)
			return
		}
	})
	return globalEnforcer.Enforcer
}

// 规则处理
func Enforce(args ...interface{}) (ok bool, err error) {
	// 如果配置未启动Casbin则直接通过
	if !cfg.DefaultConfig.Casbin.Enable {
		return true, nil
	}

	// 初始化Casbin处理器
	e := GetGlobalEnforcer()
	if e == nil {
		err = ErrEnforcerNil
		return
	}
	return e.Enforce(args...)
}

// 创建新的Casbin规则处理器
func newEnforcer(dbAdapter *gormAdapter) (e *casbin.Enforcer, err error) {
	// 读取模型
	model, err := casbinModel.NewModelFromString(modelStr)
	if err != nil {
		return
	}

	// 检测数据库表是否存在
	if !dbAdapter.Migrator().HasTable(new(CasbinRule)) {
		// 创建数据库表
		if err = dbAdapter.createTable(); err != nil {
			return
		}
	}

	// 创建处理器
	e, err = casbin.NewEnforcer(model, dbAdapter)
	if err != nil {
		return
	}

	// 设置日志
	e.SetLogger(&logrusAdapter{
		Entry:  logrus.StandardLogger().WithField("module", "casbin"),
		enable: true,
	})
	return
}

// 初始化Casbin的角色继承机制
func initCasbinRoles(e *casbin.Enforcer) (err error) {
	for i := len(models.RoleArray) - 1; i > 0; i-- {
		if _, err = e.AddGroupingPolicy(
			cast.ToString(models.RoleArray[i]),
			cast.ToString(models.RoleArray[i-1]),
		); err != nil {
			return
		}
	}
	return
}

// logrus日志适配器
type logrusAdapter struct {
	*logrus.Entry
	enable bool
}

var _ casbinLog.Logger = (*logrusAdapter)(nil)

func (l *logrusAdapter) EnableLog(enable bool) {
	l.enable = enable
}

func (l *logrusAdapter) IsEnabled() bool {
	return l.enable
}

func (l *logrusAdapter) LogModel(model [][]string) {
	l.Info(model)
}

func (l *logrusAdapter) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	l.WithFields(logrus.Fields{
		"matcher":  matcher,
		"request":  request,
		"result":   result,
		"explains": explains,
	}).Debug("casbin enforce debug")
}

func (l *logrusAdapter) LogRole(roles []string) {
	l.Info(roles)
}

func (l *logrusAdapter) LogPolicy(policy map[string][][]string) {
	l.Info(policy)
}

// Casbin规则的定义
type CasbinRule struct {
	gorm.Model
	PType string `gorm:"size:40;uniqueIndex:unique_index"`
	V0    string `gorm:"size:40;uniqueIndex:unique_index"`
	V1    string `gorm:"size:40;uniqueIndex:unique_index"`
	V2    string `gorm:"size:40;uniqueIndex:unique_index"`
	V3    string `gorm:"size:40;uniqueIndex:unique_index"`
	V4    string `gorm:"size:40;uniqueIndex:unique_index"`
	V5    string `gorm:"size:40;uniqueIndex:unique_index"`
}

// GORM V2适配器
type gormAdapter struct {
	*gorm.DB
}

// 将数据库数据转换成规则
func (*gormAdapter) loadPolicyLine(line CasbinRule, model casbinModel.Model) {
	var p = []string{line.PType, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}

	var lineText string
	if line.V5 != "" {
		lineText = strings.Join(p, ", ")
	} else if line.V4 != "" {
		lineText = strings.Join(p[:6], ", ")
	} else if line.V3 != "" {
		lineText = strings.Join(p[:5], ", ")
	} else if line.V2 != "" {
		lineText = strings.Join(p[:4], ", ")
	} else if line.V1 != "" {
		lineText = strings.Join(p[:3], ", ")
	} else if line.V0 != "" {
		lineText = strings.Join(p[:2], ", ")
	}

	persist.LoadPolicyLine(lineText, model)
}

// 创建数据库表
func (g *gormAdapter) createTable() (err error) {
	return g.AutoMigrate(new(CasbinRule))
}

// 删除数据库表
func (g *gormAdapter) dropTable() (err error) {
	entry := new(CasbinRule)
	if !g.Migrator().HasTable(entry) {
		return
	}
	return g.Migrator().DropTable(entry)
}

// 读取规则
func (g *gormAdapter) LoadPolicy(model casbinModel.Model) (err error) {
	var rules []CasbinRule
	if err = g.Order("id").Find(&rules).Error; err != nil {
		return
	}
	for _, rule := range rules {
		g.loadPolicyLine(rule, model)
	}
	return nil
}

// 创建新的CasbinRule
func (*gormAdapter) newPolicyLine(ptype string, rule []string) CasbinRule {
	var line CasbinRule

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// 保存规则
func (g *gormAdapter) SavePolicy(model casbinModel.Model) (err error) {
	if err = g.dropTable(); err != nil {
		return
	}
	if err = g.createTable(); err != nil {
		return
	}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := g.newPolicyLine(ptype, rule)
			if err = g.Create(&line).Error; err != nil {
				return err
			}
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := g.newPolicyLine(ptype, rule)
			if err = g.Create(&line).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// 添加规则 可选
func (g *gormAdapter) AddPolicy(_ string, ptype string, rule []string) error {
	line := g.newPolicyLine(ptype, rule)
	return g.Create(&line).Error
}

// 移除规则 可选
func (g *gormAdapter) RemovePolicy(_ string, ptype string, rule []string) error {
	line := g.newPolicyLine(ptype, rule)
	return g.Unscoped().Where(&line).Delete(&line).Error
}

// 移除指定规则 可选
func (g *gormAdapter) RemoveFilteredPolicy(_ string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := new(CasbinRule)
	line.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	return g.Unscoped().Where(&line).Delete(&line).Error
}

// 批量追加规则 可选
func (g *gormAdapter) AddPolicies(_ string, ptype string, rules [][]string) error {
	return g.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			line := g.newPolicyLine(ptype, rule)
			if err := tx.Create(&line).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 批量删除规则 可选
func (g *gormAdapter) RemovePolicies(_ string, ptype string, rules [][]string) error {
	return g.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			line := g.newPolicyLine(ptype, rule)
			if err := tx.Unscoped().Where(&line).Delete(&line).Error; err != nil { //can't use db.Delete as we're not using primary key http://jinzhu.me/gorm/crud.html#delete
				return err
			}
		}
		return nil
	})
}

var (
	_ persist.Adapter      = (*gormAdapter)(nil)
	_ persist.BatchAdapter = (*gormAdapter)(nil)
)
