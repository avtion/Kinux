package services

import (
	"context"
	"errors"
	"sync"
)

/* Websocket调度器中心 */
var scheduleCenter = new(sync.Map)

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

type ScheduleCenterListResult struct {
	ID          int    `json:"id"` // 用户ID
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	IsPty       bool   `json:"is_pty"`
	PtyMetaData string `json:"pty_meta_data"`
}

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
					return ws.pty.metaData
				}
				return ""
			}(),
		})
		return true
	})
	return
}
