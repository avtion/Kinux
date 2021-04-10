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
        <!-- 筛选栏 -->
        <div>
          <!-- 广播 -->
          <a-button
            @click="
              openSendMsgModal({
                id: 0,
                username: '',
                created_at: '',
                is_pty: false,
                lesson: 0,
                mission: 0,
                exam: 0,
                pty_meta_data: '',
              })
            "
            >广播</a-button
          >
          <a-divider type="vertical" />
          <!-- 刷新数据 -->
          <a-button @click="getListData">刷新</a-button>
        </div>

        <!-- 表格 -->
        <a-divider />
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
              <a-button
                type="primary"
                :disabled="!record.is_pty"
                @click="gotoWatcher(record)"
                >终端监控</a-button
              >
              <a-button type="primary" @click="openSendMsgModal(record)"
                >发送消息</a-button
              >
              <a-popconfirm
                placement="top"
                ok-text="是"
                cancel-text="否"
                @confirm="forceLogout({ target_id: record.id })"
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

    <!-- 发送消息使用的Modal -->
    <a-modal
      v-model:visible="sendMsgModalVisible"
      title="发送消息"
      :afterClose="sendMsgModalAfterClose"
      width="720px"
      @ok="sendMsg"
    >
      <a-form-item label="消息内容">
        <a-textarea
          auto-size
          v-model:value="msgText"
          placeholder="请输入要发送的消息"
          style="width: 420px"
        />
      </a-form-item>
    </a-modal>
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
    title: '实验信息',
    dataIndex: 'pty_meta_data',
    key: 'pty_meta_data',
  },
  {
    title: '登陆时间',
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

  lesson: number
  mission: number
  exam: number
}

type ListResults = ListResult[]

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>('/v1/ws/list/', {
    params: params,
  })
}

const sendMsgAPI = (params: { target_id: number; text: string }) => {
  return defaultClient.post<BaseResponse>('/v1/ws/msg/', params)
}

const forceLogoutAPI = (params: { target_id: number }) => {
  return defaultClient.post<BaseResponse>('/v1/ws/logout/', params)
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
        console.log(res)
        return res.data.Data
      },
    })

    //发送消息
    const { run: sendMsgFn } = useRequest(sendMsgAPI, {
      formatResult: (res): ListResults => {
        return res.data.Data
      },
      manual: true,
    })
    const msgTargetID = ref<number>()
    const msgText = ref<string>()
    const sendMsgModalVisible = ref<boolean>(false)

    const sendMsg = () => {
      sendMsgModalVisible.value = false
      sendMsgFn({ target_id: msgTargetID.value, text: msgText.value })
    }

    const sendMsgModalAfterClose = () => {
      msgTargetID.value = 0
      msgText.value = ''
    }
    const openSendMsgModal = (record: ListResult) => {
      sendMsgModalVisible.value = true
      msgTargetID.value = record.id
    }

    // 强制下线
    const { run: forceLogout } = useRequest(forceLogoutAPI, {
      formatResult: (res): ListResults => {
        return res.data.Data
      },
      manual: true,
    })

    // 前往监控
    const gotoWatcher = (item: ListResult) => {
      routers.push({
        name: 'shellWatcher',
        params: {
          mission: item.mission + '',
          exam: item.exam + '',
          lesson: item.lesson + '',
          target: Number(item.id),
        },
      })
    }

    // 当前管理界面的管理类型
    return {
      listData,
      columns,
      isListDataLoading,

      // 发送消息
      sendMsg,
      msgText,
      sendMsgModalVisible,
      sendMsgModalAfterClose,
      openSendMsgModal,
      forceLogout,

      // 获取当前链接信息
      getListData,

      // 跳转至监控
      gotoWatcher,
    }
  },
})
</script>

<style lang="less" scoped>
</style>