<template>
  <a-card title="课程选择" :bordered="false" :loading="isProjectDataLoading">
    <a-list :data-source="departmentLessonData">
      <template #renderItem="{ item, index }">
        <a-list-item>
          <!-- 元数据 -->
          <a-list-item-meta :description="item.desc">
            <!-- 标题 -->
            <template #title>
              <span class="font-semibold">{{ item.name }}</span>
            </template>
            <!-- 头像 -->
            <template #avatar>
              <a-avatar :src="numberCreatorFn(index + 1)" />
            </template>
          </a-list-item-meta>
          <!-- 操作 -->
          <template #actions>
            <a-button @click="junmpToMission(item.id)">进入课程 </a-button>
          </template>
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

// 图标生成
import Avatars from '@dicebear/avatars'
import sprites from '@dicebear/avatars-initials-sprites'

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

    // 序号
    const numberCreator = new Avatars(sprites, {
      dataUri: true,
      background: '#3B82F6',
    })
    const numberCreatorFn = (str: any): string => {
      return numberCreator.create(str + '')
    }

    return {
      isProjectDataLoading: false,
      departmentLessonData,
      junmpToMission,
      numberCreatorFn,
    }
  },
}
</script>

<style>
</style>