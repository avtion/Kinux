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
          console.log(res)
          resolve(<missionScore>res.data.Data)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }

  // 获取考试成绩
  GetExamScore = () => {
    return new Promise((resolve: (res: ExamScore) => void, reject) => {
      defaultClient
        .get('/v2/score/exam/', {
          params: {
            dp: this.department,
            lesson: this.lesson,
            exam: this.exam,
          },
        })
        .then((res: AxiosResponse<BaseResponse>) => {
          resolve(<ExamScore>res.data.Data)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }

  // 管理员获取考试成绩
  GetExamScoreForAdmin = () => {
    return new Promise(
      (resolve: (res: ExamScoreForAdmin[]) => void, reject) => {
        defaultClient
          .get('/v2/score/exam/', {
            params: {
              dp: this.department,
              lesson: this.lesson,
              exam: this.exam,
            },
          })
          .then((res: AxiosResponse<BaseResponse>) => {
            resolve(<ExamScoreForAdmin[]>res.data.Data)
          })
          .catch((err) => {
            reject(err)
          })
      }
    )
  }

  // 管理员获取实验成绩
  GetMissionScoreForAdmin = () => {
    return new Promise(
      (resolve: (res: MissionScoreForAdmin[]) => void, reject) => {
        defaultClient
          .get('/v2/score/mission/', {
            params: {
              dp: this.department,
              lesson: this.lesson,
              mission: this.mission,
            },
          })
          .then((res: AxiosResponse<BaseResponse>) => {
            resolve(<MissionScoreForAdmin[]>res.data.Data)
          })
          .catch((err) => {
            reject(err)
          })
      }
    )
  }
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

export interface ExamScore {
  exam_id: number
  exam_name: string
  exam_desc: string
  exam_begin_at: number
  exam_end_at: number
  score: number
  mission_scores: missionScore[]
  total: number
}

export interface ExamScoreForAdmin {
  exam_id: number
  exam_name: string
  exam_desc: string
  exam_begin_at: number
  exam_end_at: number
  score: number
  mission_scores: missionScore[]
  total: number
  pos: number
  id: number
  role: number
  profile: number
  username: string
  real_name: string
  department: string
  department_id: number
}

export interface MissionScoreForAdmin {
  mission_id: number
  mission_name: string
  mission_desc: string
  finish_score_counter: number
  all_score_counter: number
  score: number
  score_details: ScoreDetail[]
  total: number
  pos: number
  id: number
  role: number
  profile: number
  username: string
  real_name: string
  department: string
  department_id: number
}
