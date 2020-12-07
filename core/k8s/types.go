package k8s

import (
	"io"
	"k8s.io/client-go/tools/remotecommand"
)

// PtyHandler 定义remoteCommand所需要的方法集
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue // 调整终端大小
	Done()                          // 终止
}
