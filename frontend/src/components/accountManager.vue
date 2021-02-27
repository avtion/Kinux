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
          <!-- 用户名筛选 -->
          <a-input-search
            placeholder="输入需要查询的用户名"
            style="width: 200px"
            v-model:value="currentNameFilter"
            @search="onSearch"
          />
          <a-divider type="vertical" />
          <!-- 班级筛选 -->
          <a-select
            style="width: 140px"
            v-model:value="currentDepartmentFilter"
            @change="onSearch"
          >
            <a-select-option :value="0">班级筛选</a-select-option>
            <a-select-option v-for="v in departmentOptions" :key="v.id">{{
              v.name
            }}</a-select-option>
          </a-select>
          <a-divider type="vertical" />
          <!-- 角色筛选 -->
          <a-select
            style="width: 140px"
            v-model:value="currentRoleFilter"
            @change="onSearch"
          >
            <a-select-option :value="0">权限筛选</a-select-option>
            <a-select-option v-for="v in roleOptions" :key="v.id">{{
              v.name
            }}</a-select-option>
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
          <template #role="{ record }">
            <span>
              <a-tag :color="record.role_id == 2 ? 'green' : 'red'">
                {{ record.role }}
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
        <a-form-item label="用户编号">
          <a-input
            v-model:value="formRef.ID"
            :disabled="true"
            style="width: 300px"
          />
        </a-form-item>

        <a-form-item label="用户名" required>
          <a-input style="width: 300px" v-model:value="formRef.username" />
        </a-form-item>
        <a-form-item label="真实姓名">
          <a-input style="width: 300px" v-model:value="formRef.real_name" />
        </a-form-item>
        <a-form-item label="密码" required>
          <a-input
            style="width: 300px"
            type="password"
            v-model:value="formRef.password"
          />
        </a-form-item>
        <a-form-item label="权限" required>
          <a-select style="width: 300px" v-model:value="formRef.role">
            <a-select-option :value="0">无</a-select-option>
            <a-select-option v-for="v in roleOptions" :key="v.id">{{
              v.name
            }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="班级" required>
          <a-select style="width: 300px" v-model:value="formRef.department">
            <a-select-option :value="0">无</a-select-option>
            <a-select-option v-for="v in departmentOptions" :key="v.id">{{
              v.name
            }}</a-select-option>
          </a-select>
        </a-form-item>
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
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
  },
  {
    title: '真实姓名',
    dataIndex: 'real_name',
    key: 'real_name',
  },
  {
    title: '权限',
    dataIndex: 'id',
    key: 'id',
    slots: { customRender: 'role' },
  },
  {
    title: '班级名称',
    dataIndex: 'department',
    key: 'department',
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

type ListParams = {
  page: number
  size: number
  name: string
  role: number
  department: number
}

type ListResult = {
  id: number
  profile: number
  role: string
  username: string
  real_name: string
  department: string
  created_at: string
  role_id: number
  department_id: number
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  username: string
  real_name: string
  password: string
  role: number
  department: number
}

enum ModalStatus {
  add = 1,
  edit,
}

type OptionResult = {
  id: number
  name: string
}[]

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.ac.list, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.ac.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(paths.ac.add, params)
}

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(paths.ac.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(paths.ac.delete + params + '/')
}

const departmentOptionAPI = () => {
  return defaultClient.get<BaseResponse>(paths.department.quick)
}

const roleOptionAPI = () => {
  return defaultClient.get<BaseResponse>(paths.role.quick)
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

    // 选项
    const currentDepartmentFilter = ref<number>(0)
    const { data: departmentOptions } = useRequest(departmentOptionAPI, {
      formatResult: (res): OptionResult => {
        return <OptionResult>res.data.Data
      },
    })
    const currentRoleFilter = ref<number>(0)
    const { data: roleOptions } = useRequest(roleOptionAPI, {
      formatResult: (res): OptionResult => {
        return <OptionResult>res.data.Data
      },
    })

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
          role: currentRoleFilter.value,
          department: currentDepartmentFilter.value,
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
    const onSearch = (value?: string) => {
      getListData(<ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        name: currentNameFilter.value,
        role: currentRoleFilter.value,
        department: currentDepartmentFilter.value,
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
        role: currentRoleFilter.value,
        department: currentDepartmentFilter.value,
      })
    }

    // 追加和更新窗口
    const modalVisible = ref<boolean>(false)
    const modalTitle = ref<string>('添加')
    const modalStatus = ref<ModalStatus>(ModalStatus.add)
    const formRef = reactive({
      ID: '',
      username: '',
      real_name: '',
      password: '',
      role: 0,
      department: 0,
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.ID = ''
      formRef.username = ''
      formRef.real_name = ''
      formRef.password = ''
      formRef.role = 0
      formRef.department = 0
    }

    // 添加数据
    const { run: addDepartment, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [
        <UpdateParams>{
          ID: Number(formRef.ID),
          username: formRef.username,
          real_name: formRef.real_name,
          password: formRef.password,
          role: formRef.role,
          department: formRef.department,
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
            username: formRef.username,
            real_name: formRef.real_name,
            password: formRef.password,
            role: formRef.role,
            department: formRef.department,
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
      formRef.username = record.username
      formRef.real_name = record.real_name
      formRef.role = record.role_id
      formRef.department = record.department_id
    }
    const deleteFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getListData(<ListParams>{
          page: currentPage.value,
          size: currentSize.value,
          name: currentNameFilter.value,
          role: currentRoleFilter.value,
          department: currentDepartmentFilter.value,
        })
      })
    }

    // 提交表单
    const commitModal = () => {
      switch (modalStatus.value) {
        case ModalStatus.add:
          addDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            username: formRef.username,
            real_name: formRef.real_name,
            password: formRef.password,
            role: formRef.role,
            department: formRef.department,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name: currentNameFilter.value,
              role: currentRoleFilter.value,
              department: currentDepartmentFilter.value,
            })
            modalVisible.value = false
          })
          break
        case ModalStatus.edit:
          editDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            username: formRef.username,
            real_name: formRef.real_name,
            password: formRef.password,
            role: formRef.role,
            department: formRef.department,
          }).finally(() => {
            getListData(<ListParams>{
              page: currentPage.value,
              size: currentSize.value,
              name: currentNameFilter.value,
              role: currentRoleFilter.value,
              department: currentDepartmentFilter.value,
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
      currentDepartmentFilter,
      currentRoleFilter,
      departmentOptions,
      roleOptions,
    }
  },
})
</script>

<style lang="less" scoped>
</style>