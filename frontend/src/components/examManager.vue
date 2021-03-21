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
          <a-select
            v-model:value="currentNamespaceFilter"
            mode="multiple"
            style="width: 200px"
            placeholder="请选择需要查询的课程"
            @change="onSearch"
          >
            <a-select-option v-for="v in namespaceOptions" :key="v">{{
              v
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
          <!-- 指南文档 -->
          <template #guide="{ record }">
            <span>
              <a-button size="small" @click="openGuideEditor(record.id)"
                >点击编辑文档</a-button
              >
            </span>
          </template>

          <!-- 展开描述 -->
          <template #expandedRowRender="{ record }">
            <p style="margin: 0">描述：{{ record.desc }}</p>
          </template>

          <!-- 编辑框 -->
          <template #operation="{ record }">
            <a-button-group size="small">
              <a-button type="primary" @click="openMissionScoreManager(record)"
                >成绩查询</a-button
              >
              <a-button type="primary" @click="editFn(record)"
                >修改配置</a-button
              >
              <a-button
                type="primary"
                @click="openMissionCheckpointManager(record.id)"
                >修改检查点</a-button
              >
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
        <a-form-item label="实验编号">
          <a-input
            v-model:value="formRef.ID"
            :disabled="true"
            style="width: 300px"
          />
        </a-form-item>
        <a-form-item label="实验名称" required>
          <a-input v-model:value="formRef.Name" style="width: 300px" />
        </a-form-item>
        <a-form-item label="实验描述">
          <a-textarea
            v-model:value="formRef.desc"
            style="width: 300px"
            placeholder="请填写实验相关描述"
            :auto-size="{ minRows: 2 }"
          />
        </a-form-item>
        <a-form-item label="命名空间">
          <a-input
            v-model:value="formRef.namespace"
            style="width: 300px"
            placeholder="default"
          />
        </a-form-item>
        <a-form-item label="实验总分" required>
          <a-input-number v-model:value="formRef.total" style="width: 300px" />
        </a-form-item>
        <a-form-item label="实验配置" required>
          <a-select v-model:value="formRef.deployment" style="width: 300px">
            <a-select-option :key="0">请选择实验配置</a-select-option>
            <a-select-option v-for="v in deploymentOptions" :key="v.id">{{
              v.name
            }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="默认容器">
          <a-input
            v-model:value="formRef.exec_container"
            style="width: 300px"
            placeholder="默认为实验配置的首个容器"
          />
        </a-form-item>
        <a-form-item label="默认命令">
          <a-input
            v-model:value="formRef.command"
            style="width: 300px"
            placeholder="bash"
          />
        </a-form-item>
        <a-form-item label="可用容器">
          <a-select
            v-model:value="formRef.containers"
            mode="tags"
            style="width: 300px"
            placeholder="默认放行全部容器可用"
          >
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 成绩查询 -->
    <a-modal
      v-model:visible="missionScoreVisiable"
      title="成绩查询"
      :footer="null"
      width="920px"
      :destroyOnClose="true"
    >
      <missionScoreManager
        :missionID="targetMissionID"
        :namespace="targetNamepsce"
      ></missionScoreManager>
    </a-modal>

    <!-- 实验检查点编辑 -->
    <a-modal
      v-model:visible="missionCheckpointVisiable"
      title="实验检查点编辑"
      :footer="null"
      width="920px"
      :destroyOnClose="true"
    >
      <missionCheckpointManager
        :missionID="targetMissionID"
      ></missionCheckpointManager>
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

// 实验检查点管理
import missionCheckpointManager from '@/components/missionCheckpointManager.vue'

// 成绩查询
import missionScoreManager from '@/components/missionScore.vue'

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
    title: '实验名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '命名空间',
    dataIndex: 'namespace',
    key: 'namespace',
  },
  {
    title: '实验文档',
    dataIndex: 'id',
    slots: { customRender: 'guide' },
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
  ns: string[]
}

type ListResult = {
  id: number
  name: string
  desc: string
  namespace: string
  total: number
  force_order: boolean
  begin_at: string
  end_at: string
  created_at: string
  time_limit: number
}

type ListResults = ListResult[]

type UpdateParams = {
  id: number
  name: string
  desc: string
  namespace: string
  total: number
  begin_at: number
  end_at: number
  force_order: boolean
  time_limit: number
}

enum ModalStatus {
  add = 1,
  edit,
}

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.exam.list, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.exam.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(paths.exam.add, params)
}

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(paths.exam.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(paths.exam.delete + params + '/')
}

const namespacesAPI = (params: number) => {
  return defaultClient.get<BaseResponse>(paths.exam.ns)
}

interface updateGuideParams {
  id: number
  text: string
}

type optionsResult = {
  id: number
  name: string
}[]

export default defineComponent({
  components: {
    missionCheckpointManager,
    missionScoreManager,
    SmileOutlined,
    DownOutlined,
  },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 获取命名空间选项
    const currentNamespaceFilter = ref<string[]>([])
    const { data: namespaceOptions } = useRequest(namespacesAPI, {
      formatResult: (res): string[] => {
        return <string[]>res.data.Data
      },
    })

    // 表格
    const currentPage = ref<number>(1)
    const currentSize = ref<number>(10)

    const getListParams = (): ListParams => {
      return <ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        ns: currentNamespaceFilter.value,
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
      id: '',
      name: '',
      desc: '',
      namespace: '',
      total: 0,
      begin_at: 0,
      end_at: 0,
      force_order: false,
      time_limit: 0,
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.id = ''
      formRef.name = ''
      formRef.desc = ''
      formRef.namespace = ''
      formRef.total = 0
      formRef.begin_at = 0
      formRef.end_at = 0
      formRef.force_order = false
      formRef.time_limit = 0
    }
    const getUpdateParams = (): UpdateParams => {
      return <UpdateParams>{
        id: Number(formRef.id),
        name: formRef.name,
        desc: formRef.desc,
        namespace: formRef.namespace,
        total: formRef.total,
        begin_at: formRef.begin_at,
        end_at: formRef.end_at,
        force_order: formRef.force_order,
        time_limit: formRef.time_limit,
      }
    }

    // 添加数据
    const { run: addDepartment, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [getUpdateParams()],
      manual: true,
    })

    // 修改数据
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

    const addFn = () => {
      modalVisible.value = true
      modalTitle.value = '添加'
      modalStatus.value = ModalStatus.add
    }
    const editFn = (record: ListResult) => {
      modalVisible.value = true
      modalTitle.value = '修改'
      modalStatus.value = ModalStatus.edit
      formRef.desc = record.desc
      formRef.namespace = record.namespace
      formRef.total = record.total
    }
    const deleteFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getListData(getListParams())
      })
    }

    // 提交表单
    const commitModal = () => {
      switch (modalStatus.value) {
        case ModalStatus.add:
          addDepartment(getUpdateParams()).finally(() => {
            getListData(getListParams())
            modalVisible.value = false
          })
          break
        case ModalStatus.edit:
          editDepartment(getUpdateParams()).finally(() => {
            getListData(getListParams())
            modalVisible.value = false
          })
          break
      }
    }

    // 检查点编辑和成绩查询
    const targetMissionID = ref<number>(0)
    const targetNamepsce = ref<string>('')
    const missionCheckpointVisiable = ref<boolean>(false)
    const missionScoreVisiable = ref<boolean>(false)
    const openMissionCheckpointManager = (id: number) => {
      targetMissionID.value = id
      missionCheckpointVisiable.value = true
    }
    const openMissionScoreManager = (record: ListResult) => {
      targetMissionID.value = record.id
      targetNamepsce.value = record.namespace
      missionScoreVisiable.value = true
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
      editFn,
      deleteFn,
      modalVisible,
      modalTitle,
      commitModal,
      isAddLoading,
      isEditLoading,
      formRef,
      clearForm,
      currentNamespaceFilter,
      namespaceOptions,
      missionCheckpointVisiable,
      targetMissionID,
      openMissionCheckpointManager,
      missionScoreVisiable,
      openMissionScoreManager,
      targetNamepsce,
    }
  },
})
</script>

<style lang="less" scoped>
</style>