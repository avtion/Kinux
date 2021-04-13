<template>
  <!-- 班级考试成绩 -->
  <div class="w-full">
    <!-- 表格 -->
    <a-table :dataSource="dataSource" :columns="columns" :pagination="false">
      <template #time="{ text }">
        <span>{{
          text === 0 ? '未开始' : moment.unix(text).format('lll')
        }}</span>
      </template>
    </a-table>
  </div>
</template>

<script lang="ts" type="module">
import { Score, ExamScoreForAdmin } from '@/apis/score'
import { defineComponent, ref, watch, reactive } from 'vue'
import { useRequest } from 'vue-request'

import { moment } from '@/utils/time'

export default defineComponent({
  props: {
    dp: { type: Number, default: 1 },
    lesson: { type: Number, default: 1 },
    exam: { type: Number, default: 1 },
    isSaveMode: { type: Boolean, default: false },
  },
  setup(props) {
    console.log(props.lesson, props.exam, props.isSaveMode)
    const score = new Score(props.dp, props.lesson, props.exam)
    const scoreData = reactive<ExamScoreForAdmin[]>([])
    const { run: getScoreData } = useRequest(score.GetExamScoreForAdmin, {
      formatResult: (res) => {
        res.forEach((v) => {
          scoreData.push(reactive(v))
        })
        console.log(scoreData)
      },
    })
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
          title: '考试开始时间',
          dataIndex: 'exam_begin_at',
          key: 'exam_begin_at',
          slots: { customRender: 'time' },
        },
        {
          title: '考试结束时间',
          dataIndex: 'exam_begin_at',
          key: 'exam_begin_at',
          slots: { customRender: 'time' },
        },
      ],
    }
  },
})
</script>

<style>
</style>