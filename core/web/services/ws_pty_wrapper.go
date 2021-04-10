package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"Kinux/tools/bytesconv"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"k8s.io/client-go/tools/remotecommand"
)

// WsPtyWrapper 用于终端的websocket装饰器
type WsPtyWrapper struct {
	ws       *WebsocketSchedule
	reader   io.Reader
	writer   io.Writer
	sizeChan chan remotecommand.TerminalSize

	// 输出监听者 - 输入监听者需要在调度器中进行注入
	stdoutListener io.Writer

	// 并发控制
	ChildCtx context.Context
	cancelFn context.CancelFunc

	// 额外参数
	metaData PtyMeta // 用于展示当前终端正在执行的任务信息
}

type WsPtyWrapperOption = func(w *WsPtyWrapper)

var _ k8s.PtyHandler = (*WsPtyWrapper)(nil)

// InitPtyWrapper 初始化终端装饰器
func (ws *WebsocketSchedule) InitPtyWrapper(opts ...WsPtyWrapperOption) *WsPtyWrapper {
	// 埋点 - 终止上一个Pty终端
	ws.SayGoodbyeToPty()

	// 使用io管道对输入的数据进行拷贝
	r, w := io.Pipe()

	// FIX 初始化并发控制
	childCtx, cancel := context.WithCancel(ws.Context)
	wrapper := &WsPtyWrapper{
		ws:       ws,
		reader:   r,
		writer:   w,
		sizeChan: make(chan remotecommand.TerminalSize),
		ChildCtx: childCtx,
		cancelFn: cancel,
	}
	defer func() {
		ws.pty = wrapper
	}()

	for _, opt := range opts {
		if opt != nil {
			opt(wrapper)
		}
	}
	return wrapper
}

// 对调度器进行封装用于适配终端场景
func (pw *WsPtyWrapper) Read(p []byte) (n int, err error) {
	select {
	case <-pw.ChildCtx.Done():
		return 0, io.EOF
	default:
	}
	return pw.reader.Read(p)
}

// 对调度器进行封用于装适配终端
func (pw *WsPtyWrapper) Write(p []byte) (n int, err error) {
	select {
	case <-pw.ChildCtx.Done():
		return len(p), nil
	default:
	}
	// 监听器输出
	if pw.stdoutListener != nil {
		_, _ = pw.stdoutListener.Write(p)
	}
	raw, _ := jsoniter.Marshal(&WebsocketMessage{
		Op:   wsOpStdout,
		Data: bytesconv.BytesToString(p),
	})
	if pw.ws == nil {
		return len(p), nil
	}

	pw.ws.SendData(raw)
	return len(p), nil
}

// Next 实现 remotecommand.TerminalSizeQueue 接口
func (pw *WsPtyWrapper) Next() *remotecommand.TerminalSize {
	select {
	case size := <-pw.sizeChan:
		return &size
	case <-pw.ChildCtx.Done():
		return nil
	}
}

func (pw *WsPtyWrapper) Close() (err error) {
	pw.cancelFn()
	//pw.ws.SayGoodbyeToPty()
	// FIX 修复容器延迟关闭导致的当前PTY关闭问题
	//pw.ws.pty = nil
	return nil
}

// CombineWsPtyWrapperOptions 组合多个Pty装饰器
func CombineWsPtyWrapperOptions(wrappers ...WsPtyWrapperOption) WsPtyWrapperOption {
	return func(w *WsPtyWrapper) {
		for _, fn := range wrappers {
			if fn != nil {
				fn(w)
			}
		}
	}
}

// SetWsPtyMetaDataOption 设置pty的元数据
func SetWsPtyMetaDataOption(metaData PtyMeta) WsPtyWrapperOption {
	return func(w *WsPtyWrapper) {
		w.metaData = metaData
	}
}

// PtyMeta 元数据
type PtyMeta interface {
	GetType() MetaType
	StrFormat() string
	GetLessonID() uint
	GetMissionID() uint
	GetExamID() uint
}

type MetaType = uint

const (
	MissionMetaType MetaType = iota
	ExamMetaType
)

// MissionMeta 实验元数据
type MissionMeta struct {
	Account   *models.Account
	Lesson    *models.Lesson
	Mission   *models.Mission
	Container string
}

var _ PtyMeta = (*MissionMeta)(nil)

func (mm *MissionMeta) GetType() MetaType {
	return MissionMetaType
}

func (mm *MissionMeta) StrFormat() string {
	return fmt.Sprintf("实验: %s(%d)", mm.Mission.Name, mm.Mission.ID)
}

func (mm *MissionMeta) GetLessonID() uint {
	return mm.Lesson.ID
}

func (mm *MissionMeta) GetMissionID() uint {
	return mm.Mission.ID
}
func (mm *MissionMeta) GetExamID() uint {
	return 0
}

// NewMissionMeta 创建新的实验元数据
func NewMissionMeta(ac *models.Account, lesson *models.Lesson, ms *models.Mission, container string) PtyMeta {
	return &MissionMeta{
		Account:   ac,
		Lesson:    lesson,
		Mission:   ms,
		Container: container,
	}
}

// ExamMeta 考试元数据
type ExamMeta struct {
	Account   *models.Account
	Lesson    *models.Lesson
	Exam      *models.Exam
	Mission   *models.Mission
	Container string
}

var _ PtyMeta = (*ExamMeta)(nil)

func (em *ExamMeta) GetType() MetaType {
	return ExamMetaType
}

func (em *ExamMeta) StrFormat() string {
	return fmt.Sprintf("考试: %s-%s(%d)", em.Exam.Name, em.Mission.Name, em.Mission.ID)
}

func (em *ExamMeta) GetLessonID() uint {
	return em.Lesson.ID
}

func (em *ExamMeta) GetMissionID() uint {
	return em.Mission.ID
}
func (em *ExamMeta) GetExamID() uint {
	return em.Exam.ID
}

// NewExamMeta 创建新的考试元数据
func NewExamMeta(ac *models.Account, lesson *models.Lesson, ms *models.Mission, ex *models.Exam, container string) PtyMeta {
	return &ExamMeta{
		Account:   ac,
		Lesson:    lesson,
		Exam:      ex,
		Mission:   ms,
		Container: container,
	}
}

// NewMeta 创建原数据标签
func NewMeta(ac *models.Account, lesson *models.Lesson,
	ms *models.Mission, ex *models.Exam, container string) PtyMeta {
	if ex == nil {
		return NewMissionMeta(ac, lesson, ms, container)
	} else {
		return NewExamMeta(ac, lesson, ms, ex, container)
	}
}
