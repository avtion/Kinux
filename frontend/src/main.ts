import { createApp } from "vue"
// TypeScript error? Run VSCode command
// TypeScript: Select TypeScript version - > Use Workspace Version
import App from "./App.vue"

const app = createApp(App)

// antd
import Antd from "ant-design-vue"
import "ant-design-vue/dist/antd.css"
app.use(Antd)

// vuex
import {store, key} from "./store/store"
app.use(store, key)

// vue-router
import router from '@/routers/routers'
app.use(router)

app.mount("#app")
