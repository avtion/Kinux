<template>
  <a-card title="考试实验" :bordered="false" :loading="isListDataLoading">
    <a-list item-layout="horizontal" :data-source="missionData">
      <template #renderItem="{ item, index }">
        <a-list-item>
          <!-- 元数据 -->
          <a-list-item-meta :description="item.desc">
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
              :type="GetMissionButtonType(item.status)"
              :loading="GetMissionButtonLoadingStatus(item.status)"
              @click="MissionHandler(index, item)"
              :disabled="item.status == missionStatus.Done"
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
import { reactive, ref, inject, onMounted } from 'vue'

// apis
import { mission, missionList, missionStatus } from '@api/mission'

// 图标生成
import Avatars from '@dicebear/avatars'
import AvatarsSprites from '@dicebear/avatars-male-sprites'
import sprites from '@dicebear/avatars-initials-sprites'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'
import { Profile } from '@/store/interfaces'

// websocket
import {
  WebSocketConn,
  WebsocketMessage,
  WebsocketOperation,
} from '@/utils/websocketConn'

// vue-request
import { useRequest } from 'vue-request'
import { BaseResponse, defaultClient } from '@/apis/request'

export default {
  setup(props, ctx) {
    const router = useRouter()

    // 获取课程参数
    const examID = Number(router.currentRoute.value.params.exam)
    console.log(examID)

    // 从上下文中获取对象
    const ws: WebSocketConn = inject<WebSocketConn>('websocket')

    // 加载任务数据
    interface missionReqParams {
      page: number
      size: number
      lesson: number
    }
    const departmentLessonDataAPI = (params: missionReqParams) => {
      return defaultClient.get<BaseResponse>('/v2/ms/', {
        params: params,
      })
    }
    const getListParams = (): missionReqParams => {
      return <missionReqParams>{
        page: 0,
        size: 0,
        lesson: examID,
      }
    }
    const { data: missionData, loading: isListDataLoading } = useRequest(
      departmentLessonDataAPI,
      {
        defaultParams: [getListParams()],
        formatResult: (res): missionList[] => {
          return res.data.Data
        },
      }
    )

    // 任务处理函数
    const MissionHandler = (index: number, m: missionList) => {
      console.log(m)
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
            params: { mission: m.id, lesson: examID },
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
    })
    const numberCreatorFn = (str: any): string => {
      return numberCreator.create(str + '')
    }

    return {
      MissionHandler,
      missionStatus,
      GetMissionButtonType,
      GetMissionButtonLoadingStatus,
      GetMissionButtonDesc,
      numberCreatorFn,
      isListDataLoading,
      missionData,
    }
  },
}

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
    default:
      return 'dashed'
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
    default:
      return ''
  }
}
</script>

<style>
</style>