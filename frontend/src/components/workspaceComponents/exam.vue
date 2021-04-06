<template>
  <a-card
    title="考试选择"
    :bordered="false"
    :loading="isProjectDataLoading"
    class="examSelector"
  >
    <a-list :data-source="departmentLessonData">
      <template #renderItem="{ item, index }">
        <a-list-item>
          <!-- 元数据 -->
          <a-list-item-meta :description="examDescCreator(item)">
            <!-- 标题 -->
            <template #title>
              <a>{{ item.name }}</a>
            </template>
            <!-- 头像 -->
            <template #avatar>
              <a-avatar :src="numberCreatorFn(index + 1)" />
            </template>
          </a-list-item-meta>
        </a-list-item>
      </template>
    </a-list>
  </a-card>
</template>

<script lang="ts" type="module">
// antd
import { RightOutlined } from '@ant-design/icons-vue'

// vue-request
import { useRequest } from 'vue-request'
import { BaseResponse, defaultClient } from '@/apis/request'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

// vue-router
import { Profile } from '@/store/interfaces'

// 图标生成
import Avatars from '@dicebear/avatars'
import sprites from '@dicebear/avatars-initials-sprites'

const apiPath = {
  list: '/v1/exam/dp/',
}

export default {
  components: {
    RightOutlined,
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const router = useRouter()

    // 用户资料
    const profile = <Profile>store.getters.GetProfile

    // 获取课程数据
    type ListParams = {
      dp: number
    }
    type ListResult = {
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
    }
    const departmentLessonDataAPI = (params: ListParams) => {
      return defaultClient.get<BaseResponse>(apiPath.list, {
        params: params,
      })
    }
    const getListParams = (): ListParams => {
      return <ListParams>{
        dp: Number(profile.dpID),
      }
    }
    const {
      data: departmentLessonData,
      run: getListData,
      loading: isListDataLoading,
    } = useRequest(departmentLessonDataAPI, {
      defaultParams: [getListParams()],
      formatResult: (res): ListResult[] => {
        return res.data.Data
      },
    })

    // 实验跳转
    const junmpToMission = (exam: number) => {
      router.push({ name: 'examMissionSelector', params: { exam: exam } })
    }

    // 序号
    const numberCreator = new Avatars(sprites, {
      dataUri: true,
    })
    const numberCreatorFn = (str: any): string => {
      return numberCreator.create(str + '')
    }

    // 考试描述生成器
    const examDescCreator = (item: ListResult): string => {
      return `考试时间${item.begin_at}-${item.end_at} \n限时:${item.time_limit}`
    }

    return {
      isProjectDataLoading: false,
      departmentLessonData,
      junmpToMission,
      numberCreatorFn,
      examDescCreator,
    }
  },
}
</script>

<style>
.examSelector {
  width: 100%;
}
</style>