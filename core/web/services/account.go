package services

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"errors"
	"fmt"
	GinJWT "github.com/appleboy/gin-jwt/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"
)

var (
	ErrAccountVerifyFailed = errors.New("用户账号或密码错误")
)

func init() {
	RegisterWebsocketOperation(wsOpAuth, JWTRegister)
}

// LoginAccount 登陆账号
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

// GetAccountFromCtx 从上下文中获取用户账户
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

// JWTRegister websocket的JWT鉴权处理器注册器
func JWTRegister(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	rawToken := any.Get("data").ToString()
	if strings.TrimSpace(rawToken) == "" {
		return errors.New("密钥为空")
	}

	// 解析密钥
	token, err := middlewares.TokenCentral.ParseTokenString(rawToken)
	if err != nil {
		ws.RequireClientAuth()
		return
	}

	// 写入ws
	ws.userToken = token

	// 解构用户参数
	userPayload := middlewares.ClaimsToTokenPayload(GinJWT.ExtractClaimsFromToken(ws.userToken))
	ws.Account, err = models.GetAccountByUsername(ws, userPayload.Username)
	if err != nil {
		return
	}

	// 挂载调度中心
	if err = RegisterWsConn(int(ws.Account.ID), ws); err != nil {
		return
	}

	// 向用户发送通知
	if err = ws.SendMsg(msg.BuildSuccess(fmt.Sprintf("%s您好，websocket通信建立成功！", userPayload.Username))); err != nil {
		return
	}

	// 创建新协程守护刷新密钥
	go func() {
		const refreshT = 10 * time.Minute
		// 发送新密钥方法
		var sendNewTokenFn = func() {
			newToken, ttl, _err := middlewares.TokenCentral.TokenGenerator(userPayload)
			if _err != nil {
				logrus.WithField("payload", userPayload).WithField("err", _err).Error("用户JWT刷新失败")
				return
			}
			data, _ := jsoniter.Marshal(&WebsocketMessage{
				Op: wsOpRefreshToken,
				Data: map[string]interface{}{
					"token": newToken,
					"ttl":   strconv.FormatInt(ttl.Unix(), 10),
				},
			})

			ws.SendData(data)
		}

		// 脉冲定时器
		t := time.NewTicker(refreshT)
		defer t.Stop()

		claims := token.Claims.(jwt.MapClaims)
		oldTTL := cast.ToInt64(claims["exp"])
		if oldTTL == 0 {
			logrus.WithFields(logrus.Fields{
				"payload": userPayload,
				"err":     "无法确定用户原本密钥的TTL",
				"claims":  claims,
			}).Error("用户JWT刷新失败")
			return
		}

		// 密钥过期时间小于刷新时间则直接推送一次
		if time.Unix(oldTTL, 0).Sub(middlewares.TokenCentral.TimeFunc()) < refreshT {
			sendNewTokenFn()
		}

		// 循环发送
		for {
			select {
			case <-ws.Done():
				// 上下文结束也退出
				return
			case <-t.C:
				// 定期推送刷新新的密钥
				sendNewTokenFn()
			}
		}
	}()

	// 鉴权完毕之后检查一下用户是否有正在进行的考试
	go func() {
		_ew, isExist := ExamWatchers.Load(ws.Account.ID)
		if !isExist {
			return
		}
		ew, _ := _ew.(*ExamWatcher)
		raw, _ := jsoniter.Marshal(&WebsocketMessage{
			Op:   wsOpExamRunning,
			Data: NewExamRunningInfo(ew),
		})
		ws.SendData(raw)
	}()

	return
}
