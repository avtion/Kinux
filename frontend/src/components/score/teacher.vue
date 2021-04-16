<template>
  <!-- 老师查询成绩 -->

  <div class="w-full h-full">
    <div class="w-full h-full p-5">
      <div class="w-full h-full bg-white rounded">
        <!-- 表单 -->
        <div class="pt-5 pl-5">
          <a-space>
            <a-select
              placeholder="班级"
              style="width: 150px"
              v-model:value="dp"
            >
              <a-select-option
                v-for="(item, index) in dpList"
                :key="index"
                :value="item.id"
              >
                {{ item.name }}
              </a-select-option>
            </a-select>
            <!-- 课程 -->
            <a-select
              placeholder="课程"
              style="width: 150px"
              v-model:value="lesson"
            >
              <a-select-option
                v-for="(item, index) in lessonList"
                :key="index"
                :value="item.id"
              >
                {{ item.name }}
              </a-select-option>
            </a-select>
            <!-- 存档类型 -->
            <a-select placeholder="记录类型" v-model:value="recordType">
              <a-select-option :value="recordTypes.now">实时</a-select-option>
              <a-select-option :value="recordTypes.save">存档</a-select-option>
            </a-select>
            <!-- 实验或者考试 -->
            <a-select placeholder="实验/考试" v-model:value="scoreType">
              <a-select-option :value="missionOrExam.mission"
                >实验</a-select-option
              >
              <a-select-option :value="missionOrExam.exam"
                >考试</a-select-option
              >
            </a-select>
            <!-- 查询目标 -->
            <a-select
              placeholder="查询目标"
              style="min-width: 250px"
              v-model:value="targetID"
            >
              <a-select-option
                v-for="(item, index) in targets"
                :key="index"
                :value="item.id"
              >
                {{ item.name }}
              </a-select-option>
            </a-select>
          </a-space>
        </div>
        <div class="pt-5 pl-5">
          <a-space>
            <!-- 下载成绩 -->
            <a-button
              type="primary"
              :disabled="!isExamScoreShow && !isMissionScoreShow"
            >
              下载成绩
            </a-button>
            <!-- 存档 -->
            <a-button
              type="default"
              :disabled="!isShowSaveButton"
              :loading="isSaveButtonLoading"
              @click="saveScore"
              >存档</a-button
            >
          </a-space>
        </div>
        <a-divider>成绩</a-divider>
        <a-empty
          :description="false"
          v-show="!isExamScoreShow && !isMissionScoreShow"
        />
        <!-- 考试成绩 -->
        <div class="pl-5 pr-5">
          <tex
            v-if="isExamScoreShow"
            :dp="dp"
            :lesson="lesson"
            :exam="targetID"
            :isSaveMode="recordType === recordTypes.save"
          ></tex>
        </div>
        <!-- 实验成绩查询 -->
        <div class="pl-5 pr-5">
          <tms
            v-if="isMissionScoreShow"
            :dp="dp"
            :lesson="lesson"
            :mission="targetID"
            :isSaveMode="recordType === recordTypes.save"
          ></tms>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" type="module">
import { defineComponent, ref, watch, computed } from 'vue'

import tex from '@/components/score/tex.vue'
import tms from '@/components/score/tms.vue'

enum recordTypes {
  now = 1,
  save,
}

enum missionOrExam {
  mission = 1,
  exam,
}

// vue-request
import { useRequest } from 'vue-request'
import { BaseResponse, defaultClient, paths } from '@/apis/request'
import { Score } from '@/apis/score'

// 时间处理
import { moment } from '@/utils/time'

type dpListResult = {
  id: number
  name: string
  creat_at: string
  updated_at: string
}

type lessonListResut = {
  id: number
  name: string
  desc: string
}

// 目标内容
type targetResult = {
  id: number
  name: string
}

