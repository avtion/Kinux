<template>
  <!-- 学生查询成绩 -->
  <div class="w-full h-full">
    <div class="w-full h-full p-5">
      <div class="w-full h-full bg-white rounded p-5">
        <a-space :size="30" class="pl-5">
          <div class="text-lg font-sans">课程 ></div>
          <!-- 课程选择 -->
          <a-select
            v-model:value="lessonSelect"
            style="width: 160px"
            placeholder="请选择课程"
          >
            <a-select-option
              v-for="(value, index) in departmentLessonData"
              :key="index"
              :value="value.id"
              >{{ value.name }}</a-select-option
            >
          </a-select>
        </a-space>

        <a-divider>成绩查询</a-divider>

        <!-- 折叠面板 -->
        <a-collapse v-model:activeKey="activeKey">
          <!-- 实验列表 -->
          <a-collapse-panel key="1" header="实验">
            <a-list
              item-layout="horizontal"
              :data-source="missionData"
              v-if="missionData !== undefined && missionData.length > 0"
            >
              <template #renderItem="{ item, index }">
                <a-list-item>
                  <!-- 元数据 -->
                  <a-list-item-meta :description="item.desc">
                    <!-- 标题 -->
                    <template #title>
                      <div>{{ item.name }}</div>
                    </template>
                    <!-- 头像 -->
                    <template #avatar>
                      <a-avatar :src="numberCreatorFn(index + 1)" />
                    </template>
                  </a-list-item-meta>
                  <!-- 操作 -->
                  <template #actions>
                    <a-button @click="jumpToMissionScore(item.id)"
                      >查询
                    </a-button>
                  </template>
                </a-list-item>
              </template>
            </a-list>
            <a-empty description="暂无实验数据" v-else />
          </a-collapse-panel>
          <!-- 考试列表 -->
          <a-collapse-panel key="2" header="考试">
            <a-list
              item-layout="horizontal"
              :data-source="examData"
              v-if="examData !== undefined && examData.length > 0"
            >
              <template #renderItem="{ item, index }">
                <a-list-item>
                  <!-- 元数据 -->
                  <a-list-item-meta :description="item.desc">
                    <!-- 标题 -->
                    <template #title>
                      <div>{{ item.name }}</div>
                    </template>
                    <!-- 头像 -->
                    <template #avatar>
                      <a-avatar :src="numberCreatorFn(index + 1)" />
                    </template>
                  </a-list-item-meta>
                  <!-- 操作 -->
                  <template #actions>
                    <a-button>查询 </a-button>
                  </template>
                </a-list-item>
              </template>
            </a-list>
            <a-empty description="暂无考试数据" v-else />
          </a-collapse-panel>
        </a-collapse>
      </div>
    </div>
  </div>
</template>

<script lang="ts" type="module">
import { defineComponent, ref, watch } from 'vue'

// antd
import { RightOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

// vue-request
import { useRequest } from 'vue-request'
import { BaseResponse, defaultClient } from '@/apis/request'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

// vue-router
import { Profile } from '@/store/interfaces'

// apis
import { mission, missionList, missionStatus } from '@api/mission'

// 图标生成
import Avatars from '@dicebear/avatars'
import sprites from '@dicebear/avatars-initials-sprites'

// 考试状态
import { examStatus, examRunningInfo, exam } from '@/apis/exam'

const lessonAPIPaths = {
  list: '/v2/dl/list',
}

const examAPIPaths = {
  list: '/v1/exam/dp/',
  check: '/v1/exam/check/',
}

export default defineComponent({
  setup(props) {
    // vue相关变量
    const store = GetStore()
    const router = useRouter()

    // 折叠面板
    const activeKey = ref(['1', '2'])

    // 用户资料
    const profile = <Profile>store.getters.GetProfile

    // 序号
    const numberCreator = new Avatars(sprites, {
      dataUri: true,
      background: '#60A5FA',
    })
    const numberCreatorFn = (str: any): string => {
      return numberCreator.create(str + '')
    }

    // 课程选择
    const lessonSelect = ref<number>()
    type lessonParams = {
      page: number
      size: number
      department: number
    }
    type LessonResult = {
      id: number
      name: string
      desc: string
    }
    const {
      data: departmentLessonData,
      run: getLessonData,
      loading: isLessonDataLoading,
    } = useRequest(
      (params: lessonParams) => {
        return defaultClient.get<BaseResponse>(lessonAPIPaths.list, {
          params: params,
        })
      },
      {
        defaultParams: [
          <lessonParams>{
            page: 0,
            size: 0,
            department: Number(profile.dpID),
          },
        ],
        formatResult: (res): LessonResult[] => {
          const _res = <LessonResult[]>res.data.Data
          if (lessonSelect.value === undefined && _res.length > 0) {
            lessonSelect.value = _res[0].id
          }
          return _res
        },
      }
    )
    // 监控课程选择以获取实验和考试列表
    watch(lessonSelect, (newValue) => {
      if (newValue !== undefined || newValue !== 0) {
        getMissionData(<missionReqParams>{
          page: 0,
          size: 0,
          lesson: lessonSelect.value,
        })
        getExamData(<examParams>{
          dp: Number(profile.dpID),
        })
      }
    })

    // 实验列表
    interface missionReqParams {
      page: number
      size: number
      lesson: number
    }
    const { data: missionData, run: getMissionData } = useRequest(
      (params: missionReqParams) => {
        return defaultClient.get<BaseResponse>('/v2/ms/', {
          params: params,
        })
      },
      {
        formatResult: (res): missionList[] => {
          return res.data.Data
        },
        manual: true,
      }
    )

    // 考试列表
    type examParams = {
      dp: number
    }
    type examResult = {
      id: number
      name: string
      desc: string
      total: number
      force_order: boolean
      begin_at: string
      end_at: string
      created_at: string
      time_limit: string
      begin_at_unix: number
      end_at_unix: number
      created_at_unix: number
      time_limit_unix: number
      lesson: number
      lesson_name: string
      lesson_desc: string
      exam_status: examStatus
    }
    const { data: examData, run: getExamData } = useRequest(
      (params: examParams) => {
        return defaultClient.get<BaseResponse>(examAPIPaths.list, {
          params: params,
        })
      },
      {
        formatResult: (res): examResult[] => {
          return res.data.Data
        },
      }
    )

    // 跳转至实验成绩查询
    const jumpToMissionScore = (missionID: number) => {
      router.push({
        name: 'msScore',
        params: { lessonID: lessonSelect.value, missionID: missionID },
      })
    }
    return {
      activeKey,

      // 序号
      numberCreatorFn,

      // 课程
      departmentLessonData,
      lessonSelect,

      // 实验
      missionData,

      // 考试
      examData,

      // 跳转
      jumpToMissionScore,
    }
  },
})
</script>

<style>
</style>