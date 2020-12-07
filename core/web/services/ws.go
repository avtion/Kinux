// WebSocket业务层
package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type shell struct {
}

// webSocket默认升级器
var defaultUpgrader = websocket.Upgrader{
	// TODO 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 连接到具体的容器
func WebSocketContainerService(c *gin.Context, account, job string, container string, cmd []string) {
	// 协议升级
	wsConn, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, msg.BuildFailed(err, msg.WithLogPrint(logrus.ErrorLevel)))
		return
	}
	defer func() {
		_ = wsConn.Close()
	}()
	logrus.Debug("websocket连接建立成功")

	// 再次升级成容器连接
	cs := NewContainerSessionAdapter(c, wsConn)

	// 获取默认POD
	pods, err := k8s.GetPods(c, account, job)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err, msg.WithLogPrint(logrus.ErrorLevel)))
		return
	}
	if len(pods.Items) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(
			"the length of pod items is zero", msg.WithLogPrint()))
		return
	}
	pod := pods.Items[0]

	// 将升级的Session连接连接到对应的K8S容器
	if err = k8s.ConnectToPod(c, &pod, container, cs, cmd); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err, msg.WithLogPrint(logrus.ErrorLevel)))
		return
	}
}
