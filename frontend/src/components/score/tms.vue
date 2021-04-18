<template>
  <!-- 班级实验成绩 -->
  <div class="w-full">
    <!-- 表格 -->
    <a-table
      :dataSource="dataSource"
      :columns="columns"
      :pagination="false"
      rowKey="id"
    >
      <template #percent="{ record }">
        <a-progress
          :percent="
            (record.finish_score_counter / record.all_score_counter) * 100
          "
          :showInfo="false"
          size="small"
        />
      </template>
    </a-table>
  </div>
</template>

<script lang="ts" type="module">
import { Score, MissionScoreForAdmin } from '@/apis/score'
import { defineComponent, reactive, watch } from 'vue'
import { useRequest } from 'vue-request'

import { moment } from '@/utils/time'

export default defineComponent({
  props: {
    dp: { type: Number, default: 0 },
    lesson: { type: Number, default: 0 },
    mission: { type: Number, default: 0 },
    isSaveMode: { type: Boolean, default: false },
  },
  setup(props) {
    console.log(props.lesson, props.mission, props.isSaveMode)
    const score = new Score(props.dp, props.lesson, 0, props.mission)
    const scoreData = reactive<MissionScoreForAdmin[]>([])

    // 兼容存档
    let client = score.GetMissionScoreForAdmin
    if (props.isSaveMode) {
      client = () => {
        return new Promise<MissionScoreForAdmin[]>((resolve, reject) => {
          score
            .GetSaveScore(1, props.mission)
            .then((res) => {
              resolve(<MissionScoreForAdmin[]>res)
            })
            .catch((err) => {
              reject(err)
            })
        })
      }
    }
    const { run: getScoreData } = useRequest(client, {
      formatResult: (res) => {
        res.forEach((v) => {
          scoreData.push(reactive(v))
        })
        console.log(scoreData)
      },
    })

    // TODO 修复无法实时查询成绩

    return {
      // 时间处理
      moment,
      // 数据源
      dataSource: scoreData,

      columns: [
        {
          title: '排名',
          dataIndex: 'pos',
          key: 'pos',
        },
        {
          title: '用户名',
          dataIndex: 'username',
          key: 'username',
        },
        {
          title: '真实姓名',
          dataIndex: 'real_name',
          key: 'real_name',
        },
        {
          title: '成绩',
          dataIndex: 'score',
          key: 'score',
        },
        {
          title: '考点完成度',
          dataIndex: 'all_score_counter',
          slots: { customRender: 'percent' },
        },
      ],
    }
  },
})
</script>

<style>
</style>