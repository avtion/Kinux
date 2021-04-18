<template>
  <div class="w-full h-full">
    <div class="w-full h-full p-5">
      <div class="w-full bg-white rounded p-5">
        <!-- 顶部 -->
        <a-page-header title="考试成绩" sub-title="" @back="back()">
          <a-divider>汇总</a-divider>
          <!-- 分数和考点 -->
          <a-row type="flex" justify="space-around">
            <a-col :span="4">
              <!-- 总分 -->
              <a-statistic title="成绩" :value="examScore.score" />
            </a-col>
            <a-col :span="4">
              <!-- 总考点数量 -->
              <a-statistic
                title="开始时间"
                :value="
                  moment.unix(examScore.exam_begin_at).format('MM-DD HH:mm')
                "
              />
            </a-col>
            <a-col :span="4">
              <!-- 已完成考点 -->
              <a-statistic
                title="结束时间"
                :value="
                  moment.unix(examScore.exam_end_at).format('MM-DD HH:mm')
                "
              />
            </a-col>
          </a-row>
          <!-- 表格 -->
        </a-page-header>
        <a-divider>统计</a-divider>
        <div
          class="pl-10 pr-10 pt-10 h-1/4 bg-gray-50 rounded"
          v-if="examScore.mission_scores.length > 0"
        >
          <a-row>
            <a-col :span="12">
              <liquid-chart v-bind="liquidChartConfig" />
            </a-col>
            <a-col :span="12">
              <pie-chart v-bind="pieChartConfig" />
            </a-col>
          </a-row>
        </div>
        <a-empty description="暂无数据" v-else />
        <a-divider>本场考试实验</a-divider>
        <div v-if="examScore.mission_scores.length > 0" class="pr-5 pl-5">
          <a-list
            item-layout="horizontal"
            :data-source="examScore.mission_scores"
          >
            <template #renderItem="{ item, index }">
              <a-list-item>
                <!-- 元数据 -->
                <a-list-item-meta
                  :description="`成绩占比: ${item.percent}% \n考点描述: ${item.mission_desc}`"
                >
                  <!-- 标题 -->
                  <template #title>
                    <div>
                      <span class="font-bold mr-5">{{
                        item.mission_name
                      }}</span>
                      <div class="w-2/4 inline-block">
                        <a-progress
                          :percent="
                            (item.finish_score_counter /
                              item.all_score_counter) *
                            100
                          "
                          :showInfo="false"
                          size="small"
                        />
                      </div>
                    </div>
                  </template>
                  <!-- 头像 -->
                  <template #avatar>
                    <a-avatar :src="numberCreatorFn(index + 1)" />
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>
        </div>
        <a-empty description="暂无数据" v-else />
      </div>
    </div>
  </div>
</template>

<script lang="ts" type="module">
import { Score, ExamScore, missionScore } from '@/apis/score'
import { defineComponent, reactive, watch } from 'vue'

import { moment } from '@/utils/time'

// vue-router
import { useRouter } from 'vue-router'

// 图标生成
import Avatars from '@dicebear/avatars'
import sprites from '@dicebear/avatars-initials-sprites'

// 折线图
import { LiquidChart, PieChart } from '@opd/g2plot-vue'

export default defineComponent({
  components: {
    LiquidChart,
    PieChart,
  },
  props: {
    lessonID: String,
    examID: String,
  },
  setup(props) {
    console.log(props.lessonID, props.examID)
    // 获取路由
    const router = useRouter()

    // 序号
    const numberCreator = new Avatars(sprites, {
      dataUri: true,
      background: '#60A5FA',
    })
    const numberCreatorFn = (str: any): string => {
      return numberCreator.create(str + '')
    }

    // 返回
    const back = () => {
      router.push({ name: 'stuScore' })
    }

    const examScore = reactive<ExamScore>({
      exam_id: 0,
      exam_name: '',
      exam_desc: '',
      exam_begin_at: 0,
      exam_end_at: 0,
      score: 0,
      mission_scores: reactive(<missionScore[]>[]),
      total: 0,
    })
    const score = new Score(0, Number(props.lessonID), Number(props.examID), 0)
    score.GetExamScore().then((res) => {
      examScore.exam_id = res.exam_id
      examScore.exam_name = res.exam_name
      examScore.exam_desc = res.exam_desc
      examScore.exam_begin_at = res.exam_begin_at
      examScore.exam_end_at = res.exam_end_at
      examScore.mission_scores = reactive(res.mission_scores)
      examScore.score = res.score
      if (res.total !== 0) {
        liquidChartConfig.percent = res.score / res.total
      }
      if (res.mission_scores.length > 0) {
        const finishCounter = { type: '已完成考点', value: 0 }
        const unFinishCounter = { type: '未完成考点', value: 0 }
        res.mission_scores.forEach((v) => {
          v.score_details.forEach((detail) => {
            if (detail.is_finish) {
              finishCounter.value++
            } else {
              unFinishCounter.value++
            }
          })
        })
        pieChartConfig.data = [finishCounter, unFinishCounter]
      }
      console.log(res)
    })

    // 水波图
    const liquidChartConfig = reactive({
      smooth: true,
      autoFit: true,
      percent: 0,
      height: 200,
      outline: {
        border: 4,
        distance: 8,
      },
      wave: {
        length: 128,
      },
    })

    // 饼状图
    const pieChartConfig = reactive({
      data: [
        { type: '已完成考点', value: 0 },
        { type: '未完成考点', value: 0 },
      ],
      height: 200,
      angleField: 'value',
      colorField: 'type',
      radius: 0.9,
      label: {
        type: 'inner',
        content: ({ percent }) => `${(percent * 100).toFixed(0)}%`,
        style: {
          fontSize: 14,
          textAlign: 'center',
        },
      },
      interactions: [{ type: 'element-active' }],
    })

    return {
      back,
      // 时间处理
      moment,

      // 考试成绩
      examScore,

      // 序号生成
      numberCreatorFn,

      // 水波图
      liquidChartConfig,

      // 饼状图
      pieChartConfig,
    }
  },
})
</script>

<style>
</style>