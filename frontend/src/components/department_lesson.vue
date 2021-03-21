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
                placeholder="输入课程的名字搜索"
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
              <a-button @click="addFn">添加课程</a-button>
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
          <!-- 描述 -->
          <template #desc="{ record }">
            <div style="width: 300px">
              <a-tooltip :title="record.desc">
                <span class="desc">{{ record.desc }}</span>
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
                <a-button type="danger">移除课程</a-button>
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
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '描述',
    dataIndex: 'desc',
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
  department: number
}

type ListResult = {
  id: number
  name: string
  desc: string
}

type ListResults = ListResult[]

type UpdateParams = {
  department: number
  lesson: number
}

const apiPath = {
  list: '/v2/dl/list',
  count: '/v2/dl/count',
  add: '/v2/dl/',
  edit: '/v2/dl/',
  delete: '/v2/dl/',
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
    departmentID: {
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
        department: props.departmentID,
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

    const formRef = reactive({
      mission: undefined,
      priority: 0,
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.mission = undefined
      formRef.priority = 0
    }

    // 添加数据
    const getUpdateParams = (): UpdateParams => {
      return <UpdateParams>{
        department: Number(props.departmentID),
        lesson: Number(formRef.mission),
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
      name: string
      desc: string
    }
    type missionListParams = {
      page: number
      size: number
      name: string
    }
    const listMissionDataAPI = (params: missionListParams) => {
      return defaultClient.get<BaseResponse>('/v2/lesson/list', {
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