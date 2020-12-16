package controllers

import (
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 账号
type Account struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登陆入口
func LoginAccount(c *gin.Context) {
	var ac = new(Account)
	if err := c.ShouldBindJSON(ac); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	token, err := services.LoginAccount(c, ac.Username, ac.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"msg":   "🛫️登陆成功",
		"token": token,
	})
	return
}
