import { AxiosResponse } from 'axios'
import { defaultClient, paths, BaseResponse } from './request'

export class mission {
  ms: missionStatus

  // 获取所有的任务
  list(): Promise<missionList[]> {
    return new Promise<missionList[]>(
      (
        resolve: (value: missionList[]) => void,
        reject: (reason: string) => void
      ) => {
        defaultClient
          .get(paths.ms.list)
          .then((res: AxiosResponse<BaseResponse>) => {
            const ml: missionList[] = res.data.Data
            if (ml.length == 0) {
              return reject('无可用数据')
            }
            return resolve(ml)
          })
        return
      }
    )
  }

  // 获取任务的容器名
  listContainersNames(missionID: string): Promise<String[]> {
    return new Promise<string[]>(
      (
        resolve: (value: string[]) => void,
        reject: (reason: string) => void
      ) => {
        defaultClient
          .get(paths.ms.listContainersNames + missionID + '/')
          .then((res: AxiosResponse<BaseResponse>) => {
            const ml: string[] = res.data.Data
            if (ml.length == 0) {
              return reject('无可用数据')
            }
            return resolve(ml)
          })
        return
      }
    )
  }

  // 获取任务的实验文档
  getGuide(missionID: string): Promise<string> {
    return new Promise<string>(
      (resolve: (value: string) => void, reject: (reason: string) => void) => {
        defaultClient
          .get(paths.ms.getGuide + missionID + '/')
          .then((res: AxiosResponse<BaseResponse>) => {
            const guide: string = res.data.Data
            return resolve(guide)
          })
        return
      }
    )
  }

  // 删除正在进行的deployment
  deleteDeployment(missionID: string): Promise<string> {
    return new Promise<string>(
      (resolve: (value: string) => void, reject: (reason: string) => void) => {
        defaultClient
          .delete(paths.ms.dpOperation + missionID + '/')
          .then((res: AxiosResponse<BaseResponse>) => {
            if (!res.data.IsSuccess()) {
              return reject(res.data.Data)
            }
            return resolve(res.data.Data)
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
  Block, // 阻塞
}

// 任务列表
export interface missionList {
  id: number
  name: string
  desc: string
  guide: string
  status: missionStatus
}
