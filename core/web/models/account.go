package models

import (
	"Kinux/tools"
	"Kinux/tools/bytesconv"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func init() {
	migrateQueue = append(migrateQueue, &Account{}, &Profile{})
}

var (
	ErrAccountPasswordWrong = errors.New("用户密码错误")
)

// Account 用户
type Account struct {
	gorm.Model
	Username string `gorm:"unique"` // 用户名
	Password string // 密码
	Role     uint   // 角色
	Profile  uint   `gorm:"unique"` // 个人资料
}

// Profile 个人资料
type Profile struct {
	gorm.Model
	RealName   string // 真实姓名
	Department uint   // 部门
	AvatarSeed string // 头像种子
}

// AccountWithProfile 用户和个人资料
type AccountWithProfile struct {
	Account
	Profile
}

// Password 密码操作对象
type Password struct {
	Raw        string
	Salt       string
	Iterations int
}

// Check 校验是否符合规则
func (p *Password) Check() (err error) {
	if strings.EqualFold(strings.TrimSpace(p.Raw), "") {
		return errors.New("密码为空")
	}
	// TODO 完成
	return
}

// Encode 加密
func (p *Password) Encode() (string, error) {
	// 一共三个参数，分别是原始密码、盐、迭代次数
	var salt, iterations = tools.GetRandomString(12), 180000

	if strings.TrimSpace(p.Salt) != "" {
		salt = p.Salt
	}

	// 确保盐不包含美元$符号
	if strings.Contains(salt, "$") {
		return "", errors.New("salt contains dollar sign ($)")
	}

	if p.Iterations > 0 {
		iterations = p.Iterations
	}

	// pbkdf2加密
	hash := pbkdf2.Key(bytesconv.StringToBytes(p.Raw), bytesconv.StringToBytes(salt),
		iterations, sha256.Size, sha256.New)

	// base64编码成为固定长度的字符串
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	// 最终字符串拼接成pbkdf2_sha256密钥格式
	return fmt.Sprintf("%s$%d$%s$%s", "pbkdf2_sha256", iterations, salt, b64Hash), nil
}

// Verify 校验
func (p *Password) Verify(encoded string) (bool, error) {
	// 输入两个参数，分别是原始密码、需要校验的密钥（数据库中存储的密码）
	// 输出校验结果（布尔值）、错误

	// 先根据美元$符号分割密钥为4个子字符串
	s := strings.Split(encoded, "$")

	// 如果分割结果不是4个子字符串，则认为不是pbkdf2_sha256算法的结果密钥，跳出错误
	if len(s) != 4 {
		return false, errors.New("hashed password components mismatch")
	}

	// 分割子字符串的结果分别为算法名、迭代次数、盐和base64编码
	// ---> 这里可以获得加密用的盐
	algorithm, iterations, salt := s[0], s[1], s[2]

	// 如果密钥算法名不是pbkdf2_sha256算法，跳出错误
	if algorithm != "pbkdf2_sha256" {
		return false, errors.New("algorithm mismatch")
	}

	// 将迭代次数转换成int数据类型 -->这里可以获得加密用的迭代次数
	i, err := strconv.Atoi(iterations)
	if err != nil {
		return false, errors.New("unreadable component in hashed password")
	}

	// 将原始密码用上面获取的盐、迭代次数进行加密
	p.Salt, p.Iterations = salt, i
	newEncoded, err := p.Encode()
	if err != nil {
		return false, err
	}

	// 最终用hmac.Equal函数判断两个密钥是否相同
	return hmac.Equal(bytesconv.StringToBytes(newEncoded), bytesconv.StringToBytes(encoded)), nil
}

// NewAccounts 创建新用户
func NewAccounts(ctx context.Context, acs ...*AccountWithProfile) (err error) {
	if len(acs) == 0 {
		return errors.New("创建新用户失败: 无可创建的用户")
	}
	for _, ac := range acs {
		if err = ac.newAccount(ctx, ac.Profile); err != nil {
			return fmt.Errorf("创建新用户(%s)发生错误: %v", ac.Username, err)
		}
	}
	return
}

// 创建新用户内部实现
func (a *Account) newAccount(ctx context.Context, p Profile) (err error) {
	if strings.EqualFold(strings.TrimSpace(a.Username), "") {
		return errors.New("用户名为空")
	}

	// 密码校验
	pw := &Password{Raw: a.Password}
	if err = pw.Check(); err != nil {
		return
	}

	// 避免匿名用户的创建
	if a.Role == RoleAnonymous {
		a.Role = RoleNormalAccount
	}

	// 密码加密
	a.Password, err = pw.Encode()
	if err != nil {
		return
	}

	// 使用事务创建用户
	if err = GetGlobalDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先创建个人资料
		if p.AvatarSeed == "" {
			p.AvatarSeed = tools.GetRandomString(6)
		}
		if _err := tx.Create(&p).Error; _err != nil {
			return _err
		}
		a.Profile = p.ID

		// 再创建用户
		if _err := tx.Create(a).Error; _err != nil {
			if strings.Contains(_err.Error(), "Duplicate") || strings.Contains(_err.Error(), "UNIQUE") {
				// 友好提示
				_err = errors.New("用户已经存在")
			}
			return _err
		}
		return nil
	}); err != nil {
		return
	}
	return
}

