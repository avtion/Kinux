package controllers

import (
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
)

// Websocket处理器
func WebSocketHandlerV1(c *gin.Context) {
	// TODO 身份校验 / 模式校验
	// TODO 指定容器
	services.WebSocketContainerService(c, "", "", "", []string{})
}
