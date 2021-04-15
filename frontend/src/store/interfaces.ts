// Json Web Token数据
export interface JWT {
  Token: string
  TTL: number
}

// 角色
enum Role {
  RoleAnonymous = 1, // 游客
  RoleNormalAccount,  // 普通用户
  RoleManager, // 管理员
  RoleAdmin, // 系统管理员
}

// 角色
export { Role }

export interface Profile {
  username: string
  realName: string
  role: string
  department: string
  avatarSeed: string
  dpID: string
  roleID: Role
}
