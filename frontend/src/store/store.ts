import { createStore, Store, useStore } from 'vuex'
import { InjectionKey } from 'vue'
import createPersistedState from 'vuex-persistedstate'
import { JWT, Profile } from './interfaces'
import { IsTimeOutLine } from '@/utils/time'

declare module '@vue/runtime-core' {
  // 定义类型
  interface State {
    JWT: JWT
    Profile: Profile
  }

  interface ComponentCustomProperties {
    $store: Store<State>
  }
}

// 定义类型
export interface State {
  JWT: JWT
  Profile: Profile
}

// Vuex Store
export const store = createStore({
  plugins: [createPersistedState({ key: 'kinux' })],
  state: {
    JWT: <JWT>{},
    Profile: <Profile>{},
  },
  mutations: {
    UpdateJWT(state, payload: JWT) {
      state.JWT = payload
    },
    ClearJWT(state) {
      state.JWT = <JWT>{}
    },
    UpdateProfile(state, payload: Profile) {
      state.Profile = payload
    },
    ClearProfile(state) {
      state.Profile = <Profile>{}
    },
    updateAvatarSeed(state, payload: string) {
      state.Profile.avatarSeed = payload
    },
  },
  getters: {
    // 获取JWT密钥
    GetJWTToken(state) {
      // 判断密钥是否过期
      if (IsTimeOutLine(state.JWT.TTL)) {
        return ''
      }
      return state.JWT.Token
    },

    // 获取用户资料
    GetProfile(state) {
      return state.Profile
    },

    // 头像种子
    GetAvatarSeed(state) {
      return state.Profile.avatarSeed == ''
        ? state.Profile.realName == ''
          ? state.Profile.username
          : state.Profile.realName
        : state.Profile.avatarSeed
    },
  },
  actions: {},
})

// 定义注入的Key
export const key: InjectionKey<Store<State>> = Symbol('kinux-store')

// 定义获取Store方法
export function GetStore(): Store<State> {
  return useStore(key)
}
