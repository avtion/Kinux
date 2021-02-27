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
            placeholder="输入需要查询的配置"
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
        <a-form-item label="配置编号">
          <a-input
            v-model:value="formRef.ID"
            :disabled="true"
            style="width: 400px"
          />
        </a-form-item>
        <a-form-item label="配置名称">
          <a-input v-model:value="formRef.Name" style="width: 400px" />
        </a-form-item>
        <span>配置详细: </span>
        <v-ace-editor
          v-model:value="formRef.Raw"
          lang="yaml"
          theme="github"
          style="height: 500px"
        />
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

// vue3-ace-editor
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
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
    title: '配置名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
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
  name: string
}

type ListResult = {
  id: number
  name: string
  raw: string
  created_at: string
  updated_at: string
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  Name: string
  Raw: string
}

enum ModalStatus {
  add = 1,
  edit,
}

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.deployment.list, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.deployment.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(paths.deployment.add, params)
}

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(paths.deployment.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(
    paths.deployment.delete + params + '/'
  )
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
        },
      ],
      formatResult: (res): ListResults => {
        getTotal(<ListParams>{
          name: currentNameFilter.value,
        })
        return <ListResults>res.data.Data
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
      })
    }

    // 追加和更新窗口
    const modalVisible = ref<boolean>(false)
    const modalTitle = ref<string>('添加')
    const modalStatus = ref<ModalStatus>(ModalStatus.add)
    const formRef = reactive({
      ID: '',
      Name: '',
      Raw: '',
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.ID = ''
      formRef.Name = ''
      formRef.Raw = ''
    }

    // 添加数据
    const { run: addDepartment, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [
        <UpdateParams>{
          ID: Number(formRef.ID),
          Name: formRef.Name,
          Raw: formRef.Raw,
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
            Raw: formRef.Raw,
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
      formRef.Raw = record.raw
    }
    const deleteFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getListData(<ListParams>{
          page: currentPage.value,
          size: currentSize.value,
          name: currentNameFilter.value,
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
            Raw: formRef.Raw,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name: currentNameFilter.value,
            })
            modalVisible.value = false
          })
          break
        case ModalStatus.edit:
          editDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            Name: formRef.Name,
            Raw: formRef.Raw,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name: currentNameFilter.value,
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