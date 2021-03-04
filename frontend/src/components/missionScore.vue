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
        <a-row>
          <a-select
            style="width: 140px"
            v-model:value="currentDepartmentFilter"
            @change="onSearch"
            placeholder="请选择班级"
          >
            <a-select-option v-for="v in departmentOptions" :key="v.id">{{
              v.name
            }}</a-select-option>
          </a-select>
        </a-row>
        <a-divider />
        <template v-if="listData != undefined">
          <a-row>
            <a-statistic title="满分" :value="listData.total" />
            <a-divider type="vertical" />
            <a-statistic title="总人数" :value="totalAccount" />
            <a-divider type="vertical" />
            <a-statistic title="完成人数" :value="listData.finish_count" />
            <a-divider type="vertical" />
            <a-statistic
              title="未完成人数"
              :value="totalAccount - listData.finish_count"
            />
          </a-row>
        </template>
        <template v-else>
          <a-skeleton active :paragraph="{ rows: 2 }" />
        </template>

        <!-- 表格 -->
        <a-divider />
        <a-table
          :columns="columns"
          :data-source="resultData"
          :pagination="pagination"
          :loading="isListDataLoading"
          row-key="id"
          @change="handleTableChange"
        >
          <!-- 检查点统计 -->
          <template #checkpointCounter="{ record }">
            <span>
              <a-progress
                :percent="(record.finish_checkpoints / listData.cps_num) * 100"
                size="small"
              />
            </span>
          </template>

          <!-- 编辑框 -->
          <template #operation="{ record }">
            <a-button-group size="small">
              <a-button type="primary" @click="editFn(record)">修改</a-button>
              <a-popconfirm
                placement="top"
                ok-text="是"
                cancel-text="否"
                @confirm="deleteFn(record)"
              >
                <template #title>
                  <p>确定删除？</p>
                </template>
                <a-button type="danger">删除</a-button>
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
import {
  defineComponent,
  ref,
  onMounted,
  UnwrapRef,
  reactive,
  computed,
  watch,
} from 'vue'

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
    title: '序号',
    dataIndex: 'id',
    key: 'id',
  },

  {
    title: '用户',
    dataIndex: 'username',
    key: 'username',
  },
  {
    title: '班级',
    dataIndex: 'department',
    key: 'department',
  },
  {
    title: '成绩',
    dataIndex: 'score',
    key: 'score',
  },
  {
    title: '检查点统计',
    dataIndex: 'finish_checkpoints',
    slots: { customRender: 'checkpointCounter' },
  },
]

type ListParams = {
  department: number
  mission: number
}

type ListResult = {
  total: number
  cps_num: number
  finish_count: number
  data: ListResults
}

type ListResults = {
  id: number
  account: number
  department_id: number
  username: string
  department: string
  score: number
  finish_checkpoints: number
}[]

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.score.msList, {
    params: params,
  })
}

const departmentOptionsAPI = (params: departmentOptionsParams) => {
  return defaultClient.get<BaseResponse>(paths.department.quick, {
    params: params,
  })
}

type departmentOptionsParams = {
  ns: string
}

type departmentOptionsResult = {
  id: number
  name: string
}[]

export default defineComponent({
  components: {
    SmileOutlined,
    DownOutlined,
  },
  props: {
    missionID: {
      type: Number,
      required: true,
    },
    namespace: {
      type: String,
      required: true,
    },
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 班级选项
    const currentDepartmentFilter = ref<number>()
    const { data: departmentOptions } = useRequest(departmentOptionsAPI, {
      formatResult: (res): departmentOptionsResult => {
        if (
          currentDepartmentFilter.value == undefined &&
          (<departmentOptionsResult>res.data.Data).length > 0
        ) {
          currentDepartmentFilter.value = (<departmentOptionsResult>(
            res.data.Data
          ))[0].id
        }
        return <departmentOptionsResult>res.data.Data
      },
      defaultParams: [
        <departmentOptionsParams>{
          ns: props.namespace,
        },
      ],
    })
    watch(
      () => currentDepartmentFilter.value,
      () => {
        onSearch()
      }
    )

    const getListParams = (): ListParams => {
      return <ListParams>{
        mission: props.missionID,
        department: currentDepartmentFilter.value,
      }
    }

    // 表格
    const currentPage = ref<number>(1)
    const currentSize = ref<number>(5)
    const {
      data: listData,
      run: getListData,
      loading: isListDataLoading,
    } = useRequest(listDataAPI, {
      defaultParams: [getListParams()],
      formatResult: (res): ListResult => {
        return res.data.Data
      },
      manual: true,
    })
    const total = computed(() =>
      listData.value == undefined ? 0 : listData.value.data.length
    )
    const resultData = computed(
      (): ListResults =>
        listData.value == undefined ? <ListResults>[] : listData.value.data
    )

    // 分页组件
    const pagination: UnwrapRef<Pagination> = reactive({
      current: currentPage,
      pageSize: currentSize,
      total: total,
      showSizeChanger: true,
      pageSizeOptions: ['5', '10', '15', '20'],
    })

    // 搜索按钮
    const onSearch = (value?: string) => {
      getListData(getListParams())
    }

    // 更新表格
    const handleTableChange = (pag: Pagination) => {
      currentPage.value = pag.current
      currentSize.value = pag.pageSize
      getListData(getListParams())
    }

    // 当前管理界面的管理类型
    return {
      listData,
      columns,
      pagination,
      currentDepartmentFilter,
      onSearch,
      isListDataLoading,
      handleTableChange,
      departmentOptions,
      resultData,
      totalAccount: total,
    }
  },
})
</script>

<style lang="less" scoped>
</style>