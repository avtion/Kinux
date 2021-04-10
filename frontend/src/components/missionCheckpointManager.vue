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
        <!-- 添加栏 -->
        <a-form layout="inline">
          <!-- 检查点 -->
          <a-form-item>
            <a-select
              style="width: 180px"
              placeholder="请选择检查点"
              show-search
              @search="searchCheckoptions"
              :show-arrow="false"
              :filter-option="false"
              :loading="isCheckpointsOpitonsLoading"
              v-model:value="formRef.check_point"
            >
              <a-select-option
                v-for="v in checkpointsOptions"
                :key="v.id"
                :title="v.name"
                ><span>
                  <a-tag
                    :color="
                      v.method == CheckpointMethod.MethodExec ? 'green' : 'cyan'
                    "
                  >
                    {{
                      v.method == CheckpointMethod.MethodExec
                        ? '输入流'
                        : '输出流'
                    }} </a-tag
                  >{{ v.name }}
                </span>
              </a-select-option></a-select
            >
          </a-form-item>

          <!-- 容器选项 -->
          <a-form-item>
            <a-select
              style="width: 140px"
              placeholder="请选择目标容器"
              v-model:value="formRef.target_container"
            >
              <a-select-option v-for="v in containersOptions" :key="v"
                >{{ v }}
              </a-select-option>
            </a-select>
          </a-form-item>

          <!-- 成绩比例 -->
          <a-form-item :label="'完成度占比（剩余' + restPercent + ')'">
            <a-input-number
              :min="0"
              :max="restPercent"
              :defaultValue="restPercent"
              v-model:value="formRef.percent"
            />
          </a-form-item>

          <!-- 权重 -->
          <a-form-item label="优先级">
            <a-input-number
              :defaultValue="0"
              v-model:value="formRef.priority"
            />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" html-type="submit" @click="addFn">
              新增
            </a-button>
          </a-form-item>
        </a-form>

        <!-- 表格 -->
        <a-divider />
        <a-table
          :columns="columns"
          :data-source="listData"
          :pagination="pagination"
          :loading="isListDataLoading"
          row-key="id"
          @change="handleTableChange"
        >
          <!-- 比例 -->
          <template #percent="{ record }">
            <div style="width: 60px">
              <a-input-number
                v-model:value="editableData[record.id].percent"
                :min="0"
                :max="
                  restPercent >= editableData[record.id].percent
                    ? restPercent
                    : editableData[record.id].percent
                "
                style="width: 60px"
                v-if="editableData[record.id]"
              />
              <a-progress
                :stroke-color="{
                  '0%': '#108ee9',
                  '100%': '#87d068',
                }"
                :percent="record.percent"
                v-else
              />
            </div>
          </template>

          <!-- 权重 -->
          <template #priority="{ record }">
            <div style="width: 60px">
              <a-input-number
                v-model:value="editableData[record.id].priority"
                v-if="editableData[record.id]"
                style="width: 60px"
              />
              <a-tag :color="record.priority <= 0 ? 'green' : 'red'" v-else>
                {{ record.priority }}
              </a-tag>
            </div>
          </template>

          <!-- 编辑框 -->
          <template #operation="{ record }">
            <!-- 修改状态 -->
            <a-button-group size="small" v-if="editableData[record.id]">
              <a-button type="primary" @click="saveFn(record)">保存</a-button>
              <a-button @click="cancelFn(record)">取消</a-button>
            </a-button-group>
            <!-- 非修改状态 -->
            <a-button-group size="small" v-else>
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
import { defineComponent, ref, UnwrapRef, reactive, watch, computed } from 'vue'

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

import { cloneDeep } from 'lodash-es'

// 分页组件定义
type Pagination = TableState['pagination']

// 表格的列参数
const columns = [
  {
    title: '检查点',
    dataIndex: 'checkpoint',
    key: 'checkpoint',
  },
  {
    title: '完成度比例',
    dataIndex: 'percent',
    slots: { customRender: 'percent' },
  },
  {
    title: '排序优先级',
    dataIndex: 'priority',
    slots: { customRender: 'priority' },
  },
  {
    title: '目标容器',
    dataIndex: 'target_container',
    key: 'target_container',
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
  mission: number
  containers: string[]
}

type ListResult = {
  id: number
  percent: number
  mission_id: number
  checkpoint_id: number
  priority: number
  mission: string
  checkpoint: string
  target_container: string
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  mission: number
  check_point: number
  percent: number
  priority: number
  target_container: string
}

enum ModalStatus {
  add = 1,
  edit,
}

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.mcp.list, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.mcp.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(paths.mcp.add, params)
}

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(paths.mcp.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(paths.mcp.delete + params + '/')
}

