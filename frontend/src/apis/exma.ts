import { useRequest } from 'vue-request'
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
type examRunningInfo = {
  account: number
  exam_id: number
  exam_name: string
  left_time: number
}

// 获取考试运行时
const getExamRunningInfo = new Promise((resolve, reject) => {
  defaultClient
    .get<BaseResponse>(apiPath.check, {})
    .then((res) => {
      resolve(<examRunningInfo>res.data.Data)
    })
    .catch((err) => {
      reject(err)
    })
})

// 获取考试运行时
export { examRunningInfo, getExamRunningInfo }
