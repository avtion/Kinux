package services

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrAccountVerifyFailed = errors.New("用户账号或密码错误")
)

// 登陆账号
func LoginAccount(c *gin.Context, username, password string) (ac *models.Account, err error) {
	ac = &models.Account{Username: username}
	if err = ac.Verify(c, password); err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"password": password,
			"err":      err,
		}).Debug("用户身份验证失败")

		// 统一错误处理
		err = ErrAccountVerifyFailed
		return
	}
	return
}

// 从上下文中获取用户账户
func GetAccountFromCtx(c *gin.Context) (ac *models.Account, err error) {
	_u, isExist := c.Get(middlewares.TokenIdentityKey)
	if !isExist {
		err = errors.New("上下文不存在用户信息")
		return
	}
	u, ok := _u.(*middlewares.TokenPayload)
	if !ok {
		err = errors.New("上下用户信息转换失败")
		return
	}
	ac, err = models.GetAccountByUsername(c, u.Username)
	return
}
