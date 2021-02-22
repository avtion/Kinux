// Json Web Token数据
export interface JWT {
  Token: string
  TTL: number
}

export interface Profile {
  username: string
  realName: string
  role: string
  department: string
}
