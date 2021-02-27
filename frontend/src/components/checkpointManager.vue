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
          <a-input-search
            placeholder="输入需要查询的检查点"
            style="width: 200px"
            v-model:value="currentNameFilter"
            @search="onSearch"
          />
          <a-divider type="vertical" />
          <a-select
            v-model:value="currentCheckpointFilter"
            style="width: 200px"
            @change="onSearch"
          >
            <a-select-option :key="0">不筛选检测方式</a-select-option>
            <a-select-option
              v-for="item in CheckpointMethodMapper"
              :key="item[0]"
              >{{ item[1] }}</a-select-option
            >
          </a-select>
          <a-divider type="vertical" />
          <a-button @click="addFn">新增</a-button>
        </div>

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
          <!-- 徽章 -->
          <template #method="{ record }">
            <span>
              <a-tag
                :color="
                  record.method == CheckpointMethod.MethodExec ? 'green' : 'red'
                "
              >
                {{
                  record.method == CheckpointMethod.MethodExec
                    ? '输入流'
                    : '输出流'
                }}
              </a-tag>
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

    <!-- 追加和更新框 -->
    <a-modal
      :title="modalTitle"
      v-model:visible="modalVisible"
      :confirm-loading="isAddLoading || isEditLoading"
      @ok="commitModal"
      :afterClose="clearForm"
    >
      <a-form>
        <a-form-item label="检查点编号">
          <a-input
            v-model:value="formRef.ID"
            :disabled="true"
            style="width: 300px"
          />
        </a-form-item>
        <a-form-item label="检查点名称">
          <a-input v-model:value="formRef.Name" style="width: 300px" />
        </a-form-item>
        <a-form-item label="检查点描述">
          <a-input v-model:value="formRef.Desc" style="width: 300px" />
        </a-form-item>
        <a-form-item label="检测方式">
          <a-select v-model:value="formRef.Method" style="width: 300px">
            <a-select-option
              v-for="item in CheckpointMethodMapper"
              :key="item[0]"
              >{{ item[1] }}</a-select-option
            >
          </a-select>
        </a-form-item>
        <span>输入检测内容: </span>
        <v-ace-editor
          v-model:value="formRef.In"
          lang="powershell"
          theme="github"
          style="height: 200px"
        />
        <span>输出检测内容: </span>
        <v-ace-editor
          v-model:value="formRef.Out"
          lang="powershell"
          theme="github"
          style="height: 200px"
        />
      </a-form>
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
import { BaseResponse, defaultClient, paths } from '@/apis/request'

// vue3-ace-editor
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-powershell'
import 'ace-builds/src-noconflict/theme-github'

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
    title: '检测方法',
    dataIndex: 'method',
    key: 'method',
    slots: { customRender: 'method' },
  },
  {
    title: '描述',
    dataIndex: 'desc',
    key: 'desc',
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

// 检查点方法
enum CheckpointMethod {
  MethodExec = 1,
  MethodStdout,
  MethodTargetPort,
}

const CheckpointMethodMapper: Map<CheckpointMethod, string> = new Map([
  [CheckpointMethod.MethodExec, '输入流'],
  [CheckpointMethod.MethodStdout, '输出流'],
  [CheckpointMethod.MethodTargetPort, '目标端口'],
])

type ListParams = {
  page: number
  size: number
  name: string
  method: CheckpointMethod
}

type ListResult = {
  id: number
  method: CheckpointMethod
  name: string
  desc: string
  in: string
  out: string
  created_at: string
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  Name: string
  Desc: string
  In: string
  Out: string
  Method: CheckpointMethod
}

enum ModalStatus {
  add = 1,
  edit,
}

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.cp.list, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.cp.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(paths.cp.add, params)
}

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(paths.cp.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(paths.cp.delete + params + '/')
}

