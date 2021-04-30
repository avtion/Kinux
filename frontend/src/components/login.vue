<template>
  <div class="bg-white font-family-karla h-screen">
    <div class="w-full flex flex-wrap">
      <div class="w-full md:w-1/2 flex flex-col">
        <div
          class="flex justify-center md:justify-start pt-12 md:pl-12 md:-mb-24"
        >
          <a href="#" class="bg-black text-white font-bold text-xl p-4">登陆</a>
        </div>

        <div
          class="flex flex-col justify-center md:justify-start my-auto pt-8 md:pt-0 px-8 md:px-24 lg:px-32"
        >
          <!-- title -->
          <p class="text-center text-3xl">Kinux 实验考试平台</p>

          <!-- 表单 -->
          <a-form
            class="flex flex-col pt-3 md:pt-8"
            onsubmit="event.preventDefault();"
            :model="loginFormData"
            :rules="loginFormRules"
          >
            <!-- 用户名 -->
            <div class="flex flex-col pt-4">
              <label for="username" class="text-lg">账号</label>
              <a-form-item
                class="form-item-custom"
                :wrapperCol="{ style: { width: '100%' } }"
                name="username"
                v-bind="validateInfos.username"
              >
                <a-input
                  placeholder="默认为学号"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mt-1 leading-tight"
                  v-model:value="loginFormData.username"
                  allow-clear
                >
                  <template #prefix
                    ><UserOutlined style="color: rgba(0, 0, 0, 0.25)"
                  /></template>
                </a-input>
              </a-form-item>
            </div>

            <!-- 密码 -->
            <div class="flex flex-col pt-4">
              <label for="password" class="text-lg">密码</label>
              <a-form-item
                class="form-item-custom"
                :wrapperCol="{ style: { width: '100%' } }"
                name="password"
                v-bind="validateInfos.password"
              >
                <a-input-password
                  placeholder="默认为姓名的拼音"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mt-1 leading-tight"
                  v-model:value="loginFormData.password"
                  allow-clear
                >
                  <template #prefix
                    ><LockOutlined style="color: rgba(0, 0, 0, 0.25)"
                  /></template>
                </a-input-password>
              </a-form-item>
            </div>

            <!-- 登陆按钮 -->
            <a-button
              type="primary"
              style="margin-top: 30px"
              size="large"
              :loading="isLoging"
              @click="onSubmit"
            >
              登陆
            </a-button>
          </a-form>
          <div class="text-center pt-12 pb-12">
            <p>
              v1.0.0 Power By Kubernetes
              <!-- <a href="register.html" class="underline font-semibold"
                >找回密码</a
              > -->
            </p>
            <p>© 2020-2021 广东第二师范学院 计算机科学系 马乃韬</p>
          </div>
        </div>
      </div>

      <!-- 图片展示 -->
      <div class="w-1/2 shadow-2xl">
        <img
          class="object-cover w-full h-screen hidden md:block"
          :src="backgroundImg"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" type="module">
import { reactive, ref, defineComponent, onBeforeMount, computed } from 'vue'

// antd
import { useForm } from '@ant-design-vue/use'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'

// api
import { Account } from '@api/user'

// 背景图
import backgroundImg from '@image/login-wallpaper.png'

// store
import { GetStore } from '@/store/store'

// vue-router
import { useRouter } from 'vue-router'

import { Profile, Role } from '@/store/interfaces'

// 账号
interface account {
  username: string
  password: string
}

export default defineComponent({
  name: 'login',
  components: { UserOutlined, LockOutlined },
  setup(props, ctx) {
    // vue相关变量
    const store = GetStore()
    const router = useRouter()

    // 页面载入时先判断JWT是否有效
    onBeforeMount(() => {
      const token = store.getters.GetJWTToken
      if ((token as string) && token != '') {
        console.log('存在JWT密钥，正在与服务端进行校验')
        // TODO 密钥校验
        notification.success({ message: '您已经成功登陆' })
        router.push('dashboard')
      }
    })

    // 登陆状态进行时
    const isLoging = ref(false)

    // 表单数据
    const loginFormData: account = reactive({
      username: '',
      password: '',
    })

    // 表单校验规则
    const loginFormRules = reactive({
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    })

    // 登陆表单useForm解构
    const { resetFields, validate, validateInfos } = useForm(
      loginFormData,
      loginFormRules
    )

    // 登陆按钮触发函数
    const onSubmit = (e: { preventDefault: () => void }) => {
      e.preventDefault()

      // 修改登陆状态
      validate().then(() => {
        isLoging.value = true
        // 登陆
        new Account(loginFormData.username, loginFormData.password)
          .login()
          .then(() => {
            notification.success({ message: '登陆成功' })

            // 获取当前用户角色并跳转至不同的页面
            const p: Profile = store.getters.GetProfile
            if (p.roleID == Role.RoleAdmin || p.roleID === Role.RoleManager) {
              router.push({ name: 'counterManager' })
            } else {
              router.push('dashboard')
            }
          })
          .finally(() => {
            isLoging.value = false
          })
      })
    }

    return {
      backgroundImg,
      isLoging,
      loginFormData,
      loginFormRules,
      onSubmit,
      validateInfos,
    }
  },
})
</script>

<style>
.form-item-custom {
  margin-bottom: 0px;
}
</style>
