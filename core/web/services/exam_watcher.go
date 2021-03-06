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

// ExamWatchers 考试监控者 uint -> ExamWatcher
var ExamWatchers sync.Map

// ExamWatcher 监控实例
type ExamWatcher struct {
	ELog           *models.ExamLog
	ManualFinishCh chan struct{}
	RestTime       time.Duration
	StartedAt      time.Time
}

// NewExamWatcher 考试监考者
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
		"account = ? AND exam = ?", eLog.Account, eLog.Exam).First(eLog).Error
	if errors.Is(_err, gorm.ErrRecordNotFound) {
		if err = models.GetGlobalDB().WithContext(ctx).Create(eLog).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"ac":   ac,
				"exam": examID,
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
			ExamWatchers.Delete(ac)
			// 告诉用户退出考试界面
			leaveExam(ctx, eLog.Account)
		}()

		ctx = context.Background()
		finishT := time.NewTimer(ew.RestTime)
		tickerT := time.NewTicker(1 * time.Minute)
		defer func() {
			finishT.Stop()
			tickerT.Stop()
		}()

		// 检查用户Websocket是否在线并主动发送时间校验信息
		var sendExamInfoToAccountFn = func() {
			_ws, _isExist := scheduleCenter.Load(int(ac))
			if _isExist {
				ws := _ws.(*WebsocketSchedule)
				raw, __err := jsoniter.Marshal(&WebsocketMessage{
					Op:   wsOpExamRunning,
					Data: NewExamRunningInfo(ew),
				})
				if __err == nil {
					ws.SendData(raw)
				}
			}
		}
		sendExamInfoToAccountFn()

		for {
			select {
			case <-tickerT.C:
				// 检查用户Websocket是否在线并主动发送时间校验信息
				sendExamInfoToAccountFn()

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
	_ws, isExist := scheduleCenter.Load(int(ac))
	if !isExist {
		return
	}
	ws, _ := _ws.(*WebsocketSchedule)

	data, _ := jsoniter.Marshal(&WebsocketMessage{
		Op:   wsOpLeaveExam,
		Data: "考试结束",
	})
	ws.SendData(data)
}

// GetLeftTime 获取考试剩余时间
func GetLeftTime(_ context.Context, ac uint) (res time.Duration, err error) {
	// 如果
	_ew, isExist := ExamWatchers.Load(ac)
	if !isExist {
		return 0, errors.New("用户没有考试")
	}
	ew, _ := _ew.(*ExamWatcher)
	return ew.RestTime - time.Now().Sub(ew.StartedAt), nil
}

// GetExamInfo 获取考试的信息
func GetExamInfo(ctx context.Context, ac uint) (res *models.Exam, err error) {
	// 如果
	_ew, isExist := ExamWatchers.Load(ac)
	if !isExist {
		return nil, errors.New("用户没有考试")
	}
	ew, _ := _ew.(*ExamWatcher)
	return models.GetExam(ctx, ew.ELog.Exam)
}

// ExamRunningInfo 实验进行中信息
type ExamRunningInfo struct {
	Account  uint          `json:"account"`   // 用户ID
	ExamID   uint          `json:"exam_id"`   // 实验ID
	ExamName string        `json:"exam_name"` // 实验名称
	LeftTime time.Duration `json:"left_time"` // 剩余时间
}

// NewExamRunningInfo 创建实验进行中信息
func NewExamRunningInfo(ew *ExamWatcher) (res ExamRunningInfo) {
	exam, _ := models.GetExam(context.Background(), ew.ELog.Exam)
	leftTime := ew.RestTime - time.Now().Sub(ew.StartedAt)
	res = ExamRunningInfo{
		Account:  ew.ELog.Account,
		ExamID:   exam.ID,
		ExamName: exam.Name,
		LeftTime: leftTime,
	}
	return
}

// InitExamWatcher 初始化考试监考者
func InitExamWatcher(ctx context.Context) {
	defer logrus.Info("考试监考者初始完毕")
	// 从数据库查找没有结束时间的考试记录并恢复监考者
	var eLogs []*models.ExamLog
	if err := models.GetGlobalDB().WithContext(ctx).Model(new(models.ExamLog)).Where(
		"end_at = ?", time.Time{}).Find(&eLogs).Error; err != nil {
		logrus.WithField("err", err).Error("初始化考试监考者失败")
		return
	}
	if len(eLogs) == 0 {
		return
	}
	for _, v := range eLogs {
		if err := NewExamWatcher(ctx, v.Account, v.Exam); err != nil {
			logrus.WithField("err", err).Error("初始化考试监考者失败")
		}
	}
}
