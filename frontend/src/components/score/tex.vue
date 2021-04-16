<template>
  <!-- 班级考试成绩 -->
  <div class="w-full">
    <!-- 表格 -->
    <a-table
      :dataSource="dataSource"
      :columns="columns"
      :pagination="false"
      rowKey="id"
    >
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
    dp: { type: Number, default: 0 },
    lesson: { type: Number, default: 0 },
    exam: { type: Number, default: 0 },
    isSaveMode: { type: Boolean, default: false },
  },
  setup(props) {
    console.log(props.lesson, props.exam, props.isSaveMode)
    const score = new Score(props.dp, props.lesson, props.exam)
    const scoreData = reactive<ExamScoreForAdmin[]>([])

    // 兼容存档
    let client = score.GetExamScoreForAdmin
    if (props.isSaveMode) {
      client = () => {
        return new Promise<ExamScoreForAdmin[]>((resolve, reject) => {
          score
            .GetSaveScore(2, props.exam)
            .then((res) => {
              resolve(<ExamScoreForAdmin[]>res)
            })
            .catch((err) => {
              reject(err)
            })
        })
      }
    }
    const { run: getScoreData } = useRequest(client, {
      formatResult: (res: ExamScoreForAdmin[]) => {
        res.forEach((v) => {
          scoreData.push(reactive(v))
        })
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
          slots: { customRender: 'time' },
        },
        {
          title: '考试结束时间',
          dataIndex: 'exam_begin_at',
          slots: { customRender: 'time' },
        },
      ],
    }
  },
})
</script>

<style>
</style>