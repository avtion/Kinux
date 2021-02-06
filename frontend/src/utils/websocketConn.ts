import { Terminal } from 'xterm'
import { BaseResponse } from '@api/request'
import { notification } from 'ant-design-vue'

// 后端默认路由
export const DefaultBackendWebsocketRoute = 'ws://127.0.0.1:9001/v1/ws/'

// 用于管理项目内的Websocket连接
export class WebSocketConn extends WebSocket {
  // 构建函数 - 发起链接之后需要使用JWT的Token进行认证校验
  constructor(url: string, token: string, protocols?: string | string[]) {
    super(url, protocols)
    // 初始化等待队列
    this.waitQueue = []

    // 挂载处理器
    this.onmessage = messageHandler

    // 用户鉴权
    this.onopen = (ev: Event): any => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.wsOpAuth,
        data: token,
      }
      this.send(JSON.stringify(msg))

      // 完成等待队列
      this.waitQueue.forEach((fn) => {
        fn(this)
      })
      this.waitQueue = []
    }
  }

  // Websocket挂钩的终端
  public term: Terminal

  public waitQueue: ((ws: WebSocketConn) => any)[]
}

// 后端定义Websocket交互对象
export interface WebsocketMessage {
  op: WebsocketOperation
  data: any
}

// 后端定义的标准操作枚举
export enum WebsocketOperation {
  newPty = 1, // 发送创建新终端
  Stdin, // 发送终端输入
  Stdout, // 接收终端输出
  Resize, // 发送终端窗口调整
  Msg, // 接收后端消息
  ResourceApply, // 发送资源请求
  wsOpAuth, // 客户端向服务端发起鉴权
}

// 后端消息处理器，用于处理接收的数据
function messageHandler(this: WebSocketConn, ev: MessageEvent): any {
  // 反序列化
  const msg: WebsocketMessage = JSON.parse(ev.data)

  console.log('接收到websocket消息 操作码:', msg.op, '数据:', msg.data)

  switch (msg.op) {
    // 终端输出
    case WebsocketOperation.Stdout:
      if (this.term != null) {
        this.term.write(msg.data)
      }
      break

    // 页面通知 - 采用Notification的形式
    case WebsocketOperation.Msg:
      const resp = new BaseResponse(msg.data)

      // 后端约束
      if (!resp.IsSuccess()) {
        notification.error({
          message: '警告通知',
          description: resp.Data,
        })
        break
      }

      notification.success({
        message: '成功通知',
        description: resp.Data,
      })
      break
    default:
      console.log('unkown websocket msg:', ev.data)
  }
  return
}
