import { BaseResponse, defaultClient } from '@api/request'
import { AxiosResponse } from 'axios'

export class Score {
  public department: number
  public exam: number
  public lesson: number
  public mission: number

  constructor(department = 0, lesson = 0, exam = 0, mission = 0) {
    this.department = department
    this.exam = exam
    this.lesson = lesson
    this.mission = mission
  }

  // 获取实验成绩
  GetMissionScore = () => {
    return new Promise((resolve: (res: missionScore) => void, reject) => {
      defaultClient
        .get('/v2/score/mission/', {
          params: {
            dp: this.department,
            lesson: this.lesson,
            mission: this.mission,
          },
        })
        .then((res: AxiosResponse<BaseResponse>) => {
          resolve(res.data.Data)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }

  //   GetExamScore(): Promise<> {
  //     return new Promise((resolve, reject) => {})
  //   }
}

export type ScoreDetail = {
  checkpoint_id: number
  percent: number
  is_finish: boolean
  finish_time: number
  target_container: string
  checkpoint_name: string
  checkpoint_desc: string
}

export type missionScore = {
  mission_id: number
  mission_name: string
  mission_desc: string
  finish_score_counter: number
  all_score_counter: number
  score: number
  score_details: ScoreDetail[]
  total: number
}
