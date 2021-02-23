import { createRouter, createWebHashHistory } from 'vue-router'

// 路由
import loginComponents from '@/components/login.vue' // 用户登陆界面
import dashboardComponents from '@/components/dashboard.vue' // 操作界面
import workspaceComponents from '@/components/workSpace.vue' // 工作间统计
import shellComponents from '@/components/shell.vue' // 终端
import managerComponents from '@/components/manager.vue' // 管理界面
import departmentManagerComponents from '@/components/departmentManager.vue' // 管理界面
import profileComponents from '@/components/profile.vue' // 个人资料
import sessionManagerComponents from '@/components/sessionManager.vue' // 实验会话

const routes = [
  {
    path: '/',
    name: 'login',
    component: loginComponents,
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: dashboardComponents,
    children: [
      {
        path: '',
        name: 'workspace',
        component: workspaceComponents,
      },
      {
        path: 'shell/:id',
        name: 'shell',
        component: shellComponents,
        props: true,
      },
      {
        path: 'profile',
        name: 'profile',
        component: profileComponents,
        props: true,
      },
      {
        path: 'session',
        name: 'session',
        component: sessionManagerComponents,
        props: true,
      },
      {
        path: 'admin/dp',
        name: 'departmentManager',
        component: departmentManagerComponents,
        props: true,
      },
      {
        path: 'admin/ac',
        name: 'AccountManager',
        component: managerComponents,
        props: true,
        meta: { managerType: 'ac' },
      },
      {
        path: 'admin/mc',
        name: 'missionManager',
        component: managerComponents,
        props: true,
        meta: { managerType: 'mc' },
      },
      {
        path: 'admin/ex',
        name: 'examManager',
        component: managerComponents,
        props: true,
        meta: { managerType: 'ex' },
      },
    ],
  },
]

export default createRouter({
  history: createWebHashHistory(),
  routes,
})
