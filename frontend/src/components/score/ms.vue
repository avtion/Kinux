<template>
  <div class="w-full h-full">
    <div class="w-full h-full p-5">
      <div class="w-full bg-white rounded">
        <!-- 顶部 -->
        <a-page-header title="实验成绩" :sub-title="title" @back="back()">
          <a-divider>汇总</a-divider>
          <!-- 分数和考点 -->
          <a-row type="flex" justify="space-around">
            <a-col :span="4">
              <!-- 总分 -->
              <a-statistic title="成绩" :value="total" />
            </a-col>
            <a-col :span="4">
              <!-- 总考点数量 -->
              <a-statistic title="总考点数量" :value="allCps" />
            </a-col>
            <a-col :span="4">
              <!-- 已完成考点 -->
              <a-statistic title="已完成考点" :value="doneCps" />
            </a-col>
            <a-col :span="4">
              <!-- 未完成考点 -->
              <a-statistic title="未完成考点" :value="restCps" />
            </a-col>
          </a-row>
          <!-- 表格 -->
        </a-page-header>
        <a-divider>统计</a-divider>
        <!-- 统计数据 -->
        <div
          class="pl-10 pr-10 pt-10 bg-gray-50 rounded"
          v-if="areaChartConfig.data.length > 0"
        >
          <a-row>
            <a-col :span="12"
              ><liquid-chart v-bind="liquidChartConfig"
            /></a-col>
            <a-col :span="12"><area-chart v-bind="areaChartConfig" /></a-col>
          </a-row>
        </div>

        <a-empty description="暂无数据" v-else />
        <a-divider>考点完成明细</a-divider>
        <div v-if="scoreItems.length > 0" class="pr-5 pl-5">
          <a-list item-layout="horizontal" :data-source="scoreItems">
            <template #renderItem="{ item, index }">
              <a-list-item>
                <!-- 元数据 -->
                <a-list-item-meta
                  :description="`成绩占比: ${item.percent}% \n考点描述: ${item.checkpoint_desc}`"
                >
                  <!-- 标题 -->
                  <template #title>
                    <div>
                      <span class="font-bold mr-5">{{
                        item.checkpoint_name
                      }}</span>
                      <a-tag color="cyan">{{ item.target_container }}</a-tag>
                      <a-tag :color="item.is_finish ? 'success' : 'warning'">{{
                        item.is_finish ? '已完成' : '未完成'
                      }}</a-tag>
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
import { defineComponent, ref, computed, reactive } from 'vue'

import { Score, missionScore, ScoreDetail } from '@api/score'
// vue-request
import { useRequest } from 'vue-request'

// 图标生成
import { IntCreator } from '@/utils/avatar'

// vue-router
import { useRouter } from 'vue-router'

// 折线图
import { LiquidChart, AreaChart } from '@opd/g2plot-vue'

import { moment } from '@/utils/time'

export default defineComponent({
  components: {
    LiquidChart,
    AreaChart,
  },
  props: {
    lessonID: {
      type: String,
    },
    missionID: {
      type: String,
    },
  },
  setup(props, ctx) {
    const router = useRouter()

    // 小标题实验名称
    const title = ref<string>('')
    // 分数
    const total = ref<number>(0)
    const allCps = ref<number>(0)
    const doneCps = ref<number>(0)
    const restCps = computed(() => {
      return allCps.value - doneCps.value
    })

    // 序号
    const numberCreatorFn = (str: any): string => {
      return IntCreator(str + '', '#60A5FA')
    }

    // 加载成绩
    const score = new Score(
      0,
      Number(props.lessonID),
      0,
      Number(props.missionID)
    )
    const scoreItems = ref<ScoreDetail[]>([])
    const { data: missionScoreData } = useRequest(score.GetMissionScore, {
      formatResult: (res: missionScore): missionScore => {
        total.value = res.score
        allCps.value = res.all_score_counter
        doneCps.value = res.finish_score_counter
        title.value = res.mission_name
        scoreItems.value = res.score_details
        console.log(res)

        // 计算统计数值
        let counter = 1
        res.score_details.forEach((v) => {
          if (v.is_finish) {
            areaChartConfig.data.push({
              完成时间: moment.unix(v.finish_time).format('MM-DD HH:mm:ss'),
              完成考点数量: counter,
            })
            counter++
          }
        })
        if (allCps.value !== 0) {
          liquidChartConfig.percent = res.score / res.total
        }
        return res
      },
    })

    // 返回
    const back = () => {
      router.push({ name: 'stuScore' })
    }

    // 折线图
    const areaChartConfig = reactive({
      smooth: true,
      autoFit: true,
      xField: '完成时间',
      yField: '完成考点数量',
      height: 200,
      data: [],
      areaStyle: () => {
        return {
          fill: 'l(270) 0:#ffffff 0.5:#7ec2f3 1:#1890ff',
        }
      },
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

    return {
      // 退回成绩查询
      back,
      // 序号
      numberCreatorFn,

      // 实验成绩
      missionScoreData,
      scoreItems,

      // 实验名称
      title,

      // 总分
      total,

      // 考点数量
      allCps,
      doneCps,
      restCps,

      // 折线图
      areaChartConfig,
      // 水波图
      liquidChartConfig,
    }
  },
})
</script>

<style>
</style>