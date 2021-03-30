// WebSocket业务层
package services

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/tools/bytesconv"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"net/http"
	"time"
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
	context.Context

	// 消息发送通道
	dataSenderCh chan []byte

	// 操作中间件
	operationMiddlewares []WsOperationHandler

	// 关闭通知
	daemonStopCh chan struct{}

	// 用户JWT密钥
	userToken *jwt.Token
	Account   *models.Account

	// SSH终端相关
	pty *WsPtyWrapper

	// 时间标记
	CreatedAt time.Time
}

type WsFn func(ws *WebsocketSchedule) (err error)

// 创建websocket调度器
func NewWebsocketSchedule(c *gin.Context, fns ...WsFn) (ws *WebsocketSchedule, err error) {
	wsConn, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	ws = &WebsocketSchedule{
		Conn:                 wsConn,
		Context:              ctx,
		daemonStopCh:         make(chan struct{}),
		operationMiddlewares: []WsOperationHandler{authMiddleware},
		dataSenderCh:         make(chan []byte, 1<<5),
		CreatedAt:            time.Now(),
	}
	if err = ws.Apply(fns...); err != nil {
		cancel()
		return
	}

	wsConn.SetCloseHandler(func(code int, text string) error {
		logrus.Debug("wsConn关闭并执行closeHandler")
		// websocket原本的方法
		message := websocket.FormatCloseMessage(code, "")
		_ = wsConn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))

		// 关闭监听者
		ws.StopDaemon()
		return nil
	})

	// 启动守护协程
	go ws.daemon(cancel)

	// 启动数据发送协程
	go ws.dataSender()

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
func (ws *WebsocketSchedule) daemon(ctxCancelFn context.CancelFunc) {
	l := logrus.WithField("module", "websocket守护协程")

	// FIX 修复并发控制问题
	// 当守护协程终止即该Websocket调度器旗下所有协程终结
	defer ctxCancelFn()

	for {
		select {
		case <-ws.daemonStopCh:
			l.Trace("接收到主动关闭消息")

			// 埋点 - 终止Pty终端
			ws.SayGoodbyeToPty()
			return
		default:
			_, message, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway) {
					l.WithField("err", err).Debug("客户端结束")
					continue
				}
				l.WithField("err", err).Debug("获取客户端数据时发生错误")
				return
			}

			// 解析数据
			any := jsoniter.Get(message)

			// 应用中间件
			if err = func() error {
				for _, fn := range ws.operationMiddlewares {
					if _err := fn(ws, any); _err != nil {
						return _err
					}
				}
				return nil
			}(); err != nil {
				continue
			}

			// 执行对应的操作
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

// 数据发送器
func (ws *WebsocketSchedule) dataSender() {
	// 修复并发写会导致panic https://github.com/gorilla/websocket/issues/380
	for {
		select {
		case <-ws.Done():
			return
		case data := <-ws.dataSenderCh:
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				logrus.WithField("err", err).Error("消息发送失败")
				return
			}
		}
	}
}

// 向客户端发送数据
func (ws *WebsocketSchedule) SendData(p []byte) {
	select {
	case <-ws.Done():
		return
	case ws.dataSenderCh <- p:
	}
}

// 发送终止信号给Pty
func (ws *WebsocketSchedule) SayGoodbyeToPty() {
	if ws.pty == nil {
		return
	}
	_, _ = ws.pty.writer.Write([]byte(EndOfTransmission))
	ws.pty = nil
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

// Websocket链接相关操作
const (
	_                   wsOperation = iota
	wsOpNewPty                      // 用于创建终端链接，由 Mission 负责实现
	wsOpStdin                       // 用于终端的输入
	wsOpStdout                      // 用于终端的输出
	wsOpResize                      // 用于终端重新调整窗体大小
	wsOpMsg                         // 服务端向客户端发送通知
	wsOpMissionApply                // 客户端发起Mission
	wsOpAuth                        // 客户端向服务端发起鉴权
	wsOpRequireAuth                 // 服务端要求客户端进行鉴权
	wsOpRefreshToken                // 刷新密钥
	wsOpShutdownPty                 // 关闭终端链接（即向终端发送 EndOfTransmission）
	wsOpResetContainers             // 重置容器
	wsOpContainersDone              // 容器部署成功

	// 2021/03/30
	wsOpAttachOtherWsWriter // 侵入其他Websocket链接
	wsOpLeaveExam           // 退出考试
	wsOpExamRunning         // 考试进行中（用于主动告诉用户正在进行考试）
)

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

	ws.SendData(data)
	return
}

// websocket处理函数映射
var wsOperationsMapper = make(map[wsOperation]WsOperationHandler)

// 注册websocket链接操作
func RegisterWebsocketOperation(op wsOperation, handler WsOperationHandler) {
	wsOperationsMapper[op] = handler
}

// 向客户端发起鉴权请求
func (ws *WebsocketSchedule) RequireClientAuth() {
	data, _ := jsoniter.Marshal(&WebsocketMessage{
		Op:   wsOpRequireAuth,
		Data: "请重新登录",
	})
	ws.SendData(data)
	return
}

// 鉴权中间件
func authMiddleware(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	if op := any.Get("op").ToUint(); op != wsOpAuth && ws.userToken == nil {
		ws.RequireClientAuth()
		return errors.New("要求客户端进行鉴权")
	}
	return nil
}
