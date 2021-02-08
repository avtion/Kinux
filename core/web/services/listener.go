package services

import (
	terminalReader "github.com/Nerdmaster/terminal"
	"github.com/sirupsen/logrus"
	"io"
)

// 监听方式
type listenWay int8

const (
	_ listenWay = iota
	ListenStdin
	ListenStdout
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
func (l *TerminalListener) DebugPrint() {
	go func() {
		for {
			line, err := l.Reader.ReadLine()
			if err != nil {
				logrus.WithField("err", err).Error("监听器发生错误")
				return
			}
			logrus.WithField("line", line).Trace("监听器Debug输出")
		}
	}()
}
