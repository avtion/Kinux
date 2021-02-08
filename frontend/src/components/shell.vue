<template>
  <div class="back">
    <a-card class="markdown" title="说明文档"> whatever content </a-card>
    <a-card class="terminal" title="实操终端">
      <template #extra>
        <a-button-group>
          <a-button type="danger">重置容器</a-button>
          <a-button type="primary">结束实验</a-button>
        </a-button-group>
      </template>
      <div class="xterm terminal-container" ref="terminalRef"></div>
    </a-card>
  </div>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, ref, onMounted, inject, watch, reactive } from 'vue'

// xterm
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { WebLinksAddon } from 'xterm-addon-web-links'
import { FitAddon } from 'xterm-addon-fit'

// websocket
import { WebSocketConn, WebsocketMessage, WebsocketOperation } from '@/utils/websocketConn'
import App from '@/App.vue'

export default defineComponent({
  components: { App },
  name: 'shell',
  props: {
    id: Number,
  },
  setup(props, ctx) {
    // 从上下文中获取对象
    const ws: WebSocketConn =inject<WebSocketConn>('websocket')

    // 终端的DOM
    const terminalRef = ref<HTMLDivElement>()

    // 创建终端对象
    const ter = new Terminal({
      fontFamily: 'monaco, Consolas, "Lucida Console", monospace',
      fontSize: 14,
      cursorStyle: 'underline', //光标样式
      cursorBlink: true, //光标闪烁
      theme: {
          background: '#1b1b1b'
      },
    })
    ter.onData((input: string): void => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.Stdin,
        data: input
      }
      console.log(msg)
      ws.send(JSON.stringify(msg))
    })

    // 插件 - DOM适应器
    const fitAddon = new FitAddon()
    ter.loadAddon(fitAddon)
    window.addEventListener('resize', (e: UIEvent) =>{
      fitAddon.fit()
    }, false)
    ter.onResize((size: {
          cols: number;
          rows: number;
      }): any => {
      const msg: WebsocketMessage = {
        op: WebsocketOperation.Resize,
        data: size
      }
      if (ws.readyState !== WebSocket.OPEN) {
        ws.waitQueue.push((_ws)=>{
          ws.send(JSON.stringify(msg))
        })
      } else {
        ws.send(JSON.stringify(msg))
      }
    })


    // 插件 - 链接检测器
    ter.loadAddon(new WebLinksAddon())

    onMounted(() => {
      ter.open(terminalRef.value)
      ter.focus()

      ws.term = ter
      if (ws.readyState !== WebSocket.OPEN) {
        ws.waitQueue.push((_ws)=>{
          connectToPOD(ws, props.id, "")
          setTimeout(() => {
            fitAddon.fit()
          }, 1);
        })
      } else {
        connectToPOD(ws, 1, "")
        setTimeout(() => {
          fitAddon.fit()
        }, 1);
      }
    })

    return {
      ter,
      terminalRef
    }
  },
})

// 建立POD链接
function connectToPOD(ws: WebSocket, id: number, container: string): void {
  const msg: WebsocketMessage =  {
    op: WebsocketOperation.newPty,
    data: {
      id: id,
      container: container
    }
  }
  ws.send(JSON.stringify(msg))
  return
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
    height: 620px;
    padding: 0;
  }
}

.terminal-container {
  height: 100%;
  width: 100%;
}
</style>
