<template>
  <a-card title="课程选择" :bordered="false" :loading="isProjectDataLoading">
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
                {{ item.name }}
              </div>
              <div class="inline-block absolute inset-y-0 right-0">></div>
            </div>
            <div class="block text-sm pr-8 leading-relaxed h-auto line-clamp-3">
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

// 当前考试状态
import { examInfo } from '@api/exam'

const apiPath = {
  list: '/v2/dl/list',
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
      page: number
      size: number
      department: number
    }
    type ListResult = {
      id: number
      name: string
      desc: string
    }
    const departmentLessonDataAPI = (params: ListParams) => {
      return defaultClient.get<BaseResponse>(apiPath.list, {
        params: params,
      })
    }
    const getListParams = (): ListParams => {
      return <ListParams>{
        page: 0,
        size: 0,
        department: Number(profile.dpID),
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
      if (examInfo.value !== undefined) {
        message.error({
          content: '处于考试状态下无法查看实验',
          key: '__toMissionBlocker',
        })
        return
      }
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
</style>