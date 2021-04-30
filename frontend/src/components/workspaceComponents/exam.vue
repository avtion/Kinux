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
          <!-- 操作 -->
          <template #actions>
            <a-button
              @click="examStartButtonFn(item)"
              :disabled="item.exam_status == examStatus.ESFinish"
            >
              {{ GetExamButtonDesc(item.exam_status) }}
            </a-button>
          </template>
        </a-list-item>
      </template>
    </a-list>
  </a-card>
</template>

<script lang="ts" type="module">
// antd
import { RightOutlined } from '@ant-design/icons-vue'

// modal
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { Modal } from 'ant-design-vue'

// 通知
import { notification } from 'ant-design-vue'

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
import { IntCreator } from '@/utils/avatar'

// 时间处理
import { moment } from '@/utils/time'

// API接口
const apiPath = {
  list: '/v1/exam/dp/',
  check: '/v1/exam/check/',
}

// 考试状态
import { examStatus, examRunningInfo, exam } from '@/apis/exam'

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
      exam_status: examStatus
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

    // 序号
    const numberCreatorFn = (str: any): string => {
      return IntCreator(str + '', '#F59E0B')
    }

    // 考试描述生成器
    const examDescCreator = (item: ListResult): string => {
      const beginAt = moment.unix(item.begin_at_unix).format('LLL')
      const endAt = moment.unix(item.end_at_unix).format('LLL')
      const time_limit = moment.duration(item.time_limit_unix, 's').minutes()
      return `考试时间: ${beginAt} - ${endAt} \n限时: ${time_limit}分钟\n描述: ${item.desc}`
    }

    // 开始考试按钮描述
    const GetExamButtonDesc = (status: examStatus): string => {
      switch (status) {
        case examStatus.ESFinish:
          return '已结束'
        case examStatus.ESNotStart:
          return '开始考试'
        case examStatus.ESRunning:
          return '正在考试'
        case examStatus.ESPassTime:
          return '不在考试时间内'
      }
    }

    // 开始考试
    const examStartButtonFn = (item: ListResult) => {
      switch (item.exam_status) {
        case examStatus.ESFinish:
          notification.open({
            message: '提醒',
            description: '该场考试已经结束',
          })
          return
        case examStatus.ESRunning:
          router.push({
            name: 'examMissionSelector',
            params: { exam: item.id },
          })
          return
        case examStatus.ESNotStart:
          Modal.confirm({
            title: '是否开始考试',
            icon: createVNode(ExclamationCircleOutlined),
            content: '如果点击确定就会开始考试计时，请谨慎考虑！',
            okText: '我确定',
            okType: 'danger',
            cancelText: '考虑一下',
            onOk() {
              return new Promise((resolve, reject) => {
                new exam()
                  .getExamRunningInfo()
                  .then((res: examRunningInfo) => {
                    if (res === undefined) {
                      resolve(res)
                      return
                    } else {
                      jumpToRunningExam(res)
                      reject(res)
                    }
                  })
                  .catch((err) => {
                    reject(err)
                  })
              })
                .then(() => {
                  router.push({
                    name: 'examMissionSelector',
                    params: { exam: item.id },
                  })
                })
                .catch((err) => {})
            },
            onCancel() {},
          })
          return
        case examStatus.ESPassTime:
          notification.open({
            message: '提醒',
            description: '该场考试未开始或已经结束',
          })
          return
      }
    }

    // 跳转至正在进行的考试
    const jumpToRunningExam = (data: examRunningInfo) => {
      Modal.confirm({
        title: '您有其他考试正在进行中',
        icon: createVNode(ExclamationCircleOutlined),
        content: `您【${data.exam_name}】考试正在进行`,
        okText: '返回正在进行的考试',
        cancelText: '返回考试选择',
        onOk() {
          router.push({
            name: 'examMissionSelector',
            params: { exam: data.exam_id },
          })
        },
        onCancel() {},
      })
    }

    return {
      isProjectDataLoading: false,
      departmentLessonData,
      numberCreatorFn,
      examDescCreator,
      examStartButtonFn,
      examStatus,
      GetExamButtonDesc,
    }
  },
}
</script>

<style>
.examSelector {
  width: 100%;
}

.ant-list-item-meta-description {
  white-space: pre-wrap;
}
</style>