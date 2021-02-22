package controllers

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"Kinux/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 账号
type Account struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登陆入口
func LoginAccount(c *gin.Context) {
	var _ac = new(Account)
	if err := c.ShouldBindJSON(_ac); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 登陆鉴权
	ac, err := services.LoginAccount(c, _ac.Username, _ac.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 获取密钥Token
	token, ttl, err := middlewares.TokenCentral.TokenGenerator(&middlewares.TokenPayload{
		Username: ac.Username,
		Role:     ac.Role,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 用户真实姓名
	var realName, avatarSeed string
	if profile, _err := ac.GetProfile(c); _err == nil {
		realName = profile.RealName
		avatarSeed = profile.AvatarSeed
	}

	// 用户部门
	var department string
	if dp, _err := ac.GetDepartment(c); _err == nil {
		department = dp.Name
	}

	c.JSON(http.StatusOK, msg.BuildSuccess(map[string]string{
		"msg":        "🛫️登陆成功",
		"token":      token,
		"ttl":        strconv.FormatInt(ttl.Unix(), 10),
		"username":   ac.Username,
		"realName":   realName,
		"role":       models.RoleTranslator(ac.Role),
		"department": department,
		"avatarSeed": avatarSeed,
	}))
	return
}

// 更新用户的头像种子
func UpdateAccountAvatarSeed(c *gin.Context) {
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	seed := tools.GetRandomString(6)
	if err = ac.UpdateAvatarSeed(c, seed); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(seed))
}

// 更新用户密码
func UpdatePassword(c *gin.Context) {
	passwordSetter := &struct {
		Old string `json:"old"`
		New string `json:"new"`
	}{}
	if err := c.ShouldBindJSON(passwordSetter); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 校验旧密码
	if err = ac.Verify(c, passwordSetter.Old); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("旧密码错误"))
		return
	}
	if err = ac.UpdatePassword(c, passwordSetter.New); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 发送鉴权失败响应要求客户端重新登录
	c.JSON(http.StatusOK, msg.Build(msg.CodeJWTAuthFailed, "密码修改成功，请重新登录"))
}
