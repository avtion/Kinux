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
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func init() {
	migrateQueue = append(migrateQueue, &Account{}, &Profile{})
}

var (
	ErrAccountPasswordWrong = errors.New("用户密码错误")
)

// 用户
type Account struct {
	gorm.Model
	Username string `gorm:"unique"` // 用户名
	Password string // 密码
	Role     uint   // 角色
	Profile  uint   `gorm:"unique"` // 个人资料
}

// 个人资料
type Profile struct {
	gorm.Model
	RealName   string // 真实姓名
	Department uint   // 部门
}

// 用户和个人资料
type AccountWithProfile struct {
	Account
	Profile
}

// 密码操作对象
type Password struct {
	Raw        string
	Salt       string
	Iterations int
}

// 校验是否符合规则
func (p *Password) Check() (err error) {
	if strings.EqualFold(strings.TrimSpace(p.Raw), "") {
		return errors.New("密码为空")
	}
	// TODO 完成
	return
}

// 加密
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

// 校验
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

// 创建新用户
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

// 根据用户名查询用户
func (a *Account) GetByUsername(ctx context.Context) (err error) {
	if a.Username == "" {
		return errors.New("username is null")
	}
	return GetGlobalDB().WithContext(ctx).First(a, &Account{Username: a.Username}).Error
}

// 校验用户密码
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
