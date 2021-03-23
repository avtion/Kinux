<template>
  <div></div>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, ref, UnwrapRef, reactive } from 'vue'

// antd
import { TableState } from 'ant-design-vue/es/table/interface'
import { SmileOutlined, DownOutlined } from '@ant-design/icons-vue'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

// vue-request
import { useRequest } from 'vue-request'
import { BaseResponse, defaultClient, paths } from '@/apis/request'

// 分页组件定义
type Pagination = TableState['pagination']

export default defineComponent({
  components: {
    SmileOutlined,
    DownOutlined,
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 当前管理界面的管理类型
    return {}
  },
})
</script>

<style lang="less" scoped>
</style>