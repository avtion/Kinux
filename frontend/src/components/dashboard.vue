<template>
  <a-layout id="dashboard">
    <!-- 左侧导航栏 -->
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <!-- Logo -->
      <div class="logo">{{ logo }}</div>
      <!-- 导航栏 -->
      <a-menu
        theme="dark"
        mode="inline"
        v-model:selectedKeys="selectedKeys"
        @click="menuClickFn"
      >
        <a-menu-item key="workspace">
          <FundOutlined />
          <span>学习空间</span>
        </a-menu-item>
        <a-menu-item key="profile">
          <UserOutlined />
          <span>个人资料</span>
        </a-menu-item>
        <a-menu-item key="session">
          <AppstoreOutlined />
          <span>实验会话</span>
        </a-menu-item>
        <a-sub-menu>
          <template #title>
            <DatabaseOutlined />
            <span>后台管理</span>
          </template>
          <a-menu-item key="departmentManager">班级管理</a-menu-item>
          <a-menu-item key="AccountManager">用户管理</a-menu-item>
          <a-menu-item key="deploymentManager">配置管理</a-menu-item>
          <a-menu-item key="missionManager">实验管理</a-menu-item>
          <a-menu-item key="examManager">考试管理</a-menu-item>
          <a-menu-item key="checkpointManager">检查点管理</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>
    <!-- 右侧内容 -->
    <a-layout>
      <!-- 头部 -->
      <a-layout-header style="background: #fff; padding: 0">
        <a-row type="flex">
          <a-col flex="100px">
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
          </a-col>
          <a-col flex="auto"> </a-col>
          <a-col flex="100px">
            <a-button @click="logout">
              <template #icon><UnlockOutlined /></template>
              注销
            </a-button>
          </a-col>
        </a-row>
      </a-layout-header>
      <!-- 内容页 -->
      <router-view></router-view>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, provide, ref, watch } from 'vue'

// antd
import { notification } from 'ant-design-vue'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

// antd
import {
  UserOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  FundOutlined,
  UnlockOutlined,
  DatabaseOutlined,
  AppstoreOutlined,
} from '@ant-design/icons-vue'

// websocket
import {
  WebSocketConn,
  DefaultBackendWebsocketRoute,
} from '@/utils/websocketConn'

export default defineComponent({
  components: {
    UserOutlined,
    MenuUnfoldOutlined,
    MenuFoldOutlined,
    FundOutlined,
    UnlockOutlined,
    DatabaseOutlined,
    AppstoreOutlined,
  },
  setup() {
    // logo
    const logo = ref<string>('Kinux 实验平台')

    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 获取JWT密钥
    const token = store.getters.GetJWTToken
    if ((token as string) && token == '') {
      console.log('JWT密钥失效')
      routers.push('/')
    }

    // 加载页面之后建立Websocket链接
    const ws = new WebSocketConn(DefaultBackendWebsocketRoute, token)
    provide('websocket', ws)

    // 开关按钮
    const collapsed = ref(false)
    watch(collapsed, (newValue, oldValue) => {
      if (newValue) {
        logo.value = 'K'
      } else {
        logo.value = 'Kinux 实验平台'
      }
    })

    // 登出
    const logout = () => {
      store.commit('ClearJWT')
      store.commit('ClearProfile')
      routers.push('/')
      notification.success({ message: '注销成功' })
    }

    // 导航栏触发函数
    const menuClickFn = ({ item, key, keyPath }) => {
      if (routers.currentRoute.value.name == key) {
        return
      }
      routers.push({ name: key })
    }

    return {
      selectedKeys: [routers.currentRoute.value.name],
      collapsed: collapsed,
      logo: logo,
      logout,
      menuClickFn,
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
    height: 64px;
    padding: 0 24px;
    text-align: center;
    overflow: hidden;
    color: rgba(255, 255, 255, 0.85);
    font-size: 18px;
    font-family: Avenir, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
      'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji',
      'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji', sans-serif;
    line-height: 64px;
    white-space: nowrap;
    text-decoration: none;
    background: #000c17;
  }
}
</style>
