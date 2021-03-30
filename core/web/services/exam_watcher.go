package services

import (
	"Kinux/core/web/models"
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

// 考试监控者 uint -> ExamWatcher
var ExamWatchers sync.Map

// 监控实例
type ExamWatcher struct {
	ELog           *models.ExamLog
	ManualFinishCh chan struct{}
	RestTime       time.Duration
	StartedAt      time.Time
}

// 考试监考者
func NewExamWatcher(ctx context.Context, ac, examID uint) (err error) {
	eWatcher, isExist := ExamWatchers.LoadAndDelete(ac)
	if isExist {
		// 如果存在旧监控者，则终止那场考试
		close(eWatcher.(*ExamWatcher).ManualFinishCh)
	}

	// 查询目标用户和对应的考试是否有效（未超时）
	eLog := &models.ExamLog{
		Model:   gorm.Model{},
		Account: ac,
		Exam:    examID,
		TickAt:  time.Time{},
		EndAt:   time.Time{},
	}
	_err := models.GetGlobalDB().WithContext(ctx).Where(
		"account = ? AND examID = ?", eLog.Account, eLog.Exam).First(eLog).Error
	if errors.Is(_err, gorm.ErrRecordNotFound) {
		if err = models.GetGlobalDB().WithContext(ctx).Create(eLog).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"ac":     ac,
				"examID": examID,
			}).Error(err)
			return
		}
	}

	// 如果考试结束时间存在则认为当前用户该考试已经结束
	if !eLog.EndAt.IsZero() {
		return errors.New("考试已经结束")
	}

	// 查询考试
	exam, err := models.GetExam(ctx, eLog.Exam)
	if err != nil {
		return err
	}

	// 主动结束通道
	manualFinishCh := make(chan struct{})
	ew := &ExamWatcher{
		ELog:           eLog,
		ManualFinishCh: manualFinishCh,
		StartedAt:      eLog.CreatedAt,
		RestTime:       exam.TimeLimit,
	}

	ExamWatchers.Store(ac, ew)

	// 如果TickAt不为零值，则认为该用户在进行考试过程中发生了中断
	if !eLog.TickAt.IsZero() {
		// 那么剩余时间从TickAt开始计算
		passT := eLog.TickAt.Sub(eLog.CreatedAt)
		ew.StartedAt = eLog.TickAt
		ew.RestTime -= passT
	}

	go func(mFinishCh chan struct{}) {
		defer func() {
			// 写入EndAt
			finishExam(ctx, eLog.ID)
			ExamWatchers.Delete(eLog.ID)
			// 告诉用户退出考试界面
			leaveExam(ctx, eLog.Account)
		}()

		ctx = context.Background()
		finishT := time.NewTimer(ew.RestTime)
		tickerT := time.NewTicker(5 * time.Minute)
		defer func() {
			finishT.Stop()
			tickerT.Stop()
		}()

		for {
			select {
			case <-tickerT.C:
				// TODO 检查用户Websocket是否在线并主动发送时间校验信息

				// 定时脉冲用于记录TickAt
				if _err := models.GetGlobalDB().WithContext(ctx).Model(eLog).Update(
					"tick_at", time.Now()).Error; _err != nil {
					logrus.WithFields(logrus.Fields{
						"ac":     ac,
						"examID": examID,
					}).Error(err)
					// 脉冲更新失败不返回
				}
			case <-finishT.C:
				// 结束时间为考试结束时间
				return
			case <-mFinishCh:
				// 主动结束考试
				return
			}
		}
	}(manualFinishCh)
	return
}

// 结束考试
func finishExam(ctx context.Context, eLogID uint) {
	if err := models.GetGlobalDB().WithContext(ctx).Model(&models.ExamLog{
		Model: gorm.Model{
			ID: eLogID,
		},
	}).Update("end_at", time.Now()).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"eLogID": eLogID,
		}).Error(err)
		return
	}
	return
}

// 退出考试
func leaveExam(_ context.Context, ac uint) {
	_ws, isExist := scheduleCenter.Load(ac)
	if !isExist {
		return
	}
	ws, _ := _ws.(*WebsocketSchedule)

	data, _ := jsoniter.Marshal(&WebsocketMessage{
		Op:   wsOpLeaveExam,
		Data: "退出考试",
	})
	ws.SendData(data)
}
