import { BaseResponse, defaultClient } from '@api/request'

// 考试状态
export enum examStatus {
  ESNotStart = 1, // 未开始
  ESRunning, // 正在进行
  ESFinish, // 考试结束
}

// API接口
const apiPath = {
  list: '/v1/exam/dp/',
  check: '/v1/exam/check/',
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
  startExam(): Promise<void> {
    return new Promise((resolve, reject) => {})
  }
}
