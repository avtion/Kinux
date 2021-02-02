import axios from "axios"
import { AxiosRequestConfig } from "axios"
import { notification } from "ant-design-vue"

// 默认Axios配置
const config: AxiosRequestConfig = {
  baseURL: "http://localhost:9001/",
}

// 路径
export const paths: object = {
  login: "v1/account/login",
}

// 默认客户端
export const defaultClient = axios.create(config)

// 基础响应
export class BaseResponse {
  public Code: number
  public Data: any

  // 构造函数
  constructor(obj: any) {
    if (
      obj == undefined ||
      obj["code"] == undefined ||
      obj["data"] == undefined
    ) {
      return
    }
    this.Code = obj["code"]
    this.Data = obj["data"]
  }

  // 后端规则约束
  IsSuccess(): boolean {
    if (2000 < this.Code || this.Code > 3000) {
      return false
    }
    return true
  }
}

// 响应拦截器
defaultClient.interceptors.response.use(
  (response) => {
    // 判断结果是否存在
    if (response.data == undefined) {
      return
    }
    const resp: BaseResponse = new BaseResponse(response.data)

    // 后端约束
    if (!resp.IsSuccess()) {
      // 前端UI框架提示
      notification.error({
        message: "请求失败",
        description: resp.Data,
      })

      // 控制台输出
      console.log("axios拦截服务端错误信息:", resp.Data)
      return Promise.reject(resp)
    }

    // 成功则返回正确响应
    return response
  },
  (error) => {
    // Do something with response error
    return Promise.reject(error)
  }
)
