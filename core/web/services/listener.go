package services

import (
	"context"
	terminalReader "github.com/Nerdmaster/terminal"
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
	context.Context
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
			w.writer = io.MultiWriter(w.writer, writer)
		case ListenStdout:
			if w.stdoutListener == nil {
				w.stdoutListener = writer
			} else {
				w.stdoutListener = io.MultiWriter(w.stdoutListener, writer)
			}
		}
		listener.Context = w.ChildCtx
	}, listener
}

// Debug输出
func (l *TerminalListener) DebugPrint() {
	if l.Context == nil {
		logrus.Error("非法Debug输出，Pty未初始化")
	}
	go func() {
		for {
			select {
			case <-l.Context.Done():
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
