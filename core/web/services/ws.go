// WebSocket业务层
package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/msg"
	"Kinux/tools/bytesconv"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"sync"
)

// 终止标识符EOT
const EndOfTransmission = "\u0004"

// webSocket默认升级器
var defaultUpgrader = &websocket.Upgrader{
	// TODO 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Websocket调度器
type WebsocketSchedule struct {
	*websocket.Conn
	*gin.Context
	isDebug bool

	daemonStopCh chan struct{}

	mutex sync.RWMutex

	// 用户JWT密钥
	userToken *jwt.Token

	// SSH终端相关
	sizeChan chan remotecommand.TerminalSize
}

type WsFn func(ws *WebsocketSchedule) (err error)

// 创建websocket调度器
func NewWebsocketSchedule(c *gin.Context, fns ...WsFn) (ws *WebsocketSchedule, err error) {
	wsConn, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	ws = &WebsocketSchedule{
		Conn:         wsConn,
		Context:      c,
		daemonStopCh: make(chan struct{}),
		sizeChan:     make(chan remotecommand.TerminalSize),
	}
	if err = ws.Apply(fns...); err != nil {
		return
	}

	// 启动守护协程
	go ws.daemon()

	return ws, nil
}

// 应用中间函数
func (ws *WebsocketSchedule) Apply(fns ...WsFn) (err error) {
	for _, fn := range fns {
		select {
		case <-ws.Done():
			return ws.Err()
		default:
			if err = fn(ws); err != nil {
				return
			}
		}
	}
	return
}

// 守护协程 - 用于读取数据并根据对应的操作进行分发
// 当该websocket链接交付给容器进行交互时，守护协程应当结束监听
func (ws *WebsocketSchedule) daemon() {
	l := logrus.WithField("module", "websocket守护协程")
	for {
		select {
		case <-ws.Done():
			l.Trace("上下文结束")
			return
		case <-ws.daemonStopCh:
			l.Trace("接收到主动关闭消息")
			return
		default:
			_, message, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway) {
					l.WithField("err", err).Debug("客户端结束")
				}
				l.WithField("err", err).Debug("获取客户端数据时发生错误")
				return
			}
			any := jsoniter.Get(message)
			op := any.Get("op").ToUint()
			if op == 0 {
				continue
			}
			fn, isExist := wsOperationsMapper[op]
			if !isExist {
				l.WithField("raw", bytesconv.BytesToString(message)).Debug(
					"无法识别客户端发送的请求")
				continue
			}
			if err = fn(ws, any); err != nil {
				l.WithField("err", err).Error("解析数据发生错误")
				_ = ws.SendMsg(msg.BuildFailed(err))
				continue
			}
		}
	}
}

// 结束守护协程
func (ws *WebsocketSchedule) StopDaemon() {
	close(ws.daemonStopCh)
}

/*
	定义通信接口
*/
type wsOperation = uint

// WebsocketMessage C/S消息协议
type WebsocketMessage struct {
	Op   wsOperation `json:"op"`
	Data interface{} `json:"data"`
}

const (
	_          wsOperation = iota
	wsOpNewPty             // 用于创建终端链接，由 Mission 负责实现
	wsOpStdin              // 用于终端的输入
	wsOpStdout             // 用于终端的输出
	wsOpResize             // 用于终端重新调整窗体大小

	wsOpMsg           // 服务端向客户端发送通知
	wsOpResourceApply // 客户端资源申请

	wsOpAuth // 客户端向服务端发起鉴权
)

// 用于终端的websocket装饰器
type WsPtyWrapper struct {
	ws *WebsocketSchedule
}

var _ k8s.PtyHandler = (*WsPtyWrapper)(nil)

func (ws *WebsocketSchedule) InitPtyWrapper() *WsPtyWrapper {
	return &WsPtyWrapper{
		ws: ws,
	}
}

// 对调度器进行封装用于适配终端场景
func (pw *WsPtyWrapper) Read(p []byte) (n int, err error) {
	_, message, err := pw.ws.ReadMessage()
	if err != nil {
		logrus.WithField("err", err).Debug("获取客户端数据时发生错误")
		return
	}
	any := jsoniter.Get(message)
	op := any.Get("op").ToUint()
	switch op {
	case wsOpStdin:
		// 进行写入操作
		return copy(p, bytesconv.StringToBytes(any.Get("data").ToString())), nil
	default:
		// 对于非终端指令兼容
		fn, isExist := wsOperationsMapper[op]
		if !isExist {
			logrus.WithField("raw", bytesconv.BytesToString(message)).Debug(
				"websocket无法识别客户端发送的请求")
			return copy(p, EndOfTransmission), nil
		}
		if err = fn(pw.ws, any); err != nil {
			logrus.WithField("err", err).Error("websocket解析数据失败")
			return copy(p, EndOfTransmission), err
		}
	}
	return
}

// 对调度器进行封用于装适配终端
func (pw *WsPtyWrapper) Write(p []byte) (n int, err error) {
	raw, err := jsoniter.Marshal(&WebsocketMessage{
		Op:   wsOpStdout,
		Data: bytesconv.BytesToString(p),
	})
	if err != nil {
		return 0, err
	}
	if err = pw.ws.WriteMessage(websocket.TextMessage, raw); err != nil {
		return 0, err
	}
	return len(p), nil
}

// 实现 remotecommand.TerminalSizeQueue 接口
func (pw *WsPtyWrapper) Next() *remotecommand.TerminalSize {
	select {
	case size := <-pw.ws.sizeChan:
		return &size
	case <-pw.ws.Done():
		return nil
	}
}

// 实现 k8s.PtyHandler 接口
func (pw *WsPtyWrapper) Done() {
	return
}

/*
	Websocket链接相关操作
*/

// ws处理函数 - any指向未解析的原数据
type WsOperationHandler func(ws *WebsocketSchedule, any jsoniter.Any) (err error)

// 给客户端发送消息 - 利用原有的消息构建框架
func (ws *WebsocketSchedule) SendMsg(result msg.Result) (err error) {
	data, err := jsoniter.Marshal(&WebsocketMessage{
		Op:   wsOpMsg,
		Data: result,
	})
	if err != nil {
		return
	}
	return ws.WriteMessage(websocket.TextMessage, data)
}

var wsOperationsMapper = map[wsOperation]WsOperationHandler{
	wsOpResourceApply: func(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
		// TODO 完成资源申请接口的实现
		return errors.New("未实现")
	},
}

// 注册websocket链接操作
func RegisterWebsocketOperation(op wsOperation, handler WsOperationHandler) {
	wsOperationsMapper[op] = handler
}
