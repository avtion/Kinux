import { AxiosResponse } from 'axios'
import { defaultClient, paths, BaseResponse } from './request'

export class mission {
  ms: missionStatus

  list(): Promise<missionList[]> {
    return new Promise<missionList[]>(
      (
        resolve: (value: missionList[]) => void,
        reject: (reason: string) => void
      ) => {
        defaultClient.get(paths.ms.list).then((res: AxiosResponse) => {
          const ml: missionList[] = new BaseResponse(res.data).Data
          if (ml.length == 0) {
            return reject('无可用数据')
          }
          return resolve(ml)
        })
        return
      }
    )
  }
}

// 任务状态
export enum missionStatus {
  Stop = 1, // 未启动
  Pending, // 正在启动
  Working, // 正在运行
  Done, // 已经完成
}

// 任务列表
export interface missionList {
  ID: number
  Name: string
  Desc: string
  Guide: string
  Status: missionStatus
}
