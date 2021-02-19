<template>
  <div class="back">
    <a-card class="markdown" title="说明文档"> whatever content </a-card>
    <a-card class="terminal" title="实操终端">
      <!-- 按钮 -->
      <template #extra>
        <!-- 容器切换 -->
        <a-dropdown>
          <template #overlay>
            <a-menu @click="changeContainer">
              <a-menu-item
                v-for="(name, index) in containersNames"
                :key="index"
                :disabled="selectContainer == name"
              >
                <CodeSandboxOutlined />
                {{ name }} {{ selectContainer == name ? '| 当前容器' : '' }}
              </a-menu-item>
            </a-menu>
          </template>
          <a-button> 容器列表 </a-button>
        </a-dropdown>
        <a-divider type="vertical" />
        <!-- 容器管理 -->
        <a-button-group>
          <a-button type="danger" @click="comfirmToResetContainer"
            >重置容器</a-button
          >
          <a-button type="primary" @click="comfirmToShutdownMission"
            >结束实验</a-button
          >
          <a-button type="default" @click="comfirmToLeave">返回空间</a-button>
        </a-button-group>
      </template>
      <!-- 终端 -->
      <div class="xterm terminal-container" ref="terminalRef"></div>
    </a-card>
  </div>
</template>

<script lang="ts" type="module">
// vue
import {
  defineComponent,
  ref,
  onMounted,
  inject,
  createVNode,
  watch,
  onUnmounted,
} from 'vue'

// xterm
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { WebLinksAddon } from 'xterm-addon-web-links'
import { FitAddon } from 'xterm-addon-fit'

// antd
import { Modal } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  WarningOutlined,
  DownOutlined,
  CodeSandboxOutlined,
} from '@ant-design/icons-vue'

// websocket
import {
  WebSocketConn,
  WebsocketMessage,
  WebsocketOperation,
} from '@/utils/websocketConn'
import App from '@/App.vue'
import { mission } from '@/apis/mission'

// vue-router
import routers from '@/routers/routers'

export default defineComponent({
  components: { App, CodeSandboxOutlined, DownOutlined },
  name: 'shell',
  props: {
    id: String,
  },
  setup(props, ctx) {
    // 从上下文中获取对象
    const ws: WebSocketConn = inject<WebSocketConn>('websocket')

    // 终端的DOM
    const terminalRef = ref<HTMLDivElement>()

    // 创建终端对象
    const ter = new Terminal({
      fontFamily: 'monaco, Consolas, "Lucida Console", monospace',
      fontSize: 14,
      cursorStyle: 'underline', //光标样式
      cursorBlink: true, //光标闪烁
      theme: {
        background: '#1b1b1b',
      },
    })
    ter.onData((input: string): void => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.Stdin,
        data: input,
      }
      console.log(msg)
      ws.send(JSON.stringify(msg))
    })

    // 插件 - DOM适应器
    const fitAddon = new FitAddon()
    ter.loadAddon(fitAddon)
    window.addEventListener(
      'resize',
      (e: UIEvent) => {
        fitAddon.fit()
      },
      false
    )
    ter.onResize((size: { cols: number; rows: number }): any => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.Resize,
        data: size,
      }
      if (ws.readyState !== WebSocket.OPEN) {
        ws.waitQueue.push((_ws) => {
          ws.send(JSON.stringify(msg))
        })
      } else {
        ws.send(JSON.stringify(msg))
      }
    })

    // 插件 - 链接检测器
    ter.loadAddon(new WebLinksAddon())

    // 切换容器
    const changeContainer = (e: Event) => {
      console.log('click', e)
    }

    // 选择的容器
    const selectContainer = ref<string>('')
    watch(selectContainer, (newValue) => {
      console.log('当前选择的新容器', selectContainer.value)

      // 将终端连接到新的控制台
      if (ws.readyState !== WebSocket.OPEN) {
        ws.waitQueue.push((_ws) => {
          connectToPOD(ws, props.id, newValue)
          setTimeout(() => {
            fitAddon.fit()
          }, 1)
        })
      } else {
        connectToPOD(ws, props.id, newValue)
        setTimeout(() => {
          fitAddon.fit()
        }, 1)
      }
    })

    // 获取容器列表
    const containersNames = ref<string[]>()
    new mission().listContainersNames(props.id).then((names: string[]) => {
      containersNames.value = names
      if (selectContainer.value == '') {
        selectContainer.value = containersNames.value[0]
      }
    })

    // 页面挂载的钩子函数
    onMounted(() => {
      // 加载终端
      ter.open(terminalRef.value)
      ter.focus()
      ws.term = ter
      setTimeout(() => {
        fitAddon.fit()
      }, 1)
    })

    // 页面卸载的钩子函数
    onUnmounted(() => {
      shutdownPtyConn(ws)
    })

    return {
      ter,
      terminalRef,
      comfirmToLeave,
      comfirmToResetContainer,
      comfirmToShutdownMission,
      changeContainer,
      containersNames,
      selectContainer,
      leaveShell,
    }
  },
})

// 建立POD链接
function connectToPOD(ws: WebSocketConn, id: String, container: string): void {
  const msg: WebsocketMessage = {
    op: WebsocketOperation.newPty,
    data: {
      id: id,
      container: container,
    },
  }
  ws.send(JSON.stringify(msg))
  return
}

// 主动关闭Pty链接
function shutdownPtyConn(ws: WebSocketConn): void {
  console.log('主动关闭pty链接')
  const msg: WebsocketMessage = {
    op: WebsocketOperation.ShutdownPty,
    data: {},
  }
  ws.send(JSON.stringify(msg))
}

// 确定是否离开当前终端页面
function comfirmToLeave(): void {
  Modal.confirm({
    title: '想要退出终端吗?',
    icon: createVNode(ExclamationCircleOutlined),
    content: '当你点击确认按钮，将会关闭终端',
    okText: '确定',
    cancelText: '取消',
    onOk() {
      leaveShell()
    },
    onCancel() {},
  })
}

// 确定是否重置容器
function comfirmToResetContainer(): void {
  Modal.confirm({
    title: '确定要重置实验容器吗?',
    icon: createVNode(WarningOutlined),
    content: '当你点击确认按钮，将会重置实验容器，一切数据将会被销毁！',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      return new Promise((resolve, reject) => {
        setTimeout(Math.random() > 0.5 ? resolve : reject, 1000)
      })
    },
    onCancel() {},
  })
}

// 确定是否结束实验
function comfirmToShutdownMission(): void {
  Modal.confirm({
    title: '确定要结束实验吗?',
    icon: createVNode(WarningOutlined),
    content:
      '当你点击确认按钮，将会结束实验并退回学习空间，一切数据将会被销毁！',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      return new Promise((resolve, reject) => {
        setTimeout(Math.random() > 0.5 ? resolve : reject, 1000)
      })
    },
    onCancel() {},
  })
}

// 离开终端
function leaveShell() {
  routers.push({ name: 'workspace' })
}
</script>

<style lang="less" scoped>
.back {
  background: #ececec;
  padding: 15px;
  height: 100%;
  width: 100%;
}
.markdown {
  height: 250px;
  margin-bottom: 15px;
}
.terminal {
  height: 700px;
  :deep(.ant-card-body) {
    height: 100%;
    padding: 0;
  }
}

.terminal-container {
  height: 100%;
  width: 100%;
}
</style>
