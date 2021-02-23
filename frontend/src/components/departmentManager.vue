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
            placeholder="输入需要查询的班级"
            style="width: 200px"
            v-model:value="currentNameFilter"
            @search="onSearch"
          />
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
          <template #namespace="{ text: namespace }">
            <span>
              <a-tag
                v-for="(name, index) in namespace"
                :key="index"
                color="green"
              >
                {{ name }}
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
        <a-form-item label="班级编号">
          <a-input
            v-model:value="formRef.ID"
            :disabled="true"
            style="width: 400px"
          />
        </a-form-item>
        <a-form-item label="班级名称">
          <a-input v-model:value="formRef.Name" style="width: 400px" />
        </a-form-item>
        <a-form-item label="命名空间">
          <a-select
            v-model:value="formRef.Namespace"
            mode="tags"
            style="width: 400px"
            placeholder="输入命名空间并摁下回车添加"
          >
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, ref, watch, UnwrapRef, reactive } from 'vue'

// antd
import {
  TableState,
  TableStateFilters,
} from 'ant-design-vue/es/table/interface'
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
    title: '班级名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '可视命名空间',
    dataIndex: 'namespace',
    slots: { customRender: 'namespace' },
    key: 'namespace',
  },
  {
    title: '创建时间',
    dataIndex: 'creat_at',
    key: 'creat_at',
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
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
  name_filter: string
}

type ListResult = {
  id: number
  name: string
  creat_at: string
  updated_at: string
  namespace: string[]
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  Name: string
  Namespace: string[]
}

enum ModalStatus {
  add = 1,
  edit,
}

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.department.list, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.department.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(paths.department.add, params)
}

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(paths.department.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(
    paths.department.delete + params + '/'
  )
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
          name_filter: currentNameFilter.value,
        },
      ],
      formatResult: (res): ListResults => {
        getTotal(<ListParams>{
          name_filter: currentNameFilter.value,
        })
        return res.data.Data
      },
    })

    // 分页组件
    const { data: total, run: getTotal } = useRequest(countDataAPI, {
      defaultParams: [
        <ListParams>{
          name_filter: currentNameFilter.value,
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
        name_filter: currentNameFilter.value,
      })
    }

    // 更新表格
    const handleTableChange = (pag: Pagination) => {
      currentPage.value = pag.current
      currentSize.value = pag.pageSize
      getListData(<ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        name_filter: currentNameFilter.value,
      })
    }

    // 追加和更新窗口
    const modalVisible = ref<boolean>(false)
    const modalTitle = ref<string>('添加')
    const modalStatus = ref<ModalStatus>(ModalStatus.add)
    const formRef = reactive({
      ID: '',
      Name: '',
      Namespace: [],
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.ID = ''
      formRef.Name = ''
      formRef.Namespace = []
    }

    // 添加数据
    const { run: addDepartment, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [
        <UpdateParams>{
          ID: Number(formRef.ID),
          Name: formRef.Name,
          Namespace: formRef.Namespace,
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
            Namespace: formRef.Namespace,
          },
        ],
        manual: true,
      }
    )
    const { run: deleteDepartment } = useRequest(deleteAPI, {
      manual: true,
    })

    // 删除数据

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
      formRef.Namespace = record.namespace
    }
    const deleteFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getListData(<ListParams>{
          page: currentPage.value,
          size: currentSize.value,
          name_filter: currentNameFilter.value,
        })
      })
    }

    // 提交表单
    const commitModal = () => {
      switch (modalStatus.value) {
        case ModalStatus.add:
          addDepartment(<UpdateParams>{
            Name: formRef.Name,
            Namespace: formRef.Namespace,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name_filter: currentNameFilter.value,
            })
            modalVisible.value = false
          })
          break
        case ModalStatus.edit:
          editDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            Name: formRef.Name,
            Namespace: formRef.Namespace,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name_filter: currentNameFilter.value,
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
    }
  },
})
</script>

<style lang="less" scoped>
</style>