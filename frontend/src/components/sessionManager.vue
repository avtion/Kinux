<template>
  <div>
    <a-layout style="padding: 24px 24px 24px">
      <a-layout-content
        :style="{
          background: '#fff',
          minHeight: '600px',
          padding: '12px',
        }"
      >
        <!-- 表格 -->
        <a-table
          :columns="columns"
          :data-source="listData"
          :pagination="false"
          :loading="isListDataLoading"
          row-key="id"
        >
          <!-- 是否使用终端 -->
          <template #pty="{ record }">
            <span>
              <a-tag :color="record.is_pty ? 'green' : 'red'">
                {{ record.is_pty ? '是' : '否' }}
              </a-tag>
            </span>
          </template>
          <!-- 编辑框 -->
          <template #operation="{ record }">
            <a-button-group size="small">
              <a-button type="primary" :disabled="!record.is_pty"
                >终端监控</a-button
              >
              <a-button type="primary">发送消息</a-button>
              <a-popconfirm
                placement="top"
                ok-text="是"
                cancel-text="否"
                @confirm="deleteFn(record)"
              >
                <template #title>
                  <p>确定强制用户下线？</p>
                </template>
                <a-button type="danger">强制下线</a-button>
              </a-popconfirm>
            </a-button-group>
          </template>
        </a-table>
      </a-layout-content>
    </a-layout>
  </div>
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
import { BaseResponse, defaultClient } from '@/apis/request'

// 分页组件定义
type Pagination = TableState['pagination']

// 表格的列参数
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },

  {
    title: '用户',
    dataIndex: 'username',
    key: 'username',
  },
  {
    title: '终端运行中',
    dataIndex: 'is_pty',
    key: 'is_pty',
    slots: { customRender: 'pty' },
  },
  {
    title: '终端数据',
    dataIndex: 'pty_meta_data',
    key: 'pty_meta_data',
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
  },
  {
    title: '操作',
    dataIndex: 'id',
    slots: { customRender: 'operation' },
  },
]

type ListParams = {}

type ListResult = {
  id: number
  username: string
  created_at: string
  is_pty: boolean
  pty_meta_data: string
}

type ListResults = ListResult[]

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>('/v1/ws/list/', {
    params: params,
  })
}

export default defineComponent({
  components: {
    SmileOutlined,
    DownOutlined,
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 表格
    const {
      data: listData,
      run: getListData,
      loading: isListDataLoading,
    } = useRequest(listDataAPI, {
      defaultParams: [<ListParams>{}],
      formatResult: (res): ListResults => {
        return res.data.Data
      },
    })

    // 当前管理界面的管理类型
    return {
      listData,
      columns,
      isListDataLoading,
    }
  },
})
</script>

<style lang="less" scoped>
</style>