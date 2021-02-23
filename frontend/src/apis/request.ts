import axios, { AxiosResponse } from 'axios'
import { AxiosRequestConfig } from 'axios'
import { notification } from 'ant-design-vue'
import { store } from '@/store/store'
import routers from '@/routers/routers'

// 默认Axios配置
export const DefaultAxiosConfig: AxiosRequestConfig = {
  baseURL: 'http://localhost:9001/',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 15 * 1000,
  timeoutErrorMessage: '无法与服务器建立链接',
}

// 默认客户端
export const defaultClient = axios.create(DefaultAxiosConfig)

// 基础响应
export class BaseResponse {
  public Code: number
  public Data: any

  // 构造函数
  constructor(obj: any) {
    if (
      obj == undefined ||
      obj['code'] == undefined ||
      obj['data'] == undefined
    ) {
      return
    }
    this.Code = obj['code']
    this.Data = obj['data']
  }

  // 后端规则约束
  IsSuccess(): boolean {
    if (2000 < this.Code || this.Code > 3000) {
      return false
    }
    return true
  }

  // 判断是否JWT鉴权失败
  IsJWTAuthFailed(): boolean {
    return this.Code == respCode.JWTAuthFailed
  }
}

// 响应代码
enum respCode {
  Success = 2000, // 通用成功
  Failed = 4000, // 通用失败
  JWTAuthFailed = 4001, // JWT鉴权失败
}

// 响应拦截器
defaultClient.interceptors.response.use(
  (response: AxiosResponse<BaseResponse>) => {
    // 判断结果是否存在
    if (response.data == undefined) {
      return
    }
    const resp: BaseResponse = new BaseResponse(response.data)

    // 后端约束
    if (!resp.IsSuccess()) {
      // 前端UI框架提示
      notification.error({
        message: '请求失败',
        description: resp.Data,
      })

      if (resp.IsJWTAuthFailed()) {
        // JWT密钥失效则清空密钥缓存并跳转至登陆界面
        console.log('JWT鉴权失效', resp.Code, resp.Data)
        store.commit('ClearJWT')
        routers.push('/')
        return
      }

      // 控制台输出
      console.log('axios拦截服务端错误信息:', resp.Data)
      return Promise.reject(resp)
    }

    response.data = resp

    // 成功则返回正确响应
    return response
  },
  (err) => {
    // Do something with response error
    notification.error({
      message: '服务器发生未知错误',
    })
    console.log(err)
    return Promise.reject(err)
  }
)

// 请求拦截器 - 解决鉴权问题
defaultClient.interceptors.request.use(
  (cfg: AxiosRequestConfig): AxiosRequestConfig => {
    const token = store.getters.GetJWTToken
    if (<string>token != '') {
      cfg.headers.Authorization = 'Bearer ' + <string>token
    }
    return cfg
  }
)

// 路径
export const paths: routePath = {
  ac: {
    login: 'v1/account/login',
    updateAvatarSeed: 'v1/account/avatar',
    UpdatePassword: 'v1/account/pw',
  },
  ms: {
    list: 'v1/mission/',
    listContainersNames: 'v1/mission/cnames/',
    getGuide: 'v1/mission/guide/',
    dpOperation: 'v1/mission/op/',
  },
  department: {
    list: 'v1/department/',
    count: 'v1/department/count/',
    add: 'v1/department/',
    edit: 'v1/department/',
    delete: 'v1/department/',
  },
}

interface routePath {
  ac: account
  ms: mission
  department: department
}

interface account {
  login: string
  updateAvatarSeed: string
  UpdatePassword: string
}

interface mission {
  list: string
  listContainersNames: string
  getGuide: string
  dpOperation: string
}

interface department {
  list: string
  count: string
  add: string
  edit: string
  delete: string
}
