import { createApp } from 'vue'
// TypeScript error? Run VSCode command
// TypeScript: Select TypeScript version - > Use Workspace Version
import App from './App.vue'

const app = createApp(App)

// antd
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'
app.use(Antd)

// vuex
import { store, key } from './store/store'
app.use(store, key)

// vue-router
import router from '@/routers/routers'
app.use(router)

// Markdown
import VueMarkdownEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js'
import '@kangc/v-md-editor/lib/theme/style/vuepress.css'
import createEmojiPlugin from '@kangc/v-md-editor/lib/plugins/emoji/index'
import '@kangc/v-md-editor/lib/plugins/tip/tip.css'
// 支持多语言
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-c'
import 'prismjs/components/prism-cpp'
import 'prismjs/components/prism-markup-templating'
import 'prismjs/components/prism-php'
import 'prismjs/components/prism-java'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-typescript'

VueMarkdownEditor.use(vuepressTheme)
VueMarkdownEditor.use(createEmojiPlugin())
app.use(VueMarkdownEditor)

app.mount('#app')
