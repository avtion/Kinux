import { createStore, Store, useStore } from "vuex"
import { ComponentCustomProperties, InjectionKey } from "vue"
import createPersistedState from "vuex-persistedstate"
import { JWT } from "./interfaces"

declare module "@vue/runtime-core" {
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
  plugins: [createPersistedState({ key: "kinux" })],
  state: {
    JWT: {
      Token: "",
      TTL: 0,
    },
  },
  mutations: {
    UpdateJWT(state, payload: JWT) {
      state.JWT = payload
    },
  },
  getters: {},
  actions: {},
})

// 定义注入的Key
export const key: InjectionKey<Store<State>> = Symbol("kinux-store")

// 定义获取Store方法
export function GetStore(): Store<State> {
  return useStore(key)
}
