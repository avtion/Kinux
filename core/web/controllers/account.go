package controllers

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
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
	var realName string
	if profile, _err := ac.GetProfile(c); _err == nil {
		realName = profile.RealName
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
	}))
	return
}
