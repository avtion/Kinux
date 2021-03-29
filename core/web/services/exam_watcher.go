package services

import (
	"Kinux/core/web/models"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

// 考试监控者
type ExamWatcher struct {
}

// 用于守护暴毙的监考对象
func NewExamWatcherDaemon() {
	// TODO 从数据库中查询未结束的考试日志
	// TODO 重新拉起监控者
}

// 考试监考者
func NewExamWatcher(ctx context.Context, ac, examID uint) (err error) {
	// TODO 查询目标用户和对应的考试是否有效（未超时）
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

	// 剩余时间
	restTime := exam.TimeLimit

	// 如果TickAt不为零值，则认为该用户在进行考试过程中发生了中断
	if !eLog.TickAt.IsZero() {
		// 那么剩余时间从TickAt开始计算
		passT := eLog.TickAt.Sub(eLog.CreatedAt)
		restTime -= passT
	}

	finishC := time.After(restTime)
	timer := time.NewTimer(5 * time.Minute)
	go func() {
		ctx = context.Background()
		defer func() {
			timer.Stop()
		}()
		for {
			select {
			case <-timer.C:
				// 定时脉冲用于记录TickAt
				if _err := models.GetGlobalDB().WithContext(ctx).Model(eLog).Update(
					"tick_at", time.Now()).Error; _err != nil {
					logrus.WithFields(logrus.Fields{
						"ac":     ac,
						"examID": examID,
					}).Error(err)
					// 脉冲更新失败不返回
				}
			case <-finishC:
				// 结束时间为考试结束时间
				if _err := models.GetGlobalDB().WithContext(ctx).Model(eLog).Update("end_at", time.Now()).Error; _err != nil {
					logrus.WithFields(logrus.Fields{
						"ac":     ac,
						"examID": examID,
					}).Error(err)
					return
				}
				return
			}
		}
	}()

	go func() {
		// TODO 告诉考生还有多久考试结束
	}()
	return
}
