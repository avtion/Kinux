import { defaultClient, paths, BaseResponse } from "./request"
import { store } from "../store/store"
import { JWT } from "../store/interfaces"

export { Account }

// 账号定义
class Account {
  // 用户名
  public username: string
  // 密码
  public password: string

  // 构建函数
  constructor(username?: string, password?: string) {
    this.username = username
    this.password = password
  }

  // 登陆
  login(): Promise<any> {
    return new Promise((resolve, reject) => {
      defaultClient
        .post(paths["login"], {
          username: this.username,
          password: this.password,
        })
        .then((res) => {
          const resp = new BaseResponse(res.data)

          // Json Web Token
          const token: Token = new Token(resp.Data["token"], resp.Data["ttl"])
          token.UpdateToken()

          resolve(res)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }
}

// JWT对象
class Token {
  // 密钥
  key: string
  // 过期时间
  ttl: number

  constructor(token: string, ttl: string) {
    this.key = token
    // 时间转换时间
    this.ttl = Number(ttl)
  }

  // 更新JWT
  UpdateToken() {
    const jwt: JWT = {
      Token: this.key,
      TTL: this.ttl,
    }
    store.commit("UpdateJWT", jwt)
    // console.log("正在更新token, key:", this.key, "过期时间:", this.ttl)
  }
}
