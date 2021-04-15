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
            placeholder="输入需要查询的实验名称"
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
              <a-button type="primary" @click="editFn(record)"
                >修改配置</a-button
              >
              <a-button
                type="primary"
                @click="openMissionCheckpointManager(record.id)"
                >修改考点</a-button
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

    <!-- 实验文档编辑 -->
    <a-modal
      v-model:visible="instructionsVisible"
      title="实验文档"
      :footer="null"
      :afterClose="instructionsTipAfterClose"
      width="920px"
    >
      <a-skeleton v-if="instructionsLoading" :active="true" />
      <v-md-editor
        v-model="guideData"
        height="800px"
        v-if="!instructionsLoading"
        @save="guideSaver"
      >
      </v-md-editor>
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
  name: string
}

type ListResult = {
  id: number
  total: number
  name: string
  desc: string
  containers: string[]
  deployment: number
  exec_container: string
  command: string
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  Name: string
  desc: string
  total: number
  containers: string[]
  deployment: number
  exec_container: string
  command: string
}

enum ModalStatus {
  add = 1,
  edit,
}

const listDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.ms.manageList, {
    params: params,
  })
}

const countDataAPI = (params: ListParams) => {
  return defaultClient.get<BaseResponse>(paths.ms.count, {
    params: params,
  })
}

const addAPI = (params: UpdateParams) => {
  return defaultClient.post<BaseResponse>(paths.ms.add, params)
}

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(paths.ms.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(paths.ms.delete + params + '/')
}

const getGuideAPI = (params: number) => {
  return defaultClient.get<BaseResponse>(paths.ms.getGuide + params + '/')
}

interface updateGuideParams {
  id: number
  text: string
}

const updateGuideAPI = (params: updateGuideParams) => {
  return defaultClient.put<BaseResponse>(paths.ms.updateGuide, params)
}

type optionsResult = {
  id: number
  name: string
}[]

const deploymentQuickAPI = () => {
  return defaultClient.get<BaseResponse>(paths.deployment.quick)
}

export default defineComponent({
  components: {
    missionCheckpointManager,
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

    const getListParams = (): ListParams => {
      return <ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        name: currentNameFilter.value,
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
      currentPage.value = 1
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
      ID: '',
      Name: '',
      desc: '',
      total: 100,
      containers: [],
      deployment: 0,
      exec_container: '',
      command: '',
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.ID = ''
      formRef.Name = ''
      formRef.desc = ''
      formRef.total = 100
      formRef.containers = []
      formRef.deployment = 0
      formRef.exec_container = ''
      formRef.command = ''
    }
    const getUpdateParams = (): UpdateParams => {
      return <UpdateParams>{
        ID: Number(formRef.ID),
        Name: formRef.Name,
        desc: formRef.desc,
        total: formRef.total,
        containers: formRef.containers,
        deployment: formRef.deployment,
        exec_container: formRef.exec_container,
        command: formRef.command,
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

    // deployment选项
    const { data: deploymentOptions } = useRequest(deploymentQuickAPI, {
      formatResult: (res): optionsResult => {
        return <optionsResult>res.data.Data
      },
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
      formRef.desc = record.desc
      formRef.total = record.total
      formRef.containers = record.containers
      formRef.deployment = record.deployment
      formRef.exec_container = record.exec_container
      formRef.command = record.command
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

    // 指南文档
    const instructionsVisible = ref<boolean>(false)
    const {
      data: guideData,
      run: getGuideData,
      loading: instructionsLoading,
    } = useRequest(getGuideAPI, {
      manual: true,
      formatResult: (res): string => {
        return <string>res.data.Data
      },
    })
    const nowGuideEditorID = ref<number>(0)
    const openGuideEditor = (id: number) => {
      nowGuideEditorID.value = id
      getGuideData(id)
      instructionsVisible.value = true
    }
    const instructionsTipAfterClose = () => {
      instructionsVisible.value = false
      guideData.value = `🤪无实验文档数据，请联系刷新页面或实验教师`
    }
    const { run: updateGuide } = useRequest(updateGuideAPI, {
      manual: true,
    })
    const guideSaver = (text: string) => {
      updateGuide(<updateGuideParams>{ id: nowGuideEditorID.value, text: text })
    }

    // 检查点编辑和成绩查询
    const targetMissionID = ref<number>(0)
    const missionCheckpointVisiable = ref<boolean>(false)
    const openMissionCheckpointManager = (id: number) => {
      targetMissionID.value = id
      missionCheckpointVisiable.value = true
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
      openGuideEditor,
      guideData,
      instructionsVisible,
      instructionsLoading,
      instructionsTipAfterClose,
      guideSaver,
      deploymentOptions,
      missionCheckpointVisiable,
      targetMissionID,
      openMissionCheckpointManager,
    }
  },
})
</script>

<style lang="less" scoped>
</style>