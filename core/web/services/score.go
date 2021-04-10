package services

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// NewScoreListener 成绩监听器
func NewScoreListener(account *models.Account, lesson *models.Lesson, exam *models.Exam, mc *models.Mission,
	container string, cps ...*models.Checkpoint) (opt WsPtyWrapperOption) {
	// 校验参数
	if account == nil || account.ID == 0 {
		return nil
	}
	if lesson == nil || lesson.ID == 0 {
		return nil
	}
	if mc == nil || mc.ID == 0 {
		return nil
	}
	if container == "" {
		return nil
	}
	if len(cps) == 0 {
		return nil
	}
	var examID uint
	if exam != nil && exam.ID != 0 {
		examID = exam.ID
	}

	// 定义监听器
	var (
		inOpt, outOpt       WsPtyWrapperOption
		inReader, outReader *TerminalListener
		inMap, outMap       map[string]func(w *WsPtyWrapper) (err error)
	)

	// 回调函数
	newCallbackFn := func(cp *models.Checkpoint) func(w *WsPtyWrapper) (err error) {
		return func(w *WsPtyWrapper) (err error) {
			if err = models.GetGlobalDB().WithContext(w.ChildCtx).Create(&models.Score{
				Model:      gorm.Model{},
				Account:    account.ID,
				Lesson:     lesson.ID,
				Exam:       examID,
				Mission:    mc.ID,
				Checkpoint: cp.ID,
				Container:  container,
			}).Error; err != nil {
				_ = w.ws.SendMsg(msg.BuildFailed(fmt.Sprintf("考点检验失败: %s", err)))
				return
			}
			return w.ws.SendMsg(msg.BuildSuccess(fmt.Sprintf("考点已完成: %s", cp.Name)))
		}
	}

	for _, cp := range cps {
		switch cp.Method {
		case models.MethodExec:
			// 初始化
			if inOpt == nil || inReader == nil {
				inOpt, inReader = NewPtyWrapperListenerOpt(ListenStdin)
			}
			if inMap == nil {
				inMap = make(map[string]func(w *WsPtyWrapper) (err error))
			}
			inMap[cp.In] = newCallbackFn(cp)
		case models.MethodStdout:
			// 初始化
			if outOpt == nil || outReader == nil {
				outOpt, outReader = NewPtyWrapperListenerOpt(ListenStdin)
			}
			if outMap == nil {
				outMap = make(map[string]func(w *WsPtyWrapper) (err error))
			}
			outMap[cp.Out] = newCallbackFn(cp)
		case models.MethodTargetPort:
			// TODO 支持
		}
	}

	// 监听方法
	listenFn := func(listener *TerminalListener, w *WsPtyWrapper,
		callbackMap map[string]func(w *WsPtyWrapper) (err error)) {
		// 比较器
		dmp := diffmatchpatch.New()
		for {
			select {
			case <-listener.Context.Done():
				return
			default:
			}
			line, _err := listener.Reader.ReadLine()
			if _err != nil {
				logrus.WithField("err", _err).Error("监听器发生错误")
				return
			}
			if callback, isExist := callbackMap[line]; isExist {
				if _err := callback(w); _err != nil {
					logrus.WithField("err", _err).Error("监听器执行回调函数失败")
					continue
				}
				delete(callbackMap, line)
			} else {
				for k := range callbackMap {
					diffs := dmp.DiffMain(line, k, false)
					logrus.WithField("diff", dmp.DiffPrettyText(diffs)).Trace("指令差异")
				}
			}
		}
	}

	return CombineWsPtyWrapperOptions(inOpt, outOpt, func(w *WsPtyWrapper) {
		// 输入和输出监听
		for _, tmp := range []struct {
			l *TerminalListener
			m map[string]func(w *WsPtyWrapper) (err error)
		}{{l: inReader, m: inMap}, {l: outReader, m: outMap}} {
			if tmp.l != nil {
				go listenFn(tmp.l, w, tmp.m)
			}
		}
	})
}
