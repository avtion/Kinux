package controllers

import (
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 任务查询
type MissionQuery struct {
	Name      string   `form:"name" `
	Namespace []string `form:"namespace"`
	Page      uint     `form:"page"`
	Size      uint     `form:"size"`
}

// 查询任务
func QueryMissions(c *gin.Context) {
	var query = new(MissionQuery)
	if err := c.ShouldBind(query); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	ms, err := services.ListMissions(c, ac, query.Name, query.Namespace, int(query.Page), int(query.Size))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(ms))
	return
}
