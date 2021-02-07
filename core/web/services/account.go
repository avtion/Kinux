package services

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"errors"
	"fmt"
	GinJWT "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	ErrAccountVerifyFailed = errors.New("用户账号或密码错误")
)

func init() {
	RegisterWebsocketOperation(wsOpAuth, JWTRegister)
}

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

// websocket的JWT鉴权处理器注册器
func JWTRegister(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	rawToken := any.Get("data").ToString()
	if strings.TrimSpace(rawToken) == "" {
		return errors.New("密钥为空")
	}

	// 解析密钥
	ws.userToken, err = middlewares.TokenCentral.ParseTokenString(rawToken)
	if err != nil {
		return
	}

	// 解构用户参数
	userPayload := middlewares.ClaimsToTokenPayload(GinJWT.ExtractClaimsFromToken(ws.userToken))

	// 将用户信息写在上下文
	ws.Set(middlewares.TokenIdentityKey, userPayload)

	// 向用户发送通知
	if err = ws.SendMsg(msg.BuildSuccess(fmt.Sprintf("%s您好，websocket通信建立成功！", userPayload.Username))); err != nil {
		return
	}
	return
}
