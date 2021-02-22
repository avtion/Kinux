import { defaultClient, paths, BaseResponse } from './request'
import { store } from '../store/store'
import { JWT, Profile } from '../store/interfaces'

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
  login(): Promise<loginRespData> {
    return new Promise((resolve, reject) => {
      defaultClient
        .post(paths.ac.login, {
          username: this.username,
          password: this.password,
        })
        .then((res) => {
          const resp: loginRespData = new BaseResponse(res.data).Data

          // Json Web Token
          new Token(resp.token, resp.ttl).UpdateToken()

          // 更新用户的资料
          store.commit('UpdateProfile', <Profile>{
            username: resp.username,
            realName: resp.realName,
            role: resp.role,
            department: resp.department,
            avatarSeed: resp.avatarSeed,
          })

          resolve(resp)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }

  // 更新头像种子
  updateAvatarSeed(): Promise<string> {
    return new Promise((resolve, reject) => {
      defaultClient
        .put(paths.ac.updateAvatarSeed)
        .then((res) => {
          const resp: string = new BaseResponse(res.data).Data

          // 更新用户的头像种子
          store.commit('updateAvatarSeed', resp)

          resolve(resp)
        })
        .catch((err) => {
          reject(err)
        })
    })
  }

  // 修改密码
  updatePassword(oldValue: string, newValue: string): Promise<void> {
    return new Promise((resolve, reject) => {
      defaultClient
        .post(paths.ac.UpdatePassword, {
          old: oldValue,
          new: newValue,
        })
        .then((res) => {
          resolve()
        })
        .catch((err) => {
          reject(err)
        })
    })
  }
}

// JWT对象
export class Token {
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
    store.commit('UpdateJWT', jwt)
    // console.log("正在更新token, key:", this.key, "过期时间:", this.ttl)
  }

  ClearToken() {
    store.commit('ClearJWT')
  }
}

// 登陆响应结果
interface loginRespData {
  msg: string
  token: string
  ttl: string
  username: string
  realName: string
  role: string
  department: string
  avatarSeed: string
}
