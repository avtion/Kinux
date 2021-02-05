<template>
  <a-layout id="dashboard">
    <!-- 左侧导航栏 -->
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <!-- Logo -->
      <div class="logo">Kinux</div>
      <!-- 导航栏 -->
      <a-menu theme="dark" mode="inline" v-model:selectedKeys="selectedKeys">
        <a-menu-item key="1">
          <user-outlined />
          <span>nav 1</span>
        </a-menu-item>
        <a-menu-item key="2">
          <video-camera-outlined />
          <span>nav 2</span>
        </a-menu-item>
        <a-menu-item key="3">
          <upload-outlined />
          <span>nav 3</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <!-- 右侧内容 -->
    <a-layout>
      <!-- 头部 -->
      <a-layout-header style="background: #fff; padding: 0">
        <menu-unfold-outlined
          v-if="collapsed"
          class="trigger"
          @click="() => (collapsed = !collapsed)"
        />
        <menu-fold-outlined
          v-else
          class="trigger"
          @click="() => (collapsed = !collapsed)"
        />
      </a-layout-header>
      <!-- 内容页 -->
      <router-view></router-view>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, provide } from 'vue'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

// antd
import {
  UserOutlined,
  VideoCameraOutlined,
  UploadOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
} from '@ant-design/icons-vue'

// websocket
import {
  WebSocketConn,
  DefaultBackendWebsocketRoute,
} from '@/utils/websocketConn'

export default defineComponent({
  components: {
    UserOutlined,
    VideoCameraOutlined,
    UploadOutlined,
    MenuUnfoldOutlined,
    MenuFoldOutlined,
  },
  setup() {
    // vue相关变量
    const store = GetStore()
    const router = useRouter()

    // 获取JWT密钥
    const token = store.getters.GetJWTToken
    if ((token as string) && token == "") {
      console.log("JWT密钥失效");
      router.push("/");
    }

    // 加载页面之后建立Websocket链接
    const ws = new WebSocketConn(DefaultBackendWebsocketRoute, token)
    provide('websocket', ws)

    return {
      selectedKeys: ['1'],
      collapsed: false,
    }
  },
})
</script>

<style lang="less" scoped>
#dashboard {
  height: 100%;
  width: 100%;
  .trigger {
    font-size: 18px;
    line-height: 64px;
    padding: 0 24px;
    cursor: pointer;
    transition: color 0.3s;
  }
  .trigger:hover {
    color: #1890ff;
  }
  .logo {
    height: 32px;
    background: rgba(255, 255, 255, 0.2);
    margin: 16px;
  }
}
</style>
