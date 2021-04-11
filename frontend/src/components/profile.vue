<template>
  <div>
    <a-layout style="padding: 24px 24px 24px">
      <a-layout-content
        :style="{
          background: '#fff',
          minHeight: '600px',
        }"
      >
        <a-card :bordered="false">
          <a-descriptions
            title="用户信息"
            layout="horizontal"
            bordered
            :column="3"
          >
            <!-- 头像 -->
            <a-descriptions-item label="用户名" :span="2">{{
              profile.username
            }}</a-descriptions-item>
            <a-descriptions-item label="头像" :span="1">
              <a-avatar :size="64" shape="square" :src="avatar"> </a-avatar>
              <a-button
                size="small"
                :style="{ marginLeft: '12px', verticalAlign: 'middle' }"
                @click="updateAvatar"
              >
                随机生成
              </a-button>
            </a-descriptions-item>

            <a-descriptions-item label="真实姓名">{{
              profile.realName
            }}</a-descriptions-item>

            <a-descriptions-item label="班级">{{
              profile.department
            }}</a-descriptions-item>

            <a-descriptions-item label="账号权限">{{
              profile.role
            }}</a-descriptions-item>

            <a-descriptions-item label="账号状态" :span="3"
              ><a-badge status="success" text="可用"
            /></a-descriptions-item>

            <!-- 修改密码表单 -->
            <a-descriptions-item label="修改密码" :span="3">
              <a-form>
                <a-form-item
                  label="原密码"
                  has-feedback
                  v-bind="validateInfos.oldValue"
                >
                  <a-input
                    type="password"
                    autocomplete="off"
                    v-model:value="pwForm.oldValue"
                    @blur="
                      validate('oldValue', { trigger: 'blur' }).catch(() => {})
                    "
                    :disabled="commitButtonLoading"
                  />
                </a-form-item>
                <a-form-item
                  label="新密码"
                  has-feedback
                  v-bind="validateInfos.newValue"
                >
                  <a-input
                    type="password"
                    autocomplete="off"
                    v-model:value="pwForm.newValue"
                    @blur="
                      validate('newValue', { trigger: 'blur' }).catch(() => {})
                    "
                    :disabled="commitButtonLoading"
                  />
                </a-form-item>
                <a-form-item
                  label="重复新密码"
                  has-feedback
                  v-bind="validateInfos.repeatValue"
                >
                  <a-input
                    type="password"
                    autocomplete="off"
                    v-model:value="pwForm.repeatValue"
                    @blur="
                      validate('repeatValue', {
                        trigger: 'blur',
                      }).catch(() => {})
                    "
                    :disabled="commitButtonLoading"
                  />
                </a-form-item>
                <a-form-item>
                  <a-button
                    type="primary"
                    @click="commitPwForm"
                    :loading="commitButtonLoading"
                    >更新密码</a-button
                  >
                </a-form-item>
              </a-form>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script lang="ts" type="module">
// vue
import { defineComponent, reactive, ref, watch, UnwrapRef } from 'vue'

// antd
import { RuleObject } from 'ant-design-vue/es/form/interface'
import { useForm } from '@ant-design-vue/use'

// 图标生成
import Avatars from '@dicebear/avatars'
import AvatarsSprites from '@dicebear/avatars-avataaars-sprites'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'
import { Profile } from '@/store/interfaces'
import { Account } from '@/apis/user'

export default defineComponent({
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const routers = useRouter()

    // 用户资料
    const profile = <Profile>store.getters.GetProfile

    // 头像
    const avatar = ref<string>(
      new Avatars(AvatarsSprites, {
        dataUri: true,
      }).create(<string>store.getters.GetAvatarSeed)
    )
    watch(
      () => <string>store.getters.GetAvatarSeed,
      (newValue) => {
        // 头像种子更新
        avatar.value = new Avatars(AvatarsSprites, {
          dataUri: true,
        }).create(newValue)
      }
    )
    // 更新用户头像种子
    const updateAvatar = () => {
      new Account().updateAvatarSeed()
    }

    // 密码重置表单
    const pwForm: UnwrapRef<pwFormFields> = reactive({
      oldValue: '',
      newValue: '',
      repeatValue: '',
    })
    const repeatPwValidator = async (rule: RuleObject, value: string) => {
      if (value === '') {
        return Promise.reject('请重新输入新密码')
      } else if (value !== pwForm.newValue) {
        return Promise.reject('新密码与重新输入的密码不一样')
      } else if (value === pwForm.oldValue) {
        return Promise.reject('新密码和旧密码一样')
      } else {
        return Promise.resolve()
      }
    }
    const pwFormRule = reactive({
      oldValue: [
        { required: true, message: '请输入当前使用的密码', trigger: 'blur' },
      ],
      newValue: [
        { required: true, message: '请输入要设置的新密码', trigger: 'blur' },
        {
          min: 4,
          message: '新密码最少需要4位',
          trigger: 'blur',
        },
      ],
      repeatValue: [{ validator: repeatPwValidator, trigger: 'change' }],
    })
    const { resetFields, validate, validateInfos } = useForm(pwForm, pwFormRule)
    const commitButtonLoading = ref<boolean>(false)
    const commitPwForm = () => {
      validate().then(() => {
        commitButtonLoading.value = true
        new Account()
          .updatePassword(pwForm.oldValue, pwForm.newValue)
          .finally(() => {
            commitButtonLoading.value = false
          })
      })
    }
    return {
      profile,
      avatar,
      updateAvatar,
      pwForm,
      commitPwForm,
      pwFormRule,
      commitButtonLoading,
      validateInfos,
      validate,
    }
  },
})

// 重置密码表单
interface pwFormFields {
  oldValue: string
  newValue: string
  repeatValue: string
}
</script>

<style lang="less" scoped>
</style>