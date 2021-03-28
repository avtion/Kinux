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
      <!-- 底部菜单 -->
      <!-- <template #footer>
        <a-tabs :default-active-key="headerTypeOption">
          <a-tab-pane key="1" tab="实验" />
          <a-tab-pane key="2" tab="考试" />
        </a-tabs>
      </template> -->

      <!-- 主内容 -->
      <a-row type="flex">
        <!-- 左侧头像欢迎面板 -->
        <a-col :span="12">
          <div class="page-header-content">
            <div class="avatar">
              <a-avatar size="large" :src="avatar" />
            </div>
            <div class="content">
              <div class="content-title">
                {{
                  profile.realName == '' ? profile.username : profile.realName
                }}
                <span class="welcome-text">欢迎</span>
              </div>
              <div>{{ profile.role }} ｜ {{ profile.department }}</div>
            </div>
          </div>
        </a-col>

        <!-- 右侧统计面板 -->
        <a-col :span="4" class="stat-item">
          <a-statistic title="项目进度" :value="78" class="demo-class">
            <template #suffix>
              <span> / 100</span>
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="4" class="stat-item">
          <a-statistic title="班级排名" :value="1" class="demo-class">
            <template #suffix>
              <span> / 100</span>
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="4" class="stat-item">
          <a-statistic title="综合评价" value="S+" />
        </a-col>
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
import AvatarsSprites from '@dicebear/avatars-male-sprites'

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

  .avatar {
    flex: 0 1 72px;

    & > span {
      display: block;
      width: 72px;
      height: 72px;
      border-radius: 72px;
    }
  }

  .content {
    position: relative;
    top: 4px;
    flex: 1 1 auto;
    margin-left: 24px;
    color: rgba(0, 0, 0, 0.85);
    line-height: 22px;

    .content-title {
      margin-bottom: 12px;
      color: rgba(0, 0, 0, 0.85);
      font-weight: 500;
      font-size: 20px;
      line-height: 28px;
    }
  }
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
