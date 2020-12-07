package services

// 终止标识符EOT
const EndOfTransmission = "\u0004"

type Operation = string

// TerminalMessage 消息协议中的规定的操作
var Operations = struct {
	Stdin  Operation // 输入
	Stdout Operation // 输出
	Resize Operation // 重新调整窗体大小
	Ping   Operation // 保持活跃
}{
	Stdin:  "stdin",
	Stdout: "stdout",
	Resize: "resize",
	Ping:   "ping",
}

// TerminalMessage 终端之间的消息协议
type TerminalMessage struct {
	Operation Operation `json:"operation"`
	Data      string    `json:"data"`
	// 窗体大小 - 对应 remotecommand.TerminalSize 结构体
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}