export default defineComponent({
  components: {
    tex,
    tms,
  },
  setup(props) {
    const dp = ref<number>(),
      lesson = ref<number>(),
      recordType = ref<recordTypes>(recordTypes.now),
      scoreType = ref<missionOrExam>(missionOrExam.mission),
      targetID = ref<number>(),
      targets = ref<targetResult[]>([])

    // 班级列表
    const { data: dpList } = useRequest(
      () => {
        return defaultClient.get<BaseResponse>(paths.department.list)
      },
      {
        formatResult: (res): dpListResult[] => {
          return res.data.Data
        },
      }
    )
    watch(dp, (newValue) => {
      getLessonList(newValue)
    })

    // 实验API
    const { run: getMissionData } = useRequest(
      (lesson: number) => {
        return defaultClient.get<BaseResponse>('/v2/lm/list', {
          params: {
            lesson: lesson,
          },
        })
      },
      {
        formatResult: (res): targetResult[] => {
          type listResult = {
            id: number
            mission_id: number
            mission_name: string
            mission_desc: string
            priority: number
          }
          const data = <listResult[]>res.data.Data
          const _res = <targetResult[]>[]
          data.forEach((v) => {
            _res.push(<targetResult>{
              id: v.id,
              name: v.mission_name,
            })
          })
          targets.value = _res
          return _res
        },
        manual: true,
      }
    )

    // 考试
    const { run: getExamListData } = useRequest(
      (lesson: number) => {
        return defaultClient.get<BaseResponse>('/v1/exam/list/', {
          params: {
            lesson: lesson,
          },
        })
      },
      {
        formatResult: (res): targetResult[] => {
          type examListResult = {
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
          }
          const _res = <examListResult[]>res.data.Data
          const data = <targetResult[]>[]
          _res.forEach((v) => {
            data.push({
              id: v.id,
              name: v.name,
            })
          })
          targets.value = _res
          return data
        },
        manual: true,
      }
    )

    // 课程
    const { data: lessonList, run: getLessonList } = useRequest(
      (dpID: number) => {
        return defaultClient.get<BaseResponse>('/v2/dl/list', {
          params: {
            department: dp.value,
          },
        })
      },
      {
        formatResult: (res): lessonListResut => {
          return res.data.Data
        },
        manual: true,
      }
    )

    // 存档
    const { run: getSaverList } = useRequest(
      () => {
        return defaultClient.get<BaseResponse>('/v2/score/quick/', {
          params: {
            dp: dp.value,
            lesson: lesson.value,
            page: 0,
            size: 0,
            type: scoreType.value,
          },
        })
      },
      {
        manual: true,
        formatResult: (res) => {
          console.log(res)
          type __listRes = {
            id: number
            raw_id: number
            raw_name: string
            created_at: string
          }
          const _res = <__listRes[]>res.data.Data
          const data = <targetResult[]>[]
          _res.forEach((v) => {
            data.push(<targetResult>{
              id: v.id,
              name: `${v.raw_name}-${moment(v.created_at).format('lll')}`,
            })
          })
          targets.value = data
        },
      }
    )

    // 监听筛选框
    watch([dp, lesson, recordType, scoreType], () => {
      if (lesson.value === undefined || lesson.value == 0) {
        return
      }
      // 置空
      targets.value = []
      targetID.value = undefined

      // 根据不同情况获取目标
      switch (recordType.value) {
        case recordTypes.now:
          // 实时
          switch (scoreType.value) {
            case missionOrExam.mission:
              getMissionData(lesson.value)
              // 实验
              break
            case missionOrExam.exam:
              getExamListData(lesson.value)
              // 考试
              break
          }
          break
        case recordTypes.save:
          getSaverList()
          // 存档
          break
      }
    })

    // 监听目标内容
    watch([dp, lesson, recordType, scoreType, targetID], () => {
      if (targetID.value === undefined || targetID.value === 0) {
        isExamScoreShow.value = false
        isMissionScoreShow.value = false
        return
      }
      switch (scoreType.value) {
        case missionOrExam.exam:
          isExamScoreShow.value = true
          break
        case missionOrExam.mission:
          isMissionScoreShow.value = true
          break
        default:
      }
    })

    // 是否显示存档按钮
    const isShowSaveButton = computed(() => {
      return (
        recordType.value == recordTypes.now &&
        targetID.value !== undefined &&
        targetID.value !== 0
      )
    })

    // 成绩存档
    const isSaveButtonLoading = ref<boolean>(false)

    // 存档成绩
    const saveScore = () => {
      new Score(dp.value, lesson.value).SaveScore(
        scoreType.value,
        targetID.value
      )
    }

    // 是否显示考试成绩
    const isExamScoreShow = ref<boolean>(false)
    // 是否显示考试成绩
    const isMissionScoreShow = ref<boolean>(false)

    return {
      // 记录类型
      recordTypes,

      // 实验OR考试
      missionOrExam,

      // 班级列表
      dpList,
      dp,

      // 课程列表
      lessonList,
      lesson,

      // 记录类型
      recordType,
      // 成绩类型
      scoreType,

      // 目标记录
      targets,
      targetID,

      // 是否显示存档按钮
      isShowSaveButton,

      // 存档按钮是否正在加载
      isSaveButtonLoading,

      // 存档成绩
      saveScore,

      // 是否显示考试成绩
      isExamScoreShow,
      // 是否显示实验成绩
      isMissionScoreShow,
    }
  },
})
</script>

<style lang="less" scoped>
</style>