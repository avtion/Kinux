import { createStore, Store, useStore } from 'vuex'
import { InjectionKey } from 'vue'
import createPersistedState from 'vuex-persistedstate'
import { JWT } from './interfaces'
import { IsTimeOutLine } from '@/utils/time'

declare module '@vue/runtime-core' {
  // 定义类型
  interface State {
    JWT: JWT
  }

  interface ComponentCustomProperties {
    $store: Store<State>
  }
}

// 定义类型
export interface State {
  JWT: JWT
}

// Vuex Store
export const store = createStore({
  plugins: [createPersistedState({ key: 'kinux' })],
  state: {
    JWT: {
      Token: '',
      TTL: 0,
    },
  },
  mutations: {
    UpdateJWT(state, payload: JWT) {
      state.JWT = payload
    },
    ClearJWT(state) {
      state.JWT = <JWT>{
        Token: '',
        TTL: 0,
      }
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
  },
  actions: {},
})

// 定义注入的Key
export const key: InjectionKey<Store<State>> = Symbol('kinux-store')

// 定义获取Store方法
export function GetStore(): Store<State> {
  return useStore(key)
}
