package middlewares

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/tools"
	"Kinux/tools/bytesconv"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

// Gin上下文参数传递关键字
const TokenIdentityKey = "__jwt_account"

// 载荷
type TokenPayload struct {
	Username string
	Role     models.RoleIdentify
}

func (t *TokenPayload) MapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"username": t.Username,
		"role":     t.Role,
	}
}

func claimsToTokenPayload(claims jwt.MapClaims) *TokenPayload {
	username, _ := claims["username"]
	role, _ := claims["role"]
	return &TokenPayload{
		Username: cast.ToString(username),
		Role:     cast.ToUint(role),
	}
}

// 中间件
var JsonWebTokenAuth = NewJsonWebTokenAuth()

// 密钥中间件管理对象
var TokenCentral = &jwt.GinJWTMiddleware{
	Realm: "Kinux",
	// 加密的密钥采用随机加密生成
	// TODO 分布式密钥生成
	Key:         bytesconv.StringToBytes(tools.GetRandomString(12)),
	Timeout:     time.Hour,
	MaxRefresh:  time.Hour,
	IdentityKey: TokenIdentityKey,
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		// 用于创建jwt自定义载体
		// 由 jwt.TokenGenerator 触发
		if payload, ok := data.(*TokenPayload); ok {
			return payload.MapClaims()
		}
		return jwt.MapClaims{}
	},
	// 登陆使用，直接忽略
	Authenticator: nil,
	IdentityHandler: func(c *gin.Context) interface{} {
		return claimsToTokenPayload(jwt.ExtractClaims(c))
	},
	Authorizator: func(data interface{}, c *gin.Context) bool {
		// RBAC 验证权限
		// 执行顺序在IdentityHandler之后，data即IdentityHandler的返回值
		payload, ok := data.(*TokenPayload)
		if !ok {
			return false
		}
		// TODO RBAC鉴权
		logrus.WithField("payload", payload).Trace("JWT鉴权")
		return true
	},
	Unauthorized: func(c *gin.Context, code int, data string) {
		logrus.WithFields(logrus.Fields{
			"code": code,
			"msg":  data,
		}).Debug("鉴权失败")

		// 鉴权失败统一返回错误信息
		c.JSON(http.StatusOK, msg.Build(msg.CodeJWTAuthFailed, "系统鉴权失败"))
		return
	},
	TokenLookup: "header:Authorization",
	TimeFunc:    time.Now,
}

// 创建密钥管理中间件
func NewJsonWebTokenAuth() gin.HandlerFunc {
	md, err := jwt.New(TokenCentral)
	if err != nil {
		panic(err)
	}
	return md.MiddlewareFunc()
}