export default defineComponent({
  components: {
    VAceEditor,
    SmileOutlined,
    DownOutlined,
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    const currentCheckpointFilter = ref<CheckpointMethod>(0)

    // 表格
    const currentPage = ref<number>(1)
    const currentSize = ref<number>(10)
    const currentNameFilter = ref<string>('')
    const {
      data: listData,
      run: getListData,
      loading: isListDataLoading,
    } = useRequest(listDataAPI, {
      defaultParams: [
        <ListParams>{
          page: currentPage.value,
          size: currentSize.value,
          name: currentNameFilter.value,
          method: currentCheckpointFilter.value,
        },
      ],
      formatResult: (res): ListResults => {
        getTotal(<ListParams>{
          name: currentNameFilter.value,
        })
        return res.data.Data
      },
    })

    // 分页组件
    const { data: total, run: getTotal } = useRequest(countDataAPI, {
      defaultParams: [
        <ListParams>{
          name: currentNameFilter.value,
        },
      ],
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
      getListData(<ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        name: currentNameFilter.value,
        method: currentCheckpointFilter.value,
      })
    }

    // 更新表格
    const handleTableChange = (pag: Pagination) => {
      currentPage.value = pag.current
      currentSize.value = pag.pageSize
      getListData(<ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        name: currentNameFilter.value,
        method: currentCheckpointFilter.value,
      })
    }

    // 追加和更新窗口
    const modalVisible = ref<boolean>(false)
    const modalTitle = ref<string>('添加')
    const modalStatus = ref<ModalStatus>(ModalStatus.add)
    const formRef = reactive({
      ID: '',
      Name: '',
      Desc: '',
      In: '',
      Out: '',
      Method: CheckpointMethod.MethodExec,
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.ID = ''
      formRef.Name = ''
      formRef.Desc = ''
      formRef.In = ''
      formRef.Out = ''
      formRef.Method = CheckpointMethod.MethodExec
    }

    // 添加数据
    const { run: addDepartment, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [
        <UpdateParams>{
          ID: Number(formRef.ID),
          Name: formRef.Name,
          Desc: formRef.Desc,
          In: formRef.In,
          Out: formRef.Out,
          Method: formRef.Method,
        },
      ],
      manual: true,
    })

    // 修改数据
    const { run: editDepartment, loading: isEditLoading } = useRequest(
      editAPI,
      {
        defaultParams: [
          <UpdateParams>{
            ID: Number(formRef.ID),
            Name: formRef.Name,
            Desc: formRef.Desc,
            In: formRef.In,
            Out: formRef.Out,
            Method: formRef.Method,
          },
        ],
        manual: true,
      }
    )
    const { run: deleteDepartment } = useRequest(deleteAPI, {
      manual: true,
    })

    const addFn = () => {
      modalVisible.value = true
      modalTitle.value = '添加'
      modalStatus.value = ModalStatus.add
    }
    const editFn = (record: ListResult) => {
      modalVisible.value = true
      modalTitle.value = '修改'
      modalStatus.value = ModalStatus.edit
      formRef.ID = record.id + ''
      formRef.Name = record.name
      formRef.Desc = record.desc
      formRef.Method = record.method
      formRef.In = record.in
      formRef.Out = record.out
    }
    const deleteFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getListData(<ListParams>{
          page: currentPage.value,
          size: currentSize.value,
          name: currentNameFilter.value,
          method: currentCheckpointFilter.value,
        })
      })
    }

    // 提交表单
    const commitModal = () => {
      switch (modalStatus.value) {
        case ModalStatus.add:
          addDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            Name: formRef.Name,
            Desc: formRef.Desc,
            In: formRef.In,
            Out: formRef.Out,
            Method: formRef.Method,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name: currentNameFilter.value,
              method: currentCheckpointFilter.value,
            })
            modalVisible.value = false
          })
          break
        case ModalStatus.edit:
          editDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            Name: formRef.Name,
            Desc: formRef.Desc,
            In: formRef.In,
            Out: formRef.Out,
            Method: formRef.Method,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name: currentNameFilter.value,
              method: currentCheckpointFilter.value,
            })
            modalVisible.value = false
          })
          break
      }
    }

    // 当前管理界面的管理类型
    return {
      listData,
      columns,
      pagination,
      currentNameFilter,
      onSearch,
      isListDataLoading,
      handleTableChange,
      addFn,
      editFn,
      deleteFn,
      modalVisible,
      modalTitle,
      commitModal,
      isAddLoading,
      isEditLoading,
      formRef,
      clearForm,
      CheckpointMethod,
      CheckpointMethodMapper,
      currentCheckpointFilter,
    }
  },
})
</script>

<style lang="less" scoped>
</style>