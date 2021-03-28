package controllers

import (
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Websocket处理器
func WebSocketHandlerV1(c *gin.Context) {
	_, err := services.NewWebsocketSchedule(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
}

// 获取当前活跃状态的会话
func ListLiveWebsocket(c *gin.Context) {
	res := services.ListScheduleCenterInfo(c)
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
