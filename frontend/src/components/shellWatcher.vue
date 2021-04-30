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
        <!-- 容器管理 -->
        <a-button-group>
          <a-button type="primary" @click="openSendMsgModal">发送消息</a-button>
          <a-button type="default" @click="comfirmToLeave">返回</a-button>
        </a-button-group>
      </template>
      <!-- 底部切换 -->
      <template #footer>
        <a-tabs defaultActiveKey="ter" @change="tabHandler">
          <a-tab-pane key="ter" tab="实验终端" />
          <a-tab-pane key="doc" tab="实验文档" />
          <a-tab-pane key="checkpoint" tab="考点状态" />
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

    <!-- 考点 -->
    <div class="w-full h-4/5 mt-2 p-3" v-show="currentTab === 'checkpoint'">
      <div class="bg-white rounded p-8">
        <!-- 折叠面板 -->
        <a-collapse v-model:activeKey="opened">
          <a-collapse-panel
            v-for="(v, k) in cps"
            :key="k + ''"
            :header="`容器: ${v.container_name}`"
            :disabled="true"
          >
            <a-list
              :data-source="v.data"
              item-layout="vertical"
              :bordered="true"
            >
              <template #renderItem="{ item }">
                <a-list-item>
                  <!-- 顶部 -->
                  <a-list-item-meta :description="`分数占比 ${item.percent}%`">
                    <template #title>
                      <a>{{ item.cp_name }}</a>
                    </template>
                  </a-list-item-meta>
                  <!-- 描述 -->
                  <div>
                    <a-alert
                      :message="`目标指令: ${item.cp_command}`"
                    ></a-alert>
                  </div>

                  <div class="mt-2">{{ item.cp_desc }}</div>
                </a-list-item>
              </template>
            </a-list>
          </a-collapse-panel>
        </a-collapse>
      </div>
    </div>

    <!-- 发送消息使用的Modal -->
    <a-modal
      v-model:visible="sendMsgModalVisible"
      title="发送消息"
      :afterClose="sendMsgModalAfterClose"
      width="720px"
      @ok="sendMsg"
    >
      <a-form-item label="消息内容">
        <a-textarea
          auto-size
          v-model:value="msgText"
          placeholder="请输入要发送的消息"
          style="width: 420px"
        />
      </a-form-item>
    </a-modal>
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
  computed,
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
import { Modal, notification, message } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  WarningOutlined,
  DownOutlined,
  CodeSandboxOutlined,
  CheckOutlined,
  LoadingOutlined,
  SmileOutlined,
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
import { Checkpoint, missionCheckpointRes } from '@/apis/checkpoint'

// 图标生成
import { IntCreator } from '@/utils/avatar'

const sendMsgAPI = (params: { target_id: number; text: string }) => {
  return defaultClient.post<BaseResponse>('/v1/ws/msg/', params)
}

export default defineComponent({
  components: {
    App,
    CodeSandboxOutlined,
    DownOutlined,
    CheckOutlined,
    LoadingOutlined,
    SmileOutlined,
  },
  name: 'shell',
  props: {
    mission: {
      type: String,
      default: '',
    },
    exam: {
      type: String,
      default: '',
    },
    lesson: {
      type: String,
      default: '',
    },
    target: {
      type: String,
      default: '',
    },
  },
  setup(props, ctx) {
    // 从上下文中获取对象
    const ws: WebSocketConn = inject<WebSocketConn>('websocket')

    // 路由
    const router = useRouter()

    // 退出实验环境
    const leaveShell = () => {
      router.push({ name: 'session' })
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
      message.error('监控状态下无法输入')
    })

    // 建立POD链接
    const connectToPOD = (ws: WebSocketConn) => {
      console.log(props.target)
      const msg: WebsocketMessage = {
        op: WebsocketOperation.AttachOtherWsWriter,
        data: {
          target_id: props.target,
        },
      }
      ter.clear()
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
    onMounted(() => {
      connectToPOD(ws)
    })
    onUnmounted(() => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.StopAttachOtherWsWriter,
        data: {},
      }
      ws.send(JSON.stringify(msg))
    })

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

    // 确定是否离开当前终端页面
    const comfirmToLeave = () => {
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
      // shutdownPtyConn(ws)
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
    let _size
    const tabHandler = (activeKey) => {
      // 修复切换面板时导致终端自闭
      if (activeKey === 'ter') {
        const size = _size
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
        _size = {}
      } else {
        if (_size === undefined || _size == {}) {
          _size = fitAddon.proposeDimensions()
        }
      }
      currentTab.value = activeKey
    }

    // 获取当前实验的考点
    const cp = new Checkpoint(
      Number(props.lesson),
      Number(props.mission),
      Number(props.exam)
    )
    const { data: cps, run: refreshCps, loading: isCpsLoading } = useRequest(
      cp.Get,
      {
        formatResult: (res) => {
          message.success('考点状态更新成功')
          return <missionCheckpointRes[]>res
        },
      }
    )
    const opened = computed(() => {
      if (cps.value === undefined) {
        return <string[]>[]
      }
      return cps.value.map((i, index) => {
        return index + ''
      })
    })

    // 序号
    const numberCreatorFn = (str: any): string => {
      return IntCreator(str + '', '#1D4ED8')
    }

    //发送消息
    const { run: sendMsgFn } = useRequest(sendMsgAPI, {
      manual: true,
    })
    const msgTargetID = ref<number>()
    const msgText = ref<string>()
    const sendMsgModalVisible = ref<boolean>(false)

    const sendMsg = () => {
      sendMsgModalVisible.value = false
      sendMsgFn({ target_id: msgTargetID.value, text: msgText.value })
    }

    const sendMsgModalAfterClose = () => {
      msgTargetID.value = 0
      msgText.value = ''
    }
    const openSendMsgModal = () => {
      sendMsgModalVisible.value = true
      msgTargetID.value = Number(props.target)
    }

    return {
      ter,
      terminalRef,
      comfirmToLeave,
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

      // 考点
      cps,
      refreshCps,
      isCpsLoading,
      opened,
      numberCreatorFn,

      // 发送消息
      sendMsg,
      msgText,
      sendMsgModalVisible,
      sendMsgModalAfterClose,
      openSendMsgModal,
    }
  },
})
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