// GetByUsername 根据用户名查询用户
func (a *Account) GetByUsername(ctx context.Context) (err error) {
	if a.Username == "" {
		return errors.New("username is null")
	}
	return GetGlobalDB().WithContext(ctx).First(a, &Account{Username: a.Username}).Error
}

// Verify 校验用户密码
func (a *Account) Verify(ctx context.Context, pw string) (err error) {
	if a.Password == "" {
		// 查询用户数据
		if err = a.GetByUsername(ctx); err != nil {
			return
		}
		// 如果密码还是空的直接判断错误
		if a.Password == "" {
			return ErrAccountPasswordWrong
		}
	}
	ok, err := (&Password{Raw: pw}).Verify(a.Password)
	if err != nil {
		return
	}
	if !ok {
		return ErrAccountPasswordWrong
	}
	return
}

// GetProfile 获取用户的个人资料
func (a *Account) GetProfile(ctx context.Context) (p *Profile, err error) {
	p = new(Profile)
	err = GetGlobalDB().WithContext(ctx).First(p, a.Profile).Error
	return
}

// GetDepartment 根据用户资料获取对应的班级
func (p *Profile) GetDepartment(ctx context.Context) (d *Department, err error) {
	d = new(Department)
	err = GetGlobalDB().WithContext(ctx).First(d, p.Department).Error
	return
}

// GetDepartment 根据用户获取对应的班级（是 GetProfile 和 GetDepartment 的快捷方式）
func (a *Account) GetDepartment(ctx context.Context) (d *Department, err error) {
	p, err := a.GetProfile(ctx)
	if err != nil {
		return
	}
	return p.GetDepartment(ctx)
}

// ListAccounts 批量查询Accounts
func ListAccounts(ctx context.Context, builder *PageBuilder) (acs []*Account, err error) {
	db := GetGlobalDB().WithContext(ctx)
	if builder != nil {
		db = builder.Build(db)
	}
	err = db.Find(&acs).Error
	return
}

// GetAccountByUsername 根据用户名获取对应的用户
func GetAccountByUsername(ctx context.Context, username string) (ac *Account, err error) {
	ac = new(Account)
	err = GetGlobalDB().WithContext(ctx).Where(&Account{Username: username}).First(ac).Error
	return
}

// UpdateAvatarSeed 更新用户的头像种子
func (a *Account) UpdateAvatarSeed(ctx context.Context, seed string) (err error) {
	if seed == "" {
		seed = tools.GetRandomString(6)
	}
	if a.ID == 0 {
		return errors.New("用户ID为空")
	}
	return GetGlobalDB().WithContext(ctx).Model(new(Profile)).Where(&Profile{
		Model: gorm.Model{
			ID: a.ID,
		},
	}).Update("avatar_seed", seed).Error
}

// UpdatePassword 更新用户密码
func (a *Account) UpdatePassword(ctx context.Context, newPw string) (err error) {
	if a.ID == 0 {
		return errors.New("用户ID为空")
	}
	// 密码校验
	pw := &Password{Raw: newPw}
	if err = pw.Check(); err != nil {
		return
	}
	// 密码加密
	newPwEncode, err := pw.Encode()
	if err != nil {
		return
	}
	return GetGlobalDB().WithContext(ctx).Model(new(Account)).Where(
		"id = ?", a.ID).Update("Password", newPwEncode).Error
}

// AccountFilterFn 用户查询过滤器
type AccountFilterFn = func(db *gorm.DB) *gorm.DB

// AccountNameFilter 用户名过滤器
func AccountNameFilter(name string) AccountFilterFn {
	return func(db *gorm.DB) *gorm.DB {
		if name == "" {
			return db
		}
		likeParams := fmt.Sprintf("%%%s%%", name)
		return db.Where("accounts.username LIKE ? OR profiles.real_name LIKE ?", likeParams, likeParams)
	}
}

