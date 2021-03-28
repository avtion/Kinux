<template>
  <a-card
    title="考试选择"
    :bordered="false"
    :loading="isProjectDataLoading"
    class="examSelector"
  >
    <a-list :data-source="departmentLessonData">
      <template #renderItem="{ item }">
        <a-list-item>
          <a
            @click="junmpToMission(item.id)"
            class="container mx-auto bg-gray-50 h-auto w-full p-4 space-y-1 rounded-lg hover:bg-gray-200 shadow"
          >
            <div class="block font-extrabold font-sans relative">
              <div class="inline-block text-blue-400 text-2xl">|</div>
              <div class="inline-block ml-2 text-gray-700 text-lg">
                {{ item.name }} ({{ item.lesson_name }})
              </div>
              <div class="inline-block absolute inset-y-0 right-0">></div>
            </div>
            <div class="block text-sm">
              <div class="inline-block pr-8 leading-relaxed h-aut">
                <div class="inline-block text-red-400 text-2xl">|</div>
                考试总分: {{ item.total }} 分
              </div>
              <div class="inline-block pr-8 leading-relaxed h-aut">
                <div class="inline-block text-green-400 text-2xl">|</div>
                考试时长: {{ item.time_limit }}
              </div>
              <div class="inline-block pr-8 leading-relaxed h-aut">
                <div class="inline-block text-yellow-400 text-2xl">|</div>
                考试时间: {{ item.begin_at }} 至 {{ item.end_at }}
              </div>
            </div>

            <div
              class="block text-sm pr-8 leading-relaxed h-auto line-clamp-3 text-gray-500"
              v-if="item.desc !== ''"
            >
              {{ item.desc }}
            </div>
          </a>
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
    const junmpToMission = (lesson: number) => {
      router.push({ name: 'missionSelector', params: { lesson: lesson } })
    }

    return {
      isProjectDataLoading: false,
      departmentLessonData,
      junmpToMission,
    }
  },
}
</script>

<style>
.examSelector {
  width: 100%;
}
</style>