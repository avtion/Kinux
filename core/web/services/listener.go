package services

import (
	"Kinux/core/web/models"
	"context"
	terminalReader "github.com/Nerdmaster/terminal"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sirupsen/logrus"
	"io"
)

// 监听方式
type listenWay int8

const (
	_            listenWay = iota
	ListenStdin            // 输入监听器
	ListenStdout           // 输出监听器
	ListenWeb              // Web监听器
)

// 终端监听器
type TerminalListener struct {
	*terminalReader.Reader
}

// 创建新的监听器
func NewTerminalListener(r io.Reader) *TerminalListener {
	t := &TerminalListener{
		Reader: terminalReader.NewReader(r),
	}
	return t
}

// 给终端创建新的监听器
func NewPtyWrapperListenerOpt(way listenWay) (opt WsPtyWrapperOption, listener *TerminalListener) {
	reader, writer := io.Pipe()
	listener = NewTerminalListener(reader)
	return func(w *WsPtyWrapper) {
		switch way {
		case ListenStdin:
			if w.ws.PtyStdin == nil {
				w.ws.PtyStdin = writer
			} else {
				w.ws.PtyStdin = io.MultiWriter(w.ws.PtyStdin, writer)
			}
		case ListenStdout:
			if w.stdoutListener == nil {
				w.stdoutListener = writer
			} else {
				w.stdoutListener = io.MultiWriter(w.stdoutListener, writer)
			}
		}
	}, listener
}

// Debug输出
func (l *TerminalListener) DebugPrint(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			line, err := l.Reader.ReadLine()
			if err != nil {
				logrus.WithField("err", err).Error("监听器发生错误")
				return
			}
			logrus.WithField("line", line).Trace("监听器Debug输出")
		}
	}()
}

// 创建检查点监听器
func NewWrapperForCheckpointCallback(ctx context.Context, ac *models.Account, exam *models.Exam,
	mission *models.Mission, container string, checkpoints ...*models.Checkpoint) (opt WsPtyWrapperOption) {
	if mission == nil || len(checkpoints) == 0 {
		return
	}

	var (
		inOpt, outOpt       WsPtyWrapperOption
		inReader, outReader *TerminalListener
		inMap, outMap       map[string]func(ctx context.Context) (err error)
	)

	for _, cp := range checkpoints {
		switch cp.Method {
		case models.MethodExec:
			// 懒加载
			if inOpt == nil || inReader == nil {
				inOpt, inReader = NewPtyWrapperListenerOpt(ListenStdin)
			}
			if inMap == nil {
				inMap = make(map[string]func(ctx context.Context) (err error))
			}
			if exam == nil {
				inMap[cp.In] = models.NewMissionScoreCallback(ac.ID, mission.ID, cp.ID, container)
			} else {
				inMap[cp.In] = models.NewExamScoreCallback(ac.ID, exam.ID, mission.ID, cp.ID, container)
			}
		case models.MethodStdout:
			// 懒加载
			if outOpt == nil || outReader == nil {
				outOpt, outReader = NewPtyWrapperListenerOpt(ListenStdout)
			}
			if outMap == nil {
				outMap = make(map[string]func(ctx context.Context) (err error))
			}
			if exam == nil {
				outMap[cp.Out] = models.NewMissionScoreCallback(ac.ID, mission.ID, cp.ID, container)
			} else {
				outMap[cp.Out] = models.NewExamScoreCallback(ac.ID, exam.ID, mission.ID, cp.ID, container)
			}
		case models.MethodTargetPort:
			// TODO 完成Web端口监听
			continue
		default:
			continue
		}
	}

	// 监听方法
	listenFn := func(ctx context.Context, listener *TerminalListener,
		callbackMap map[string]func(ctx context.Context) (err error)) {
		// 比较器
		dmp := diffmatchpatch.New()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			line, _err := listener.Reader.ReadLine()
			if _err != nil {
				logrus.WithField("err", _err).Error("监听器发生错误")
				return
			}
			if callback, isExist := callbackMap[line]; isExist {
				if _err := callback(ctx); _err != nil {
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

	// 输入和输出监听
	// TODO Web监听
	for _, tmp := range []struct {
		l *TerminalListener
		m map[string]func(ctx context.Context) (err error)
	}{{l: inReader, m: inMap}, {l: outReader, m: outMap}} {
		if tmp.l != nil {
			go listenFn(ctx, tmp.l, tmp.m)
		}
	}

	return CombineWsPtyWrapperOptions(inOpt, outOpt)
}
