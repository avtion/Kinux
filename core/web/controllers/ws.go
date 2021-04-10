package controllers

import (
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// WebSocketHandlerV1 Websocket处理器
func WebSocketHandlerV1(c *gin.Context) {
	_, err := services.NewWebsocketSchedule(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
}

// ListLiveWebsocket 获取当前活跃状态的会话
func ListLiveWebsocket(c *gin.Context) {
	res := services.ListScheduleCenterInfo(c)
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// SendMessageToTargetWs 向目标Websocket链接发送消息
func SendMessageToTargetWs(c *gin.Context) {
	params := &struct {
		TargetID int    `json:"target_id"`
		Text     string `json:"text"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	err := services.SendMessageToTargetWs(c, params.TargetID, params.Text)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("消息发送成功"))
}

// ForceTargetLogout 强制目标Websocket链接登出
func ForceTargetLogout(c *gin.Context) {
	params := &struct {
		TargetID int `json:"target_id"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	err := services.ForceTargetLogout(c, params.TargetID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("目标已登出"))
}
