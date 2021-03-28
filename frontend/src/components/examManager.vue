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
            style="width: 180px"
            placeholder="输入课程的名字搜索"
            show-search
            @search="searchMissionListOptions"
            :show-arrow="false"
            :filter-option="false"
            :loading="isMissionOpitonsLoading"
            v-model:value="currentLessonFilter"
          >
            <a-select-option
              v-for="v in missionListData"
              :key="v.id"
              :title="v.name"
              >{{ v.id }} - {{ v.name }}
            </a-select-option>
          </a-select>
          <a-divider type="vertical" />
          <a-button
            @click="addFn"
            :disabled="
              currentLessonFilter === undefined || currentLessonFilter === 0
            "
            >新增考试{{
              currentLessonFilter === undefined || currentLessonFilter === 0
                ? '(需先选课程)'
                : ''
            }}</a-button
          >
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
          @expand="expandRowFn"
        >
          <!-- 展开描述 -->
          <template #expandedRowRender="{ index }">
            <!-- 考试实验表格 -->
            <a-table
              :pagination="false"
              :columns="examMissionsColumns"
              :data-source="listData[index].missions"
              rowKey="id"
            >
              <!-- 编辑框 -->
              <template #operation="{ record }">
                <a-button-group size="small">
                  <a-button type="primary" @click="editMissionFn(record)">
                    修改实验
                  </a-button>
                  <a-button type="default">自定义考点 </a-button>

                  <a-popconfirm
                    placement="top"
                    ok-text="是"
                    cancel-text="否"
                    @confirm="deleteExamMissionFn(record)"
                  >
                    <template #title>
                      <p>确定删除？</p>
                    </template>
                    <a-button type="danger">删除</a-button>
                  </a-popconfirm>
                </a-button-group>
              </template>
            </a-table>
          </template>

          <!-- 编辑框 -->
          <template #operation="{ record }">
            <a-button-group size="small">
              <a-button type="primary" @click="addMissionFn(record)"
                >新增实验
              </a-button>
              <a-button type="default" @click="editFn(record)"
                >修改考试</a-button
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
        <a-form-item label="考试编号">
          <a-input
            v-model:value="formRef.ID"
            :disabled="true"
            style="width: 300px"
          />
        </a-form-item>
        <a-form-item label="考试名称" required>
          <a-input v-model:value="formRef.name" style="width: 300px" />
        </a-form-item>
        <a-form-item label="考试总分" required>
          <a-input-number v-model:value="formRef.total" style="width: 300px" />
        </a-form-item>
        <a-form-item label="考试时间" required>
          <a-range-picker
            :show-time="{ format: 'HH:mm' }"
            format="YYYY-MM-DD HH:mm"
            :placeholder="['开始时间', '结束时间']"
            style="width: 300px"
            v-model:value="formRef.begin_end"
          />
        </a-form-item>
        <a-form-item label="时长限制" required>
          <a-time-picker v-model:value="formRef.time_limit" />
        </a-form-item>
        <a-form-item label="按顺序完成">
          <a-switch v-model:checked="formRef.force_order" />
        </a-form-item>
        <a-form-item label="考试描述">
          <a-textarea
            v-model:value="formRef.desc"
            style="width: 300px"
            placeholder="请填写考试的描述"
            :auto-size="{ minRows: 2 }"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 新增实验和修改实验 -->
    <a-modal
      :title="missionModalTitle"
      v-model:visible="missionModalVisible"
      :confirm-loading="isAddMissionLoading || isEditMissionLoading"
      @ok="commitMissionModal"
      :afterClose="clearMissionForm"
    >
      <a-form>
        <a-form-item label="记录编号">
          <a-input
            v-model:value="missionFormRef.id"
            :disabled="true"
            style="width: 300px"
          />
        </a-form-item>
        <a-form-item label="考试编号">
          <a-input
            v-model:value="missionFormRef.exam"
            :disabled="true"
            style="width: 300px"
          />
        </a-form-item>
        <a-form-item label="实验项目" required>
          <a-select
            show-search
            v-model:value="missionFormRef.mission"
            placeholder="输入需要搜索的实验项目名"
            style="width: 300px"
            :default-active-first-option="false"
            :show-arrow="false"
            :filter-option="false"
            not-found-content="无对应的实验"
            @search="getMissionOptions"
          >
            <a-select-option v-for="v in missionOptions" :key="v.id">
              {{ v.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
          :label="'成绩占比(' + '剩余可用:' + examRestPercent + ')'"
          required
        >
          <a-input-number
            v-model:value="missionFormRef.percent"
            style="width: 200px"
            :min="1"
            :max="
              missionModalStatus == ModalStatus.edit ? 100 : examRestPercent
            "
          />
        </a-form-item>
        <a-form-item label="排序权重" required>
          <a-input-number
            v-model:value="missionFormRef.priority"
            style="width: 300px"
          />
        </a-form-item>
      </a-form>
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

// 实验检查点管理
import missionCheckpointManager from '@/components/missionCheckpointManager.vue'
import * as moment from 'moment'

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
    title: '考试名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '开始时间',
    dataIndex: 'begin_at',
    key: 'begin_at',
  },
  {
    title: '结束时间',
    dataIndex: 'end_at',
    key: 'end_at',
  },
  {
    title: '时长',
    dataIndex: 'time_limit',
    key: 'time_limit',
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
  lesson: number
}

type ListResult = {
  id: number
  name: string
  desc: string
  total: number
  force_order: boolean
  begin_at: string
  end_at: string
  created_at: string
  time_limit: string
  begin_at_unix: number
  end_at_unix: number
  created_at_unix: number
  time_limit_unix: number
  missions: examMission[]
}

// 实验数据
type examMission = {
  id: number
  exam: number
  mission: string
  percent: number
  priority: number
  mission_id: number
}

type ListResults = ListResult[]

type UpdateParams = {
  ID: number
  name: string
  desc: string
  begin_at: number
  end_at: number
  force_order: boolean
  time_limit: number
  lesson: number
}

enum ModalStatus {
  add = 1,
  edit,
}

const apiPath = {
  list: '/v1/exam/list/',
  count: '/v1/exam/count/',
  add: '/v1/exam/',
  edit: '/v1/exam/',
  delete: '/v1/exam/',
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
    const currentLessonFilter = ref<number>()

    // 课程筛选
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
        if (
          (currentLessonFilter.value === undefined ||
            currentLessonFilter.value === 0) &&
          (<missionListResult[]>res.data.Data).length > 0
        ) {
          currentLessonFilter.value = (<missionListResult[]>res.data.Data)[0].id
        }
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
    watch(
      () => currentLessonFilter.value,
      () => {
        getListData(getListParams())
      }
    )

    // 获取考试数据
    const getListParams = (): ListParams => {
      return <ListParams>{
        page: currentPage.value,
        size: currentSize.value,
        lesson: currentLessonFilter.value,
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
        const resData = <ListResults>res.data.Data
        resData.forEach((v) => [(v.missions = <examMission[]>[])])
        return resData
      },
      manual: true,
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
      name: '',
      desc: '',
      begin_end: [moment(), moment().add(1, 'd')],
      force_order: false,
      time_limit: moment(40, 'minutes'),
      total: 100,
      lesson: currentLessonFilter.value,
    })
    const clearForm = () => {
      console.log('清理数据')
      formRef.ID = ''
      formRef.name = ''
      formRef.desc = ''
      formRef.begin_end = [moment(), moment().add(1, 'd')]
      formRef.force_order = false
      formRef.time_limit = moment(40, 'minutes')
      formRef.total = 100
    }
    const getUpdateParams = (): UpdateParams => {
      return <UpdateParams>{
        ID: Number(formRef.ID),
        name: formRef.name,
        desc: formRef.desc,
        begin_at: formRef.begin_end[0].valueOf(),
        end_at: formRef.begin_end[1].valueOf(),
        force_order: formRef.force_order,
        time_limit: moment
          .duration({
            h: formRef.time_limit.hour(),
            m: formRef.time_limit.minute(),
            s: formRef.time_limit.second(),
          })
          .asSeconds(),
        lesson: currentLessonFilter.value,
        total: formRef.total,
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
      formRef.name = record.name
      formRef.desc = record.desc

      formRef.begin_end = [
        moment.unix(record.begin_at_unix),
        moment.unix(record.end_at_unix),
      ]
      formRef.force_order = record.force_order
      const t = moment.duration({ s: record.time_limit_unix })
      formRef.time_limit = moment({
        h: t.hours(),
        m: t.minutes(),
        s: t.seconds(),
      })

      formRef.total = record.total
      formRef.lesson = currentLessonFilter.value
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
    const missionCheckpointVisiable = ref<boolean>(false)
    const openMissionCheckpointManager = (id: number) => {
      targetMissionID.value = id
      missionCheckpointVisiable.value = true
    }

    // 考试的实验数据

    const examMissionsColumns = [
      {
        title: '展示优先级',
        dataIndex: 'priority',
        key: 'priority',
      },
      {
        title: '实验名称',
        dataIndex: 'mission',
        key: 'mission',
      },
      {
        title: '成绩占比',
        dataIndex: 'percent',
        key: 'percent',
      },
      {
        title: '操作',
        dataIndex: 'id',
        slots: { customRender: 'operation' },
      },
    ]
    const examMissionApiPath = {
      list: '/v1/em/',
      count: '/v1/em/count/',
      add: '/v1/em/',
      edit: '/v1/em/',
      delete: '/v1/em/',
      used: '/v1/em/percent/',
    }

    const listExamMissionAPI = (params: number) => {
      return defaultClient.get<BaseResponse>(examMissionApiPath.list, {
        params: {
          exam: params,
        },
      })
    }

    // 修改和删除
    type missionUpdateParams = {
      id: number
      exam: number
      mission: number
      percent: number
      priority: number
    }
    const addExamMissionAPI = (params: missionUpdateParams) => {
      return defaultClient.post<BaseResponse>(examMissionApiPath.add, params)
    }

    const editExamMissionAPI = (params: missionUpdateParams) => {
      return defaultClient.put<BaseResponse>(examMissionApiPath.edit, params)
    }

    const deleteExamMissionAPI = (params: number) => {
      return defaultClient.delete<BaseResponse>(
        examMissionApiPath.delete + params + '/'
      )
    }
    const { run: deleteExamMission } = useRequest(deleteExamMissionAPI, {
      manual: true,
    })

    const deleteExamMissionFn = (record: ListResult) => {
      deleteExamMission(record.id).finally(() => {
        getListData(getListParams())
      })
    }

    const { run: getExamMission } = useRequest(listExamMissionAPI, {
      manual: true,
      formatResult: (res): examMission[] => {
        return res.data.Data
      },
    })

    // 展开子行的函数
    const expandRowFn = (expanded: false, record: ListResult) => {
      if (!expanded) {
        return
      }
      getExamMission(record.id).then((res: examMission[]) => {
        const index = listData.value.findIndex((v) => {
          return v.id === record.id
        })
        listData.value[index].missions = res
        return
      })
    }

    // 新增和修改实验
    const missionModalVisible = ref<boolean>(false)
    const missionModalTitle = ref<string>('添加')
    const missionModalStatus = ref<ModalStatus>(ModalStatus.add)

    const missionFormRef = reactive({
      id: '',
      exam: 0,
      mission: undefined,
      percent: 1,
      priority: 0,
    })
    const clearMissionForm = () => {
      console.log('清理数据')
      missionFormRef.id = ''
      missionFormRef.exam = 0
      missionFormRef.mission = undefined
      missionFormRef.percent = 1
      missionFormRef.priority = 0
    }
    const getMissionUpdateParams = (): missionUpdateParams => {
      return <missionUpdateParams>{
        id: Number(missionFormRef.id),
        exam: missionFormRef.exam,
        mission: missionFormRef.mission,
        percent: missionFormRef.percent,
        priority: missionFormRef.priority,
      }
    }

    // 添加数据
    const { run: addMission, loading: isAddMissionLoading } = useRequest(
      addExamMissionAPI,
      {
        defaultParams: [getMissionUpdateParams()],
        manual: true,
      }
    )

    // 修改数据
    const { run: editMission, loading: isEditMissionLoading } = useRequest(
      editExamMissionAPI,
      {
        defaultParams: [getMissionUpdateParams()],
        manual: true,
      }
    )
    const { run: deleteMission } = useRequest(deleteAPI, {
      manual: true,
    })

    const addMissionFn = (record: ListResult) => {
      missionModalVisible.value = true
      missionModalTitle.value = '添加'
      missionModalStatus.value = ModalStatus.add

      missionFormRef.exam = record.id
      getExamUsedPercent(record.id)
    }
    const editMissionFn = (record: examMission) => {
      missionModalVisible.value = true
      missionModalTitle.value = '修改'
      missionModalStatus.value = ModalStatus.edit

      missionFormRef.id = record.id + ''
      missionFormRef.exam = record.exam
      missionFormRef.mission = record.mission_id
      missionFormRef.percent = record.percent
      missionFormRef.priority = record.priority
      getExamUsedPercent(record.exam)
    }
    const deleteMissionFn = (record: ListResult) => {
      deleteDepartment(record.id).finally(() => {
        getListData(getListParams())
      })
    }

    // 提交表单
    const commitMissionModal = () => {
      switch (missionModalStatus.value) {
        case ModalStatus.add:
          addMission(getMissionUpdateParams()).finally(() => {
            getExamMission(missionFormRef.exam).then((res: examMission[]) => {
              const index = listData.value.findIndex((v) => {
                return v.id === missionFormRef.exam
              })
              listData.value[index].missions = res
              return
            })
            missionModalVisible.value = false
          })
          break
        case ModalStatus.edit:
          editMission(getMissionUpdateParams()).finally(() => {
            getExamMission(missionFormRef.exam).then((res: examMission[]) => {
              const index = listData.value.findIndex((v) => {
                return v.id === missionFormRef.exam
              })
              listData.value[index].missions = res
              return
            })
            missionModalVisible.value = false
          })
          break
      }
    }

    // 获取实验选项
    const missionOptionsAPI = (params?: string) => {
      return defaultClient.get<BaseResponse>(paths.ms.manageList, {
        params: {
          name: params,
        },
      })
    }

    const {
      data: missionOptions,
      run: getMissionOptions,
      loading: isMissionOptionsLoading,
    } = useRequest(missionOptionsAPI, {
      manual: false,
      formatResult: (res): { id: number; name: string }[] => {
        return res.data.Data
      },
    })

    // 剩余成绩比例
    const percentAPI = (params: number) => {
      return defaultClient.get<BaseResponse>(
        examMissionApiPath.used + params + '/'
      )
    }

    const { data: examUsedPercent, run: getExamUsedPercent } = useRequest(
      percentAPI,
      {
        formatResult: (res): number => {
          return <number>res.data.Data
        },
        manual: true,
      }
    )
    const examRestPercent = computed((): number => {
      return 100 - examUsedPercent.value
    })

    // 当前管理界面的管理类型
    return {
      listData,
      columns,
      pagination,
      currentLessonFilter,
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
      missionCheckpointVisiable,
      targetMissionID,
      openMissionCheckpointManager,

      // 课程筛选
      searchMissionListOptions,
      isMissionOpitonsLoading,
      missionListData,

      // 展开子行
      expandRowFn,
      examMissionsColumns,
      deleteExamMissionFn,

      // 实验修改表单
      missionModalVisible,
      missionModalTitle,
      missionModalStatus,
      missionFormRef,
      clearMissionForm,
      getMissionUpdateParams,
      isAddMissionLoading,
      isEditMissionLoading,
      deleteMission,
      addMissionFn,
      editMissionFn,
      deleteMissionFn,
      commitMissionModal,
      examRestPercent,
      ModalStatus,

      // 实验选项
      missionOptions,
      getMissionOptions,
      isMissionOptionsLoading,
    }
  },
})
</script>

<style lang="less" scoped>
</style>