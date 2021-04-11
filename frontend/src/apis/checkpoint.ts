import { BaseResponse, defaultClient } from '@api/request'
import { AxiosResponse } from 'axios'

export type missionCheckpointMetaData = {
  id: number // 实验考点ID
  mission: number // 实验ID
  check_point: number // 考点ID
  percent: number // 成绩占比
  priority: number // 权重
  target_container: string // 目标容器
  checkpoint_id: number // 考点ID
  cp_name: string // 考点名称
  cp_desc: string // 考点描述
  cp_command: string // 考点检查的指令
  cp_method: number // 考点方法
  is_done: boolean // 考点是否已经完成
}

// 考点结果结构
export type missionCheckpointRes = {
  container_name: string
  data: missionCheckpointMetaData[]
}

// 考点
export class Checkpoint {
  private lessonID: number // 课程ID
  private missionID: number // 实验ID
  private examID: number // 考试ID
  constructor(lessonID: number, missionID: number, examID: number) {
    this.lessonID = lessonID
    this.missionID = missionID
    this.examID = examID
  }
  // 获取实验的检查点
  Get = () => {
    return new Promise<missionCheckpointRes[]>((resolve, reject) => {
      defaultClient
        .get('/v1/cp/mcp/', {
          params: {
            lesson: this.lessonID,
            exam: this.examID,
            mission: this.missionID,
          },
        })
        .then((res: AxiosResponse<BaseResponse>) => {
          resolve(<missionCheckpointRes[]>res.data.Data)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }
}
