// 容器相关的业务
package services

import (
	"Kinux/core/k8s"
	"Kinux/tools"
	"Kinux/tools/bytesconv"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"k8s.io/client-go/tools/remotecommand"
)

// ContainerSessionAdapter 实现 k8s.PtyHandler 接口
type ContainerSessionAdapter struct {
	id       string
	ctx      context.Context
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
	// TODO 考试监控
}

var _ k8s.PtyHandler = (*ContainerSessionAdapter)(nil)

// NewContainerSessionAdapter 创建一个 ContainerSessionAdapter 对象
func NewContainerSessionAdapter(ctx context.Context, conn *websocket.Conn) *ContainerSessionAdapter {
	adapter := &ContainerSessionAdapter{
		id:       tools.GetRandomString(12),
		ctx:      ctx,
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return adapter
}

// Done 结束当前 websocket 连接（由k8s组件主动调用）
func (t *ContainerSessionAdapter) Done() {
	// TODO 判断用户存活并结束掉对应的 Deployment
	close(t.doneChan)
	return
}

// Next 返回当前终端窗口的大小（需要使用阻塞，默认情况下k8s调用是不断尝试的）
func (t *ContainerSessionAdapter) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// Read 实现 io.Reader 接口将 webSocket 连接中的数据进行处理并拷贝到 p 中
func (t *ContainerSessionAdapter) Read(p []byte) (int, error) {
	// 读取数据并反序列化
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		logrus.Errorf("read message err: %v", err)
		return copy(p, EndOfTransmission), err
	}
	var msg = new(WebsocketMessage)
	if err := json.Unmarshal(message, msg); err != nil {
		logrus.Errorf("read parse message err: %v", err)
		return copy(p, EndOfTransmission), err
	}
	logrus.Trace("接收到数据", msg.Data)

	// 根据消息协议的规定进行对应的操作
	switch msg.Operation {
	case Operations.Stdin:
		// TODO 历史命令
		// TODO 20201229 兼容性测试
		return copy(p, cast.ToString(msg.Data)), nil
	case Operations.Resize:
		logrus.Trace("重新调整窗口大小", msg.Cols, msg.Rows)
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case Operations.Ping:
		return 0, nil
	default:
		logrus.Infof("unknown message type '%s'", msg.Operation)
		return copy(p, EndOfTransmission), fmt.Errorf("unknown message type '%s'", msg.Operation)
	}
}

// Write 将参数p中的数据拷贝到 websocket 连接中
func (t *ContainerSessionAdapter) Write(p []byte) (int, error) {
	tm := WebsocketMessage{
		Operation: Operations.Stdout,
		Data:      bytesconv.BytesToString(p),
	}
	msg, err := json.Marshal(tm)
	if err != nil {
		logrus.Tracef("write parse message err: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		logrus.Tracef("write message err: %v", err)
		return 0, err
	}
	logrus.Trace("输出数据", tm.Data)
	return len(p), nil
}
