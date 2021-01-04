package controllers

import (
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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

// 创建新的任务
// TODO 前端测试
func NewMission(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("任务id为空"))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	if err = services.AccountMissionOpera(c, ac, id, services.MissionCreate); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("任务创建成功"))
}

// 删除正在进行的任务
// TODO 前端测试
func DeleteMission(c *gin.Context) {
	id := cast.ToUint(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("任务id为空"))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// TODO 任务点检查

	if err = services.AccountMissionOpera(c, ac, id, services.MissionDelete); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("任务删除成功"))
}
