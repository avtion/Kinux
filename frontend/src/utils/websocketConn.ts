import { BaseResponse } from '@api/request'
import { notification } from 'ant-design-vue'
import { store } from '@/store/store'
import routers from '@/routers/routers'
import { Token } from '@api/user'
import { examRunningInfo } from '@api/exam'

// Antd全局提醒
import { message } from 'ant-design-vue'

// 时间处理
import { moment } from '@/utils/time'

// 考试状态
import { examInfo } from '@api/exam'

// 后端默认路由
export const DefaultBackendWebsocketRoute = `ws://${window.location.host}/ws`

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
        op: WebsocketOperation.Auth,
        data: token,
      }
      this.send(JSON.stringify(msg))

      // 完成等待队列
      this.waitQueue.forEach((fn) => {
        fn(this)
      })
      this.waitQueue = []
    }

    this.onclose = (ev: Event): any => {
      notification.error({
        message: '警告通知',
        description: 'websocket链接已关闭, 请刷新以重新建立链接',
      })
    }

    // 回调
    this.callbacksMapper = new Map()
    this.sendWithCallback = (data, op, callback, once) => {
      this.callbacksMapper.set(
        op,
        (ws: WebSocketConn, msg: WebsocketMessage): void => {
          callback(ws, msg)
          if (once) {
            this.callbacksMapper.delete(op)
          }
        }
      )
      super.send(data)
    }
  }

  // 等待链接成功时回调的函数队列
  public waitQueue: ((ws: WebSocketConn) => any)[]

  // 回调Hash
  public callbacksMapper: Map<
    WebsocketOperation,
    (ws: WebSocketConn, msg?: WebsocketMessage) => void
  >

  // 发送数据并挂载回调函数
  public sendWithCallback: (
    data: string | ArrayBufferLike | Blob | ArrayBufferView,
    op: WebsocketOperation,
    callback: (ws: WebSocketConn, msg?: WebsocketMessage) => void,
    once?: boolean
  ) => void
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
  MissionApply, // 发送资源请求
  Auth, // 客户端向服务端发起鉴权
  RequireAuth, // 服务端要求客户端进行鉴权
  RefreshToken, // 刷新JWT密钥
  ShutdownPty, // 关闭终端链接（即向终端发送 EndOfTransmission）
  ResetContainers, // 重置容器
  ContainersDone, // 容器重置成功

  // 2021/04/11
  AttachOtherWsWriter, // 侵入其他Websocket链接
  LeaveExam, // 退出考试
  ExamRunning, // 考试进行中（用于主动告诉用户正在进行考试）
  StopAttachOtherWsWriter, // 停止侵入其他Websocket链接
}

// 后端消息处理器，用于处理接收的数据
function messageHandler(this: WebSocketConn, ev: MessageEvent): any {
  // 反序列化
  const msg: WebsocketMessage = JSON.parse(ev.data)

  console.log('接收到websocket消息 操作码:', msg.op, '数据:', msg.data)

  switch (msg.op) {
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

    // 服务端要求客户端进行鉴权
    case WebsocketOperation.RequireAuth:
      store.commit('ClearJWT')
      routers.push('/')
      this.close()
      break

    // 刷新JWT密钥
    case WebsocketOperation.RefreshToken:
      new Token(msg.data['token'], msg.data['ttl']).UpdateToken()
      break

    // 退出考试
    case WebsocketOperation.LeaveExam:
      examInfo.value = undefined
      message.success({
        content: `考试结束`,
        duration: 5,
      })
      routers.push({
        name: 'examSelector',
      })
      break

    // 考试进行中
    case WebsocketOperation.ExamRunning:
      const _res = <examRunningInfo>msg.data
      examInfo.value = _res
      break

    default:
      const fn = this.callbacksMapper.get(msg.op)
      if (fn == undefined) {
        console.log('unkown websocket msg:', ev.data)
        return
      }
      fn(this, msg)
  }
  return
}
