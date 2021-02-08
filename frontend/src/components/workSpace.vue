<template>
  <div class="workSpace">
    <!-- 头部展示栏 -->
    <a-page-header
      :ghost="false"
      style="border: 1px solid rgb(235, 237, 240)"
      title="WorkSpace 工作间"
      :breadcrumb="{ routes }"
      sub-title=""
    >
      <!-- 底部菜单 -->
      <template #footer>
        <a-tabs :default-active-key="headerTypeOption">
          <a-tab-pane key="1" tab="实验" />
          <a-tab-pane key="2" tab="考试" />
        </a-tabs>
      </template>

      <!-- 主内容 -->
      <a-row type="flex">
        <!-- 左侧头像欢迎面板 -->
        <a-col :span="12">
          <div class="page-header-content">
            <div class="avatar">
              <a-avatar size="large" />
            </div>
            <div class="content">
              <div class="content-title">
                早上好，Avtion，吃饭了吗
                <span class="welcome-text">欢迎</span>
              </div>
              <div>学生 ｜ 计算机科学系 - 17网络工程</div>
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
    <a-layout style="padding: 24px 24px 24px">
      <a-layout-content
        :style="{
          background: '#fff',
          minHeight: '400px',
        }"
      >
        <a-card
          title="实验项目"
          :bordered="false"
          :loading="isProjectDataLoading"
          :tab-list="tabList"
          @tabChange="(key) => onTabChange(key)"
        >
          <a-list item-layout="horizontal" :data-source="dataList">
            <template #renderItem="{ item, index }">
              <a-list-item>
                <!-- 元数据 -->
                <a-list-item-meta :description="item.Desc">
                  <!-- 标题 -->
                  <template #title>
                    {{ index + 1 }} |
                    <a href="https://www.antdv.com/">{{ item.Name }}</a>
                  </template>
                  <!-- 头像 -->
                  <template #avatar>
                    <a-avatar />
                  </template>
                </a-list-item-meta>
                <!-- 操作 -->
                <template #actions>
                  <a-button
                    :type="GetMissionButtonType(item.Status)"
                    :loading="GetMissionButtonLoadingStatus(item.Status)"
                    @click="MissionHandler(item)"
                  >
                    {{ GetMissionButtonDesc(item.Status) }}
                  </a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script lang="ts" type="module">
import { reactive, ref } from 'vue'
import { mission, missionList, missionStatus } from '@api/mission'
import routers from '@/routers/routers'

export default {
  setup() {
    // 顶部breadcrumb路径
    const breadcrumbPath = reactive([
      {
        path: '/',
        breadcrumbName: 'Kinux平台',
      },
      {
        path: '/dashboard',
        breadcrumbName: '终端',
      },
      {
        path: '/workspace',
        breadcrumbName: '工作间',
      },
    ])

    // 头部类型选卡: 1-实验 2-考试
    const headerTypeOption = ref(1)

    // 项目内的分类选卡
    const tabList = reactive([
      {
        key: '1',
        tab: '未完成',
      },
      {
        key: '2',
        tab: '已完成',
      },
    ])

    // 数据加载
    const isProjectDataLoading = ref(true)

    // 加载任务数据
    const dataList = ref(<missionList[]>[])
    new mission()
      .list()
      .then((res) => {
        dataList.value = res
      })
      .finally(() => {
        isProjectDataLoading.value = false
      })

    return {
      routes: breadcrumbPath,
      isProjectDataLoading,
      tabList,
      headerTypeOption,
      dataList,
      GetMissionButtonType,
      GetMissionButtonLoadingStatus,
      GetMissionButtonDesc,
      MissionHandler
    }
  },
  methods: {
    onTabChange(key) {
      console.log(key)
    },
  },
}

// 任务处理中的Status
const missionHandingStauts: number = -1

// 获取任务按钮的类型
function GetMissionButtonType(t: number): string {
  switch (t) {
    case missionStatus.Stop:
      return 'default'
    case missionStatus.Pending:
      return 'danger'
    case missionStatus.Working:
      return 'primary'
    case missionStatus.Done:
      return 'default'
    case missionHandingStauts:
      return 'danger'
    default:
      return 'dashed'
  }
}

// 获取任务按钮的状态
function GetMissionButtonLoadingStatus(t: number): boolean {
  return t == missionStatus.Pending || t == missionHandingStauts
}

// 获取任务按钮的描述
function GetMissionButtonDesc(t: number): string {
    switch (t) {
    case missionStatus.Stop:
      return '开始'
    case missionStatus.Pending:
      return '正在加载'
    case missionStatus.Working:
      return '进入'
    case missionStatus.Done:
      return '已结束'
    case missionHandingStauts:
      return '任务正在处理'
    default:
      return ''
  }
}

// 任务处理函数
function MissionHandler(m: missionList): void{
  console.log(m)
  const status = m.Status
  m.Status = missionHandingStauts
  switch (status) {
    case missionStatus.Stop:
      return
    case missionStatus.Pending:
      return
    case missionStatus.Working:
      routers.push({name:"shell", params: {id: m.ID}})
      return
    case missionStatus.Done:
      return
    default:
      return
  }
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
</style>
