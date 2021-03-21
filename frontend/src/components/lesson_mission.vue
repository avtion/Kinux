<template>
  <div>
    <a-layout>
      <a-layout-content
        :style="{
          background: '#fff',
          minHeight: '600px',
          padding: '12px',
        }"
      >
        <!-- 筛选栏 -->
        <div style="margin-bottom: 8px">
          <a-form layout="inline">
            <!-- 实验 -->
            <a-form-item>
              <a-select
                style="width: 180px"
                placeholder="输入实验名称搜索"
                show-search
                @search="searchMissionListOptions"
                :show-arrow="false"
                :filter-option="false"
                :loading="isMissionOpitonsLoading"
                v-model:value="formRef.mission"
              >
                <a-select-option
                  v-for="v in missionListData"
                  :key="v.id"
                  :title="v.name"
                  >{{ v.id }} - {{ v.name }}
                </a-select-option></a-select
              >
            </a-form-item>
            <a-form-item>
              <a-input-number
                placeholder="优先级"
                v-model.value="formRef.priority"
              />
            </a-form-item>
            <a-form-item>
              <a-button @click="addFn">添加实验</a-button>
            </a-form-item>
          </a-form>
        </div>

        <!-- 表格 -->
        <a-table
          class="ant-table-striped"
          :columns="columns"
          :data-source="listData"
          :pagination="pagination"
          :loading="isListDataLoading"
          row-key="id"
          @change="handleTableChange"
          :rowClassName="
            (record, index) => (index % 2 === 1 ? 'table-striped' : null)
          "
          bordered
          tableLayout="auto"
        >
          <!-- 编辑框 -->
          <template #desc="{ record }">
            <div style="width: 300px">
              <a-tooltip :title="record.mission_desc">
                <span class="desc">{{ record.mission_desc }}</span>
              </a-tooltip>
            </div>
          </template>

          <!-- 编辑框 -->
          <template #operation="{ record }">
            <a-button-group size="small">
              <a-popconfirm
                placement="top"
                ok-text="是"
                cancel-text="否"
                @confirm="deleteFn(record)"
              >
                <template #title>
                  <p>确定删除？</p>
                </template>
                <a-button type="danger">移除实验</a-button>
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
import { BaseResponse, defaultClient, paths } from '@/apis/request'

// 分页组件定义
type Pagination = TableState['pagination']

// 表格的列参数
const columns = [
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
  },
  {
    title: '名称',
    dataIndex: 'mission_name',
    key: 'mission_name',
  },
  {
    title: '描述',
    dataIndex: 'mission_desc',
    slots: { customRender: 'desc' },
  },
  {
    title: '操作',
    dataIndex: 'id',
    slots: { customRender: 'operation' },
  },
]

type ListParams = {
  page: number
  size: number
  lesson: number
}

type ListResult = {
  id: number
  mission_id: number
  mission_name: string
  mission_desc: string
  priority: number
}

type ListResults = ListResult[]

type UpdateParams = {
  lesson: number
  mission: number
  priority: number
}

enum ModalStatus {
  add = 1,
  edit,
}

const apiPath = {
  list: '/v2/lm/list',
  count: '/v2/lm/count',
  add: '/v2/lm/',
  edit: '/v2/lm/',
  delete: '/v2/lm/',
}

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(apiPath.list, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(apiPath.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(apiPath.add, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(apiPath.delete + params + '/')
}

export default defineComponent({
  components: {
    SmileOutlined,
    DownOutlined,
  },
  props: {
    lessonID: {
      type: Number,
      required: true,
    },
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 表格
    const currentPage = ref<number>(1)
    const currentSize = ref<number>(10)
    const getListParams = (): ListParams => {
      return <ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        lesson: props.lessonID,
      }
    }

    const {
      data: listData,
      run: getListData,
      loading: isListDataLoading,
    } = useRequest(listDataAPI, {
      defaultParams: [getListParams()],
      formatResult: (res): ListResults => {
        getTotal(getListParams())
        return res.data.Data
      },
    })

    // 分页组件
    const { data: total, run: getTotal } = useRequest(countDataAPI, {
      defaultParams: [getListParams()],
      formatResult: (res): number => {
        return res.data.Data
      },
    })
    const pagination: UnwrapRef<Pagination> = reactive({
      current: currentPage,
      pageSize: currentSize,
      total: total,
      showSizeChanger: true,
      pageSizeOptions: ['10', '20', '30', '40'],
    })

    // 搜索按钮
    const onSearch = (value: string) => {
      getListData(getListParams())
    }

    // 更新表格
    const handleTableChange = (pag: Pagination) => {
      currentPage.value = pag.current
      currentSize.value = pag.pageSize
      getListData(getListParams())
    }

    // 追加和更新窗口
    const modalVisible = ref<boolean>(false)
    const modalTitle = ref<string>('添加')
    const modalStatus = ref<ModalStatus>(ModalStatus.add)
    const formRef = reactive({
      mission: '',
      priority: 0,
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.mission = ''
      formRef.priority = 0
    }

    // 添加数据
    const getUpdateParams = (): UpdateParams => {
      return <UpdateParams>{
        lesson: Number(props.lessonID),
        mission: Number(formRef.mission),
        priority: Number(formRef.priority),
      }
    }

    // 添加数据
    const { run: addMission, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [getUpdateParams()],
      manual: true,
    })
    const addFn = () => {
      addMission(getUpdateParams()).finally(() => {
        getListData(getListParams())
      })
    }

    // 实验列表
    const missionNameFilter = ref<string>('')
    type missionListResult = {
      id: number
      total: number
      name: string
      desc: string
      containers: string[]
      deployment: number
      exec_container: string
      command: string
    }
    type missionListParams = {
      page: number
      size: number
      name: string
    }
    const listMissionDataAPI = (params: missionListParams) => {
      return defaultClient.get<BaseResponse>(paths.ms.manageList, {
        params: params,
      })
    }
    const {
      data: missionListData,
      loading: isMissionOpitonsLoading,
      run: getMissionListOptions,
    } = useRequest(listMissionDataAPI, {
      defaultParams: [
        <missionListParams>{
          page: 0,
          size: 0,
          name: missionNameFilter.value,
        },
      ],
      formatResult: (res): missionListResult[] => {
        return res.data.Data
      },
    })
    const searchMissionListOptions = (value: string) => {
      getMissionListOptions(<missionListParams>{
        page: 0,
        size: 0,
        name: value,
      })
    }

    // 删除数据
    const { run: deleteDepartment } = useRequest(deleteAPI, {
      manual: true,
    })
    const deleteFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getListData(getListParams())
      })
    }

    // 当前管理界面的管理类型
    return {
      listData,
      columns,
      pagination,
      onSearch,
      isListDataLoading,
      handleTableChange,
      addFn,
      modalVisible,
      modalTitle,
      isAddLoading,
      formRef,
      clearForm,
      isMissionOpitonsLoading,
      missionListData,
      searchMissionListOptions,
      deleteFn,
    }
  },
})
</script>

<style lang="less" scoped>
.ant-table-striped :deep(.table-striped) {
  background-color: #fafafa;
}
.desc {
  overflow: hidden;
  -webkit-line-clamp: 1;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
}
</style>