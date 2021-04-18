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
        <a-menu-item
          key="lessonSelector"
          class="menu"
          :disabled="examLeftTime !== null"
          v-if="!isAdmin"
        >
          <AppstoreOutlined class="align-middle" />
          <span class="align-middle">在线实验</span>
        </a-menu-item>
        <a-menu-item key="examSelector" v-if="!isAdmin">
          <FundOutlined class="align-middle" v-if="!isAdmin" />
          <span class="align-middle">在线考试</span>
        </a-menu-item>
        <a-menu-item key="stuScore" v-if="!isAdmin">
          <BarChartOutlined class="align-middle" />
          <span class="align-middle">成绩查询</span>
        </a-menu-item>
        <a-menu-item key="counterManager" v-if="isAdmin">
          <ClusterOutlined class="align-middle" />
          <span class="align-middle">系统统计</span>
        </a-menu-item>
        <a-menu-item key="profile">
          <UserOutlined class="align-middle" />
          <span class="align-middle">个人资料</span>
        </a-menu-item>
        <a-menu-item key="AccountManager" v-if="isAdmin">
          <EditOutlined class="align-middle" />
          <span class="align-middle">用户管理</span>
        </a-menu-item>
        <a-menu-item key="departmentManager" v-if="isAdmin">
          <DatabaseOutlined class="align-middle" />
          <span class="align-middle">班级管理</span>
        </a-menu-item>
        <a-menu-item key="lessonManager" v-if="isAdmin">
          <DropboxOutlined class="align-middle" />
          <span class="align-middle">课程管理</span>
        </a-menu-item>
        <a-menu-item key="missionManager" v-if="isAdmin">
          <CodepenOutlined class="align-middle" />
          <span class="align-middle">实验管理</span>
        </a-menu-item>
        <a-menu-item key="checkpointManager" v-if="isAdmin">
          <DeploymentUnitOutlined class="align-middle" />
          <span class="align-middle">考点管理</span>
        </a-menu-item>
        <a-menu-item key="examManager" v-if="isAdmin">
          <CodeSandboxOutlined class="align-middle" />
          <span class="align-middle">考试管理</span>
        </a-menu-item>
        <a-menu-item key="teaScore" v-if="isAdmin">
          <DotChartOutlined class="align-middle" />
          <span class="align-middle">成绩查询</span>
        </a-menu-item>
        <a-menu-item key="session" v-if="isAdmin">
          <DingtalkOutlined class="align-middle" />
          <span class="align-middle">会话管理</span>
        </a-menu-item>
        <a-menu-item key="deploymentManager" v-if="isAdmin">
          <FormOutlined class="align-middle" />
          <span class="align-middle">容器配置</span>
        </a-menu-item>
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
          <a-col flex="auto">
            <div v-if="examLeftTime !== null">
              <span class="inline-block">考试倒计时</span>
              <a-statistic-countdown
                :value="examLeftTime"
                format="HH:mm:ss:SSS"
                class="counter"
              />
            </div>
          </a-col>
          <a-col flex="100px">
            <a-button @click="logout" shape="round">
              <template #icon><UnlockOutlined class="align-middle" /></template>
              注销
            </a-button>
          </a-col>
        </a-row>
      </a-layout-header>
      <!-- 内容页 -->
      <a-layout-content>
        <router-view></router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, provide, ref, watch, computed } from 'vue'

// antd
import { notification } from 'ant-design-vue'

// store
import { GetStore } from '@/store/store'
import { Profile } from '@/store/interfaces'

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
  EditOutlined,
  FormOutlined,
  CodepenOutlined,
  CodeSandboxOutlined,
  DropboxOutlined,
  DingtalkOutlined,
  DeploymentUnitOutlined,
  BarChartOutlined,
  DotChartOutlined,
  ClusterOutlined,
} from '@ant-design/icons-vue'

// websocket
import {
  WebSocketConn,
  DefaultBackendWebsocketRoute,
} from '@/utils/websocketConn'

// 时间处理
import { moment } from '@/utils/time'

// 考试状态
import { examInfo } from '@api/exam'

import { Role } from '@/store/interfaces'

export default defineComponent({
  components: {
    UserOutlined,
    MenuUnfoldOutlined,
    MenuFoldOutlined,
    FundOutlined,
    UnlockOutlined,
    DatabaseOutlined,
    AppstoreOutlined,
    EditOutlined,
    FormOutlined,
    CodepenOutlined,
    CodeSandboxOutlined,
    DropboxOutlined,
    DingtalkOutlined,
    DeploymentUnitOutlined,
    BarChartOutlined,
    DotChartOutlined,
    ClusterOutlined,
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
      // 从上下文中获取对象
      ws.close()
      notification.success({ message: '注销成功' })
    }

    // 导航栏触发函数
    const menuClickFn = ({ item, key, keyPath }) => {
      if (routers.currentRoute.value.name == key) {
        return
      }
      routers.push({ name: key })
    }

    // 考试剩余时间
    const examLeftTime = computed(() => {
      if (examInfo.value == undefined) {
        return null
      }
      return moment()
        .add(moment.duration(Number(examInfo.value.left_time) / 1000000))
        .valueOf()
    })

    // 获取当前用户角色
    const p: Profile = store.getters.GetProfile
    const isTeacher = computed(() => {
      return p.roleID === Role.RoleManager
    })
    const isAdmin = computed(() => {
      return p.roleID == Role.RoleAdmin || isTeacher.value
    })

    return {
      selectedKeys: [routers.currentRoute.value.name],
      collapsed: collapsed,
      logo: logo,
      logout,
      menuClickFn,

      // 考试有关
      examInfo,
      examLeftTime,

      // 角色ID
      isAdmin,
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

.menu {
  span {
    vertical-align: middle;
  }
}

.counter {
  height: 100%;
  line-height: 64px;
  display: inline-block;
  margin-left: 10px;
  :deep(.ant-statistic-content) {
    // height: 100%;
    vertical-align: middle;
  }
}
</style>
