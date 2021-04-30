<template>
  <div class="w-full h-full">
    <div class="w-full h-full p-5">
      <div class="w-full bg-white rounded">
        <!-- 头部展示栏 -->
        <a-page-header
          :ghost="false"
          style="border: 1px solid rgb(235, 237, 240)"
          title="系统状态统计"
          sub-title=""
        >
          <!-- 主内容 -->
          <a-row>
            <a-space :size="30">
              <!-- 头像 -->
              <div class="avatar">
                <a-avatar size="large" :src="avatar" />
              </div>

              <!-- 用户名 -->
              <a-statistic
                title="用户名"
                :value="
                  profile.realName == '' ? profile.username : profile.realName
                "
              >
              </a-statistic>
              <a-divider type="vertical" />
              <!-- 用户类型 -->
              <a-statistic title="用户类型" :value="profile.role">
              </a-statistic>
            </a-space>
          </a-row>
        </a-page-header>
      </div>
      <div class="pt-5">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="用户总数" :value="data.account" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="活跃用户" :value="data.session" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="班级数量" :value="data.department" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="课程数量" :value="data.lesson" />
            </a-card>
          </a-col>
        </a-row>
      </div>
      <div class="pt-5">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="实验数量" :value="data.mission" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="考试数量" :value="data.exam" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="考点数量" :value="data.checkpoint" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card :bordered="false" :hoverable="true">
              <a-statistic title="容器配置" :value="data.deployment" />
            </a-card>
          </a-col>
        </a-row>
      </div>
    </div>
  </div>
</template>

<script lang="ts"  type="module">
import { defineComponent, ref, watch, reactive } from 'vue'
// store
import { GetStore } from '@/store/store'

// vue-router
import { Profile } from '@/store/interfaces'

// 图标生成
import { ProfileAvatarCreator } from '@/utils/avatar'

import { BaseResponse, defaultClient } from '@/apis/request'
import { AxiosResponse } from 'axios'

interface counterRespData {
  account: number
  department: number
  deployment: number
  lesson: number
  mission: number
  exam: number
  checkpoint: number
  session: number
}

export default {
  setup(props) {
    // vue相关变量
    const store = GetStore()

    // 用户资料
    const profile = <Profile>store.getters.GetProfile

    // 头像
    const avatar = ProfileAvatarCreator(<string>store.getters.GetAvatarSeed)

    // 数据
    const data = ref<counterRespData>(<counterRespData>{})
    defaultClient
      .get('/v2/counter/')
      .then((res: AxiosResponse<BaseResponse>) => {
        const _res: counterRespData = res.data.Data
        data.value = _res
      })

    return {
      profile,
      avatar,

      data,
    }
  },
}
</script>

<style>
</style>