package services

import (
	"Kinux/core/web/msg"
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

func init() {
	RegisterWebsocketOperation(wsOpAttachOtherWsWriter, AttachTargetWs)
}

/* Websocket调度器中心 */
var scheduleCenter = new(sync.Map)

// 错误
var TargetWebsocketNotExistErr = errors.New("目标Websocket链接不存在")

// 注册websocket链接
func RegisterWsConn(id int, ws *WebsocketSchedule) (err error) {
	if id == 0 {
		err = errors.New("非法用户，链接无法挂载调度中心")
		return
	}
	oldWs, isExist := scheduleCenter.LoadOrStore(id, ws)

	// 如果不存在则直接结束调用
	if !isExist {
		return
	}

	// 如果存在则抛弃原本的链接并保存
	oldWs.(*WebsocketSchedule).SayGoodbyeToPty()
	_ = oldWs.(*WebsocketSchedule).Close()
	scheduleCenter.Store(id, ws)
	return
}

// Websocket链接信息
type ScheduleCenterListResult struct {
	ID          int    `json:"id"` // 用户ID
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	IsPty       bool   `json:"is_pty"`
	PtyMetaData string `json:"pty_meta_data"`
}

// 获取Websocket链接信息
func ListScheduleCenterInfo(_ context.Context) (res []*ScheduleCenterListResult) {
	res = make([]*ScheduleCenterListResult, 0)
	scheduleCenter.Range(func(key, value interface{}) bool {
		ws := value.(*WebsocketSchedule)
		res = append(res, &ScheduleCenterListResult{
			ID:        key.(int),
			Username:  ws.Account.Username,
			CreatedAt: ws.CreatedAt.Format("2006-01-02 15:04:05"),
			IsPty:     ws.pty != nil,
			PtyMetaData: func() string {
				if ws.pty != nil {
					return ws.pty.metaData.StrFormat()
				}
				return "无终端活动"
			}(),
		})
		return true
	})
	return
}

// 向目标Websocket链接发送信息
func SendMessageToTargetWs(_ context.Context, id int, text string) (err error) {
	_ws, isExist := scheduleCenter.Load(id)
	if !isExist {
		return TargetWebsocketNotExistErr
	}
	ws, _ := _ws.(*WebsocketSchedule)
	return ws.SendMsg(msg.BuildSuccess(text))
}

// 侵入目标Websocket链接的终端
func AttachTargetWs(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	// 解析数据
	params := &struct {
		TargetID int `json:"target_id"`
	}{}
	any.Get("data").ToVal(params)

	// 获取目标Websocket链接
	_ws, isExist := scheduleCenter.Load(params.TargetID)
	if !isExist {
		return TargetWebsocketNotExistErr
	}
	targetWs := _ws.(*WebsocketSchedule)

	if targetWs.pty == nil {
		return errors.New("目标Websocket链接无正在使用的终端")
	}

	// 初始化Pty装饰器
	targetWs.pty.stdoutListener = ws.InitPtyWrapper(func(w *WsPtyWrapper) {
		w.metaData = targetWs.pty.metaData
	})
	return
}
