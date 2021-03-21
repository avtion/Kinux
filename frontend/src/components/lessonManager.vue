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
        <a-button @click="addFn" style="margin-bottom: 8px">新增课程</a-button>
        <!-- 表格 -->
        <a-table
          :columns="columns"
          :data-source="listData"
          :pagination="pagination"
          :loading="isListDataLoading"
          row-key="id"
          @change="handleTableChange"
          class="ant-table-striped"
          :rowClassName="
            (record, index) => (index % 2 === 1 ? 'table-striped' : null)
          "
          bordered
          tableLayout="auto"
        >
          <!-- 实验管理 -->
          <template #mission="{ record }">
            <a-button
              type="default"
              size="small"
              @click="openLessonMission(record)"
              >编辑实验
            </a-button>
          </template>

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
              <a-button type="primary" @click="editFn(record)"
                >修改资料
              </a-button>
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
        <a-form-item label="课程ID">
          <a-input
            v-model:value="formRef.ID"
            :disabled="true"
            style="width: 300px"
          />
        </a-form-item>
        <a-form-item label="课程名称">
          <a-input v-model:value="formRef.Name" style="width: 300px" />
        </a-form-item>
        <a-form-item label="课程描述">
          <a-textarea
            v-model:value="formRef.Desc"
            style="width: 300px"
            auto-size
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 实验管理 -->
    <a-modal
      v-model:visible="lessonMissionVisible"
      title="课程实验管理"
      :footer="null"
      width="920px"
      :destroyOnClose="true"
    >
      <lessonMission :lessonID="lessonMissionLessonID"></lessonMission>
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
import { BaseResponse, defaultClient } from '@/apis/request'

// 实验管理
import lessonMission from '@/components/lesson_mission.vue'

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
    title: '实验管理',
    dataIndex: 'id',
    slots: { customRender: 'mission' },
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
  desc: string
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  Name: string
  Desc: string
}

enum ModalStatus {
  add = 1,
  edit,
}

const apiPath = {
  list: '/v2/lesson/list',
  count: '/v2/lesson/count',
  add: '/v2/lesson/',
  edit: '/v2/lesson/',
  delete: '/v2/lesson/',
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

const editAPI = (params: UpdateParams) => {
  return defaultClient.put<BaseResponse>(apiPath.edit, params)
}

const deleteAPI = (params: number) => {
  return defaultClient.delete<BaseResponse>(apiPath.delete + params + '/')
}

export default defineComponent({
  components: {
    SmileOutlined,
    DownOutlined,
    lessonMission,
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
      ID: '',
      Name: '',
      Desc: '',
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.ID = ''
      formRef.Name = ''
      formRef.Desc = ''
    }

    // 添加数据
    const { run: addDepartment, loading: isAddLoading } = useRequest(addAPI, {
      defaultParams: [
        <UpdateParams>{
          ID: Number(formRef.ID),
          Name: formRef.Name,
          Desc: formRef.Desc,
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
          addDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            Name: formRef.Name,
            Desc: formRef.Desc,
          }).finally(() => {
            getListData(getListParams())
            modalVisible.value = false
          })
          break
        case ModalStatus.edit:
          editDepartment(<UpdateParams>{
            ID: Number(formRef.ID),
            Name: formRef.Name,
            Desc: formRef.Desc,
          }).finally(() => {
            getListData(getListParams())
            modalVisible.value = false
          })
          break
      }
    }

    // 课程实验管理
    const lessonMissionVisible = ref<boolean>(false)
    const lessonMissionLessonID = ref<number>(0)
    const openLessonMission = (record: ListResult) => {
      lessonMissionLessonID.value = record.id
      lessonMissionVisible.value = true
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

      // 课程实验管理
      lessonMissionVisible,
      lessonMissionLessonID,
      openLessonMission,
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