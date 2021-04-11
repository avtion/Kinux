<template>
  <div class="workSpace">
    <!-- 头部展示栏 -->
    <a-page-header
      :ghost="false"
      style="border: 1px solid rgb(235, 237, 240)"
      title="Workspace 在线实验空间"
      :breadcrumb="{ routes }"
      sub-title=""
    >
      <!-- 主内容 -->
      <a-row>
        <a-space :size="30">
          <div></div>
          <!-- 头像 -->
          <div class="avatar">
            <a-avatar size="large" :src="avatar" />
          </div>

          <!-- 用户名 -->
          <a-statistic
            title="用户名"
            :value="
              profile.realName == '' ? profile.username : profile.realName
            "
          >
          </a-statistic>
          <a-divider type="vertical" />
          <!-- 班级 -->
          <a-statistic title="班级" :value="profile.department"> </a-statistic>

          <a-divider type="vertical" />
          <!-- 用户类型 -->
          <a-statistic title="用户类型" :value="profile.role"> </a-statistic>
        </a-space>
      </a-row>
    </a-page-header>

    <!-- 下方表格 -->
    <a-layout class="table">
      <a-layout-content
        :style="{
          background: '#fff',
          minHeight: '400px',
        }"
      >
        <!-- 二级路由 -->
        <router-view></router-view>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script lang="ts" type="module">
// vue
import { reactive, ref, inject } from 'vue'

// 图标生成
import Avatars from '@dicebear/avatars'
import AvatarsSprites from '@dicebear/avatars-avataaars-sprites'

// websocket
import { WebSocketConn } from '@/utils/websocketConn'

// store
import { GetStore } from '@/store/store'

// vue-router
import { Profile } from '@/store/interfaces'

export default {
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()

    // 从上下文中获取对象
    const ws: WebSocketConn = inject<WebSocketConn>('websocket')

    // 用户资料
    const profile = <Profile>store.getters.GetProfile

    // 顶部breadcrumb路径
    const breadcrumbPath = reactive([
      {
        path: '/',
        breadcrumbName: 'Kinux平台',
      },
      {
        path: '/dashboard',
        breadcrumbName: '在线实验',
      },
    ])

    // 头像
    const avatar = new Avatars(AvatarsSprites, {
      dataUri: true,
    }).create(<string>store.getters.GetAvatarSeed)

    return {
      routes: breadcrumbPath,
      avatar,
      profile,
    }
  },
  methods: {},
}
</script>

<style scoped lang="less">
.workSpace {
  width: 100%;
  height: 100%;
}
.page-header-content {
  display: flex;

  .welcome-text {
    display: none;
  }
}

.avatar {
  text-align: center;
  vertical-align: middle;
}

.stat-item {
  position: relative;
  display: inline-block;
  padding: 0 32px;

  &::after {
    position: absolute;
    top: 8px;
    right: 0;
    width: 1px;
    height: 40px;
    background-color: #e8e8e8;
    content: '';
  }

  &:last-child {
    padding-right: 0;

    &::after {
      display: none;
    }
  }
}

.table {
  width: 100%;
  padding: 24px 24px 24px;
}
</style>
