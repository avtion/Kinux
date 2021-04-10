package services

import (
	"Kinux/core/web/msg"
	"Kinux/tools/bytesconv"
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io"
	"sync"
)

func init() {
	RegisterWebsocketOperation(wsOpAttachOtherWsWriter, AttachTargetWs)
	RegisterWebsocketOperation(wsOpStopAttachOtherWsWriter, StopAttachTargetWs)
}

/* Websocket调度器中心 */
var scheduleCenter = new(sync.Map)

// TargetWebsocketNotExistErr 错误
var TargetWebsocketNotExistErr = errors.New("目标Websocket链接不存在")

// RegisterWsConn 注册websocket链接
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
	oldWs.(*WebsocketSchedule).Context.Done()
	oldWs.(*WebsocketSchedule).SayGoodbyeToPty()
	_ = oldWs.(*WebsocketSchedule).Close()
	scheduleCenter.Store(id, ws)
	return
}

// ScheduleCenterListResult Websocket链接信息
type ScheduleCenterListResult struct {
	ID          int    `json:"id"` // 用户ID
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	IsPty       bool   `json:"is_pty"`
	PtyMetaData string `json:"pty_meta_data"`

	Lesson  uint `json:"lesson"`
	Mission uint `json:"mission"`
	Exam    uint `json:"exam"`
}

// ListScheduleCenterInfo 获取Websocket链接信息
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

			Lesson: func() uint {
				if ws.pty != nil {
					return ws.pty.metaData.GetLessonID()
				}
				return 0
			}(),
			Mission: func() uint {
				if ws.pty != nil {
					return ws.pty.metaData.GetMissionID()
				}
				return 0
			}(),
			Exam: func() uint {
				if ws.pty != nil {
					return ws.pty.metaData.GetExamID()
				}
				return 0
			}(),
		})
		return true
	})
	return
}

// SendMessageToTargetWs 向目标Websocket链接发送信息
func SendMessageToTargetWs(_ context.Context, id int, text string) (err error) {
	_ws, isExist := scheduleCenter.Load(id)
	if !isExist {
		return TargetWebsocketNotExistErr
	}
	ws, _ := _ws.(*WebsocketSchedule)
	return ws.SendMsg(msg.BuildSuccess(text))
}

// AttachTargetWs 侵入目标Websocket链接的终端
func AttachTargetWs(ws *WebsocketSchedule, any jsoniter.Any) (err error) {
	// 解析数据
	params := &struct {
		TargetID string `json:"target_id"`
	}{}
	any.Get("data").ToVal(params)
	var target = cast.ToInt(params.TargetID)
	if target == 0 {
		return TargetWebsocketNotExistErr
	}

	// 获取目标Websocket链接
	_ws, isExist := scheduleCenter.Load(target)
	if !isExist {
		return TargetWebsocketNotExistErr
	}
	targetWs := _ws.(*WebsocketSchedule)

	if targetWs.pty == nil {
		return errors.New("目标Websocket链接无正在使用的终端")
	}

	// 初始化Pty装饰器
	ws.pty = ws.InitPtyWrapper(func(w *WsPtyWrapper) {
		w.metaData = targetWs.pty.metaData
	})
	if targetWs.pty.stdoutListener != nil {
		targetWs.pty.stdoutListener = io.MultiWriter(targetWs.pty.stdoutListener, ws.pty)
	} else {
		targetWs.pty.stdoutListener = ws.pty
	}
	_, _ = ws.pty.Write(bytesconv.StringToBytes("连接成功"))
	return
}

// StopAttachTargetWs 侵入目标Websocket链接的终端
func StopAttachTargetWs(ws *WebsocketSchedule, _ jsoniter.Any) (err error) {
	if ws.pty == nil {
		return
	}
	// 直接抛弃掉Pty
	_ = ws.pty.Close()
	ws.pty.ws = nil
	ws.pty = nil
	return
}

// ForceTargetLogout 强制目标Websocket链接登出
func ForceTargetLogout(_ context.Context, id int) (err error) {
	_ws, isExist := scheduleCenter.Load(id)
	if !isExist {
		return TargetWebsocketNotExistErr
	}
	ws, _ := _ws.(*WebsocketSchedule)
	ws.RequireClientAuth()
	return
}

// BroadcastMessage 广播消息
func BroadcastMessage(_ context.Context, text string) (err error) {
	var _res = msg.BuildSuccess(text)
	scheduleCenter.Range(func(key, value interface{}) bool {
		ws, _ := value.(*WebsocketSchedule)
		if err := ws.SendMsg(_res); err != nil {
			logrus.WithField("err", err).Error("广播消息失败")
		}
		return true
	})
	return nil
}
