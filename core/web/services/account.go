package services

import (
	"Kinux/core/web/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrAccountVerifyFailed = errors.New("用户账号或密码错误")
)

// 登陆账号
func LoginAccount(c *gin.Context, username, password string) (token string, err error) {
	// TODO 校验账号和密码返回token和错误信息
	if err = (&models.Account{Username: username}).Verify(c, password); err != nil {
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
