<template>
  <a-card title="考试实验" :bordered="false" :loading="isListDataLoading">
    <a-list item-layout="horizontal" :data-source="missionData">
      <template #renderItem="{ item, index }">
        <a-list-item>
          <!-- 元数据 -->
          <a-list-item-meta :description="descCreator(item)">
            <!-- 标题 -->
            <template #title>
              <a>{{ item.mission_name }}</a>
            </template>
            <!-- 头像 -->
            <template #avatar>
              <a-avatar :src="numberCreatorFn(index + 1)" />
            </template>
          </a-list-item-meta>
          <!-- 操作 -->
          <template #actions>
            <a-button
              :type="GetMissionButtonType(item.status)"
              :loading="GetMissionButtonLoadingStatus(item.status)"
              @click="MissionHandler(index, item)"
              :disabled="
                item.status == missionStatus.Done ||
                item.status == missionStatus.Block
              "
            >
              {{ GetMissionButtonDesc(item.status) }}
            </a-button>
          </template>
        </a-list-item>
      </template>
    </a-list>
  </a-card>
</template>

<script lang="ts" type="module">
// vue
import { inject, onMounted, defineComponent } from 'vue'

// apis
import { missionStatus } from '@api/mission'

// 图标生成
import Avatars from '@dicebear/avatars'
import sprites from '@dicebear/avatars-initials-sprites'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

// websocket
import {
  WebSocketConn,
  WebsocketMessage,
  WebsocketOperation,
} from '@/utils/websocketConn'

// vue-request
import { useRequest } from 'vue-request'
import { BaseResponse, defaultClient } from '@/apis/request'
import { exam } from '@/apis/exam'

export default defineComponent({
  setup(props, ctx) {
    const router = useRouter()

    // 获取课程参数
    const examID = Number(router.currentRoute.value.params.exam)

    // 从上下文中获取对象
    const ws: WebSocketConn = inject<WebSocketConn>('websocket')

    // 加载任务数据
    interface examReqParams {
      exam: number
    }
    type examMission = {
      exam_id: number
      id: number
      mission_id: number
      mission_name: string
      percent: number
      priority: number
      status: number
    }
    const departmentLessonDataAPI = (params: examReqParams) => {
      return defaultClient.get<BaseResponse>('/v1/em/list/', {
        params: params,
      })
    }
    const getListParams = (): examReqParams => {
      return <examReqParams>{
        exam: examID,
      }
    }
    const { data: missionData, loading: isListDataLoading } = useRequest(
      departmentLessonDataAPI,
      {
        defaultParams: [getListParams()],
        formatResult: (res): examMission[] => {
          return res.data.Data
        },
      }
    )

    // 任务处理函数
    const MissionHandler = (index: number, m: examMission) => {
      const status = m.status
      switch (status) {
        case missionStatus.Stop:
          startMission(index, m.id + '')
          return
        case missionStatus.Pending:
          return
        case missionStatus.Working:
          router.push({
            name: 'shell',
            params: { mission: m.id, lesson: m.mission_id, exam: examID },
          })
          return
        case missionStatus.Done:
          return
        default:
          return
      }
    }

    // 启动任务
    const startMission = (missionListIndex: number, missionID: string) => {
      missionData.value[missionListIndex].status = missionStatus.Pending
      const msg: WebsocketMessage = {
        op: WebsocketOperation.MissionApply,
        data: {
          id: missionID,
          exam: examID + '',
        },
      }
      const fn = (ws: WebSocketConn) => {
        ws.sendWithCallback(
          JSON.stringify(msg),
          WebsocketOperation.ContainersDone,
          (_ws: WebSocketConn): void => {
            missionData.value[missionListIndex].status = missionStatus.Working
          },
          true
        )
      }
      if (ws.readyState !== WebSocket.OPEN) {
        ws.waitQueue.push(fn)
      } else {
        fn(ws)
      }
    }

    // 序号
    const numberCreator = new Avatars(sprites, {
      dataUri: true,
      background: '#3B82F6',
    })
    const numberCreatorFn = (str: any): string => {
      return numberCreator.create(str + '')
    }

    // 描述生成
    const descCreator = (item: examMission): string => {
      return `> 成绩占比: ${item.percent}%`
    }

    // 启动实验
    new exam().startExam(examID).then((res) => {
      console.log(res)
    })

    return {
      MissionHandler,
      missionStatus,
      GetMissionButtonType,
      GetMissionButtonLoadingStatus,
      GetMissionButtonDesc,
      numberCreatorFn,
      isListDataLoading,
      missionData,
      descCreator,
    }
  },
})

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
    case missionStatus.Block:
      return 'default'
    default:
      return 'default'
  }
}

// 获取任务按钮的状态
function GetMissionButtonLoadingStatus(t: number): boolean {
  return t == missionStatus.Pending
}

// 获取任务按钮的描述
function GetMissionButtonDesc(t: number): string {
  switch (t) {
    case missionStatus.Stop:
      return '开始'
    case missionStatus.Pending:
      return '正在加载'
    case missionStatus.Working:
      return '进入终端'
    case missionStatus.Done:
      return '已完成'
    case missionStatus.Block:
      return '需完成之前实验才能开始'
    default:
      return ''
  }
}
</script>

<style>
</style>