// AccountDepartmentFilter 用户班级过滤器
func AccountDepartmentFilter(id int) AccountFilterFn {
	return func(db *gorm.DB) *gorm.DB {
		if id == 0 {
			return db
		}
		return db.Where("profiles.department = ?", id)
	}
}

// AccountRoleFilter 用户角色过滤器
func AccountRoleFilter(level RoleIdentify) AccountFilterFn {
	return func(db *gorm.DB) *gorm.DB {
		if level == 0 {
			return db
		}
		return db.Where("accounts.role = ?", level)
	}
}

// AccountsListResult 用户列表结果
type AccountsListResult struct {
	ID           uint
	Role         uint
	Profile      uint
	Username     string
	RealName     string
	Department   string
	DepartmentId uint
	CreatedAt    time.Time
}

// 获取用户列表包括个人资料（内部实现）
func listAccountsWithProfiles(ctx context.Context, builder *PageBuilder, filters ...AccountFilterFn) (
	db *gorm.DB) {
	const selectQuery = `accounts.id, accounts.username, accounts.role, profiles.real_name, 
							departments.name AS department, departments.id AS department_id, 
							accounts.created_at, accounts.profile`
	const JoinQuery = `accounts
         LEFT JOIN profiles ON accounts.profile = profiles.id
         LEFT JOIN departments ON profiles.department = departments.id`
	db = GetGlobalDB().WithContext(ctx).Table("accounts").Select(selectQuery).Joins(JoinQuery).Scopes(filters...)
	if builder != nil {
		db = builder.Build(db)
	}
	return
}

// ListAccountsWithProfiles 获取用户列表包括个人资料
func ListAccountsWithProfiles(ctx context.Context, builder *PageBuilder, filters ...AccountFilterFn) (
	res []*AccountsListResult, err error) {
	err = listAccountsWithProfiles(ctx, builder, filters...).Scan(&res).Error
	return
}

// CountAccountsWithProfiles 统计用户列表包括个人资料
func CountAccountsWithProfiles(ctx context.Context, filters ...AccountFilterFn) (res int64, err error) {
	err = listAccountsWithProfiles(ctx, nil, filters...).Count(&res).Error
	return
}

// DeleteAccount 删除用户
func DeleteAccount(ctx context.Context, id int) (err error) {
	if id == 0 {
		return errors.New("id为空")
	}
	// TODO 删除用户正在执行的容器
	return GetGlobalDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ac = new(Account)
		if _err := tx.First(ac, id).Error; _err != nil {
			return _err
		}
		if _err := tx.Unscoped().Delete(new(Profile), ac.Profile).Error; _err != nil {
			return _err
		}
		if _err := tx.Unscoped().Delete(ac).Error; _err != nil {
			return _err
		}
		return nil
	})
}

// GetAccountByID 根据ID获取账号
func GetAccountByID(ctx context.Context, id int) (ac *Account, err error) {
	if id == 0 {
		return nil, errors.New("id为空")
	}
	ac = new(Account)
	err = GetGlobalDB().WithContext(ctx).First(ac, id).Error
	return
}

// Update 更新用户数据 - 仅支持更新用户名和角色
func (a *Account) Update(ctx context.Context) (err error) {
	if a.ID == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Model(new(Account)).Select(
		"username", "role").Where("id = ?", a.ID).Updates(a).Error
	return
}

// Update 更新用户资料 - 仅支持真实姓名和班级
func (p *Profile) Update(ctx context.Context) (err error) {
	if p.ID == 0 {
		return errors.New("id为空")
	}
	err = GetGlobalDB().WithContext(ctx).Model(new(Profile)).Select(
		"real_name, department").Where("id = ?", p.ID).Updates(p).Error
	return
}

// GetAccountsUsernameMapper 获取账号用户名和ID映射
func GetAccountsUsernameMapper(ctx context.Context, id ...uint) (res map[uint]string, err error) {
	type api struct {
		ID       uint
		Username string
	}
	if len(id) == 0 {
		return nil, errors.New("没有用户ID参数")
	}

	var data = make([]*api, 0)
	GetGlobalDB().WithContext(ctx).Model(new(Account)).Where("id IN ?", id).Find(&data)

	res = make(map[uint]string, len(data))
	for _, v := range data {
		if v.Username == "" {
			res[v.ID] = cast.ToString(v.ID)
		} else {
			res[v.ID] = v.Username
		}
	}

	return
}
