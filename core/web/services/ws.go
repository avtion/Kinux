// WebSocket业务层
package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"net/http"
)

type WsSvcFn func(c *gin.Context, w *websocket.Conn) (err error)

// webSocket默认升级器
var defaultUpgrader = websocket.Upgrader{
	// TODO 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 单独入口
func NewWebSocketService(c *gin.Context, fns ...WsSvcFn) (err error) {
	// 协议升级
	wsConn, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func() {
		if _err := wsConn.Close(); _err != nil {
			logrus.Error(err)
			return
		}
	}()
	logrus.WithField("wsInfo", wsConn).Debug("WebSocket连接建立成功")

	// 按照函数对象执行对应的方法
	for _, fn := range fns {
		// 每次执行函数之前都检验一次上下文是否终止
		select {
		case <-c.Done():
			return c.Err()
		default:
		}
		if err = fn(c, wsConn); err != nil {
			return
		}
	}
	return
}

// 用于操作Websocket的对象
type _wsOperator struct{}

var WsOperator = new(_wsOperator)

// 创建新的K8S POD连接
func (_wsOperator) NewK8SPodConnection(account, job, container string, cmd []string) WsSvcFn {
	return func(c *gin.Context, w *websocket.Conn) (err error) {
		// 再次升级成容器连接
		// TODO 支持多写者和多读者
		cs := NewContainerSessionAdapter(c, w)

		// 获取默认POD
		pods, err := k8s.GetPods(c, account, job)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err, msg.WithLogPrint(logrus.ErrorLevel)))
			return
		}
		if len(pods.Items) == 0 {
			err = errors.New("the length of pod items is zero")
			return
		}
		pod := pods.Items[0]

		// 将升级的Session连接连接到对应的K8S容器
		if err = k8s.ConnectToPod(c, &pod, container, cs, cmd); err != nil {
			return
		}

		return
	}
}
