<template>
  <div class="back h-full">
    <!-- 顶部 -->
    <a-page-header
      class="border border-solid border-gray-200 h-auto"
      title="虚拟实验环境"
      :ghost="false"
    >
      <!-- 按钮组 -->
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
          <a-button>切换容器</a-button>
        </a-dropdown>
        <a-divider type="vertical" />
        <!-- 容器管理 -->
        <a-button-group>
          <a-button type="danger" @click="comfirmToResetContainer"
            >重置环境</a-button
          >
          <a-button type="primary" @click="comfirmToShutdownMission"
            >关闭环境</a-button
          >
          <a-button type="default" @click="comfirmToLeave">返回</a-button>
        </a-button-group>
      </template>
      <!-- 底部切换 -->
      <template #footer>
        <a-tabs defaultActiveKey="ter" @change="tabHandler">
          <a-tab-pane key="ter" tab="实验终端" />
          <a-tab-pane key="doc" tab="实验文档" />
          <a-tab-pane key="checkpoint" tab="考点" />
        </a-tabs>
      </template>
      <!-- 描述 -->
      <a-descriptions size="small" :column="2">
        <a-descriptions-item label="实验名称">
          {{ missionInfo.name }}
        </a-descriptions-item>
        <a-descriptions-item label="实验总分">
          <a>{{ missionInfo.total }}</a>
        </a-descriptions-item>
        <a-descriptions-item label="实验描述">
          <span>{{ missionInfo.desc }}</span>
        </a-descriptions-item>
      </a-descriptions>
    </a-page-header>

    <!-- 实操终端 -->
    <div class="w-full h-4/5 mt-2 p-3" v-show="currentTab === 'ter'">
      <div class="w-full h-full rounded p-3" style="background-color: #1f2937">
        <!-- 终端 -->
        <div class="xterm terminal-container h-full" ref="terminalRef"></div>
      </div>
    </div>

    <!-- 实验文档 -->
    <div class="w-full h-4/5 mt-2 p-3" v-show="currentTab === 'doc'">
      <v-md-editor
        class="h-full"
        v-model="instructions"
        mode="preview"
        v-if="currentTab === 'doc'"
      ></v-md-editor>
    </div>

    <div
      class="w-full h-4/5 mt-2 p-3"
      v-show="currentTab === 'checkpoint'"
    ></div>
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

// vue-router
import { useRouter } from 'vue-router'

// xterm
import 'xterm/css/xterm.css'
import { Terminal, ITheme } from 'xterm'
import { WebLinksAddon } from 'xterm-addon-web-links'
import { FitAddon } from 'xterm-addon-fit'
const defaultTheme: ITheme = {
  background: '#1F2937',
  foreground: '#F3F4F6',
  selection: '#FFFFFF40',
  black: '#000000',
  red: '#cd3131',
  green: '#0DBC79',
  yellow: '#e5e510',
  blue: '#2472c8',
  magenta: '#bc3fbc',
  cyan: '#11a8cd',
  white: '#e5e5e5',
  brightBlack: '#666666',
  brightRed: '#f14c4c',
  brightGreen: '#23d18b',
  brightYellow: '#f5f543',
  brightBlue: '#3b8eea',
  brightMagenta: '#d670d6',
  brightCyan: '#29b8db',
  brightWhite: '#e5e5e5',
}

// antd
import { Modal, notification } from 'ant-design-vue'
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

// apis
import { mission } from '@/apis/mission'
import { BaseResponse, defaultClient } from '@/apis/request'
import { useRequest } from 'vue-request'

