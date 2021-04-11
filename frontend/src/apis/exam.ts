import { BaseResponse, defaultClient } from '@api/request'

import { ref } from 'vue'
// 考试状态
export enum examStatus {
  ESNotStart = 1, // 未开始
  ESRunning, // 正在进行
  ESFinish, // 考试结束
  ESPassTime, // 考试未开始或已经结束
}

// API接口
const apiPath = {
  list: '/v1/exam/dp/',
  check: '/v1/exam/check/',
  start: '/v1/exam/start/',
}

// 考试进行时
export type examRunningInfo = {
  account: number
  exam_id: number
  exam_name: string
  left_time: number
}

// 考试对象
export class exam {
  // 获取实验运行情况
  getExamRunningInfo(): Promise<examRunningInfo> {
    return new Promise((resolve: (value: examRunningInfo) => void, reject) => {
      defaultClient
        .get<BaseResponse>(apiPath.check, {})
        .then((res) => {
          resolve(<examRunningInfo>res.data.Data)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }

  // 开始实验
  startExam(examID: number) {
    return new Promise((resolve, reject) => {
      defaultClient
        .get(apiPath.start, {
          params: {
            exam: examID,
          },
        })
        .then((res) => {
          resolve(res)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }
}

// 考试进行状态
export const examInfo = ref<examRunningInfo>()
