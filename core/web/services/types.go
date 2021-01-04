package services

// 终止标识符EOT
const EndOfTransmission = "\u0004"

type Operation = string

// WebsocketMessage 消息协议中的规定的操作
var Operations = struct {
	Stdin         Operation // 用于终端的输入
	Stdout        Operation // 用于终端的输出
	Resize        Operation // 用于终端重新调整窗体大小
	Ping          Operation // 保持ws活跃
	Msg           Operation // 向客户端发送通知
	ResourceApply Operation // 客户端资源申请
}{
	Stdin:         "stdin",
	Stdout:        "stdout",
	Resize:        "resize",
	Ping:          "ping",
	Msg:           "msg",
	ResourceApply: "apply",
}

// WebsocketMessage C/S消息协议
type WebsocketMessage struct {
	Operation Operation   `json:"operation"`
	Data      interface{} `json:"data"`
	// 窗体大小 - 对应 remotecommand.TerminalSize 结构体
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}