export default defineComponent({
  components: { App, CodeSandboxOutlined, DownOutlined },
  name: 'shell',
  props: {
    mission: String,
    exam: String,
    lesson: String,
  },
  setup(props, ctx) {
    // 从上下文中获取对象
    const ws: WebSocketConn = inject<WebSocketConn>('websocket')

    // 路由
    const router = useRouter()

    const leaveShell = (examID = '') => {
      if (examID !== '') {
        router.push({
          name: 'examMissionSelector',
          params: { exam: examID },
        })
        return
      }
      router.push({ name: 'workspace' })
      return
    }

    // 终端的DOM
    const terminalRef = ref<HTMLDivElement>()

    // 创建终端对象
    const ter = new Terminal({
      fontFamily: 'monaco, Consolas, "Lucida Console", monospace',
      fontSize: 16,
      cursorStyle: 'underline', //光标样式
      cursorBlink: true, //光标闪烁
    })
    ter.setOption('theme', defaultTheme)
    ter.onData((input: string): void => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.Stdin,
        data: input,
      }
      console.log(msg)
      ws.send(JSON.stringify(msg))
    })

    // 建立POD链接
    const connectToPOD = (
      ws: WebSocketConn,
      lesson: string,
      exam: string,
      mission: string,
      container: string
    ) => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.newPty,
        data: {
          mission_id: mission,
          lesson_id: lesson,
          exam_id: exam,
          container: container,
        },
      }
      ws.sendWithCallback(
        JSON.stringify(msg),
        WebsocketOperation.Stdout,
        (ws, msg) => {
          ter.write(msg.data)
        },
        false
      )
      return
    }

    // 插件 - DOM适应器
    const fitAddon = new FitAddon()
    ter.loadAddon(fitAddon)
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

    // 浏览器监听窗口变化
    const fitListener = (e: UIEvent) => {
      fitAddon.fit()
    }
    onMounted(() => {
      window.addEventListener('resize', fitListener, false)
    })
    onUnmounted(() => {
      window.removeEventListener('resize', fitListener, false)
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

      // 初始化函数
      const fn = () => {
        connectToPOD(ws, props.lesson, props.exam, props.mission, newValue)
        setTimeout(() => {
          const size = fitAddon.proposeDimensions()
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
        }, 1000)
      }

      // 将终端连接到新的控制台
      if (ws.readyState !== WebSocket.OPEN) {
        // 如果ws未准备就绪则压入等待队列
        ws.waitQueue.push((_ws) => {
          fn()
        })
      } else {
        fn()
      }
    })

    // 获取容器列表
    const containersNames = ref<string[]>()
    new mission().listContainersNames(props.mission).then((names: string[]) => {
      containersNames.value = names
      if (selectContainer.value == '') {
        selectContainer.value = containersNames.value[0]
      }
    })

    // 说明文档
    const instructions = ref<string>(
      `🤪无实验文档数据，请联系刷新页面或实验教师`
    )
    const instructionsLoading = ref<boolean>(true)
    new mission()
      .getGuide(props.mission)
      .then((res) => {
        instructions.value = res
      })
      .finally(() => {
        instructionsLoading.value = false
      })

    // 重置实验
    const comfirmToResetContainer = () => {
      Modal.confirm({
        title: '确定要重置实验容器吗?',
        icon: createVNode(WarningOutlined),
        content: '当你点击确认按钮，将会重置实验容器，一切数据将会被销毁！',
        okText: '确定',
        okType: 'danger',
        cancelText: '取消',
        onOk() {
          const msg: WebsocketMessage = {
            op: WebsocketOperation.ResetContainers,
            data: {
              id: props.mission,
            },
          }
          ws.sendWithCallback(
            JSON.stringify(msg),
            WebsocketOperation.ContainersDone,
            (ws) => {
              connectToPOD(
                ws,
                props.lesson,
                props.exam,
                props.mission,
                selectContainer.value
              )
              setTimeout(() => {
                fitAddon.fit()
              }, 1)
              ter.clear()
            },
            true
          )
        },
        onCancel() {},
      })
    }

    // 结束实验
    const comfirmToShutdownMission = () => {
      Modal.confirm({
        title: '确定要结束实验吗?',
        icon: createVNode(WarningOutlined),
        content:
          '当你点击确认按钮，将会结束实验并退回学习空间，一切数据将会被销毁！',
        okText: '确定',
        okType: 'danger',
        cancelText: '取消',
        onOk() {
          new mission().deleteDeployment(props.mission).then((res) => {
            notification.success({
              message: res,
            })
            leaveShell(props.exam)
          })
        },
        onCancel() {},
      })
    }

    // 确定是否离开当前终端页面
    const comfirmToLeave = () => {
      Modal.confirm({
        title: '想要退出终端吗?',
        icon: createVNode(ExclamationCircleOutlined),
        content: '当你点击确认按钮，将会关闭终端',
        okText: '确定',
        cancelText: '取消',
        onOk() {
          leaveShell(props.exam)
        },
        onCancel() {},
      })
    }

    // 页面挂载的钩子函数
    onMounted(() => {
      // 加载终端
      ter.open(terminalRef.value)
      ter.focus()
      setTimeout(() => {
        fitAddon.fit()
      }, 1)
    })

    // 页面卸载的钩子函数
    onUnmounted(() => {
      shutdownPtyConn(ws)
    })

    // 获取实验数据
    type missionInfoType = {
      id: number
      name: string
      desc: string
      total: number
    }
    const missionInfo = ref<missionInfoType>({
      id: 0,
      name: '',
      desc: '',
      total: 0,
    })
    defaultClient
      .get<BaseResponse>('/v1/mission/get/' + props.mission + '/')
      .then((res) => {
        missionInfo.value = res.data.Data
      })

    // 当前tab
    const currentTab = ref<string>('ter')
    const tabHandler = (activeKey) => {
      currentTab.value = activeKey
    }

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
      instructions,
      instructionsLoading,

      // 实验数据
      missionInfo,

      // 考试
      isExam: props.exam !== '',

      // 标签
      currentTab,
      tabHandler,
    }
  },
})

// 主动关闭Pty链接
function shutdownPtyConn(ws: WebSocketConn): void {
  console.log('主动关闭pty链接')
  const msg: WebsocketMessage = {
    op: WebsocketOperation.ShutdownPty,
    data: {},
  }
  ws.send(JSON.stringify(msg))
}
</script>

<style lang="less" scoped>
.back {
  background: #ececec;
  width: 100%;
}
.markdown {
  margin-bottom: 15px;
}
.terminal {
  :deep(.ant-card-body) {
    padding: 0;
  }
}

.terminal-container {
  width: 100%;
  :deep(.xterm) {
    height: 100%;
  }
  :deep(.xterm-viewport) {
    overflow-y: hidden;
  }
}
</style>
