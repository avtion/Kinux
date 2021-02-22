<template>
  <div>管理界面 {{ managerType }}</div>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, ref, watch } from 'vue'
// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

export default defineComponent({
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 当前管理界面的管理类型
    const managerType = ref<string>(
      routers.currentRoute.value.meta['managerType']
    )
    watch(
      () => routers.currentRoute.value.meta['managerType'],
      (newValue: string) => {
        managerType.value = newValue
      }
    )

    return {
      managerType,
    }
  },
})
</script>

<style lang="less" scoped>
</style>