const percentAPI = (params: number) => {
  return defaultClient.get<BaseResponse>(paths.mcp.percent + params + '/')
}

const checkpointsOptionsAPI = (params?: string) => {
  return defaultClient.get<BaseResponse>(paths.cp.quick, {
    params: { name: params },
  })
}

const containersOptionsAPI = (params: number) => {
  return defaultClient.get<BaseResponse>(
    paths.ms.listContainersNames + params + '/'
  )
}

type checkpointsOptionsResult = {
  id: number
  method: number
  name: string
  tag: string
}[]

// 检查点方法
enum CheckpointMethod {
  MethodExec = 1,
  MethodStdout,
  MethodTargetPort,
}

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
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 表格
    const currentPage = ref<number>(1)
    const currentSize = ref<number>(5)
    const currentContainersFilter = ref<string[]>([])
    const getListParams = (): ListParams => {
      return <ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        containers: currentContainersFilter.value,
        mission: props.missionID,
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
      pageSizeOptions: ['5', '10', '15'],
    })

    // 获取当前已经占用的成绩比例
    const { data: usedPercent, run: getUsedPercent } = useRequest(percentAPI, {
      defaultParams: [props.missionID],
      formatResult: (res): number => {
        return <number>res.data.Data
      },
    })
    watch(
      () => total,
      () => {
        getUsedPercent(props.missionID)
      }
    )
    const restPercent = computed((): number => {
      return 100 - usedPercent.value
    })

    // 检查点列表
    const {
      data: checkpointsOptions,
      loading: isCheckpointsOpitonsLoading,
      run: getCheckpointOptions,
    } = useRequest(checkpointsOptionsAPI, {
      formatResult: (res): checkpointsOptionsResult => {
        return <checkpointsOptionsResult>res.data.Data
      },
    })
    const searchCheckoptions = (value: string) => {
      getCheckpointOptions(value)
    }

    // 实验容器选项列表
    const { data: containersOptions } = useRequest(containersOptionsAPI, {
      defaultParams: [props.missionID],
      formatResult: (res): string[] => {
        return <string[]>res.data.Data
      },
    })

    // 更新表格
    const handleTableChange = (pag: Pagination) => {
      currentPage.value = pag.current
      currentSize.value = pag.pageSize
      getListData(getListParams())
    }

    // 追加和更新窗口
    const formRef = reactive(<UpdateParams>{
      mission: props.missionID,
      percent: 0,
      priority: 0,
    })
    const getUpdateParams = (): UpdateParams => {
      return formRef
    }
    // 添加数据
    const { run: addDepartment, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [getUpdateParams()],
      manual: true,
    })
    const addFn = () => {
      addDepartment(getUpdateParams()).finally(() => {
        getUsedPercent(props.missionID)
        getListData(getListParams())
      })
    }

    // 修改数据
    const editableData: UnwrapRef<Record<number, ListResult>> = reactive({})
    const { run: editDepartment, loading: isEditLoading } = useRequest(
      editAPI,
      {
        defaultParams: [getUpdateParams()],
        manual: true,
      }
    )
    const { run: deleteDepartment } = useRequest(deleteAPI, {
      manual: true,
    })

    const editFn = (record: ListResult) => {
      editableData[record.id] = cloneDeep(record)
    }

    const saveFn = (record: ListResult) => {
      const data = editableData[record.id]
      editDepartment(<UpdateParams>{
        ID: record.id,
        mission: props.missionID,
        // checkpointid不修改
        check_point: record.checkpoint_id,
        percent: data.percent,
        priority: data.priority,
        target_container: data.target_container,
      }).finally(() => {
        getUsedPercent(props.missionID)
        getListData(getListParams())
        delete editableData[record.id]
      })
    }

    const cancelFn = (record: ListResult) => {
      delete editableData[record.id]
    }

    // 删除数据
    const deleteFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getUsedPercent(props.missionID)
        getListData(getListParams())
      })
    }

    // 当前管理界面的管理类型
    return {
      listData,
      columns,
      pagination,
      currentContainersFilter,
      isListDataLoading,
      handleTableChange,
      addFn,
      editFn,
      deleteFn,
      isAddLoading,
      isEditLoading,
      formRef,
      restPercent,
      checkpointsOptions,
      CheckpointMethod,
      isCheckpointsOpitonsLoading,
      searchCheckoptions,
      containersOptions,
      editableData,
      saveFn,
      cancelFn,
    }
  },
})
</script>

<style lang="less" scoped>
</style>