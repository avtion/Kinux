import { createRouter, createWebHashHistory } from 'vue-router'

// 路由
import loginComponents from '@/components/login.vue' // 用户登陆界面
import dashboardComponents from '@/components/dashboard.vue' // 操作界面
import workspaceComponents from '@/components/workSpace.vue' // 工作间统计
import shellComponents from '@/components/shell.vue' // 终端
import departmentManagerComponents from '@/components/departmentManager.vue' // 班级管理
import AccountManagerComponents from '@/components/accountManager.vue' // 用户管理
import deploymentManagerComponents from '@/components/deploymentManager.vue' // 配置管理
import checkpointManagerComponents from '@/components/checkpointManager.vue' // 检查点管理
import missionManagerComponents from '@/components/missionManager.vue' // 实验管理
import examManagerComponents from '@/components/examManager.vue' // 考试管理

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
        component: AccountManagerComponents,
        props: true,
      },
      {
        path: 'admin/mc',
        name: 'missionManager',
        component: missionManagerComponents,
        props: true,
      },
      {
        path: 'admin/ex',
        name: 'examManager',
        component: examManagerComponents,
        props: true,
      },
      {
        path: 'admin/deployment',
        name: 'deploymentManager',
        component: deploymentManagerComponents,
        props: true,
      },
      {
        path: 'admin/cp',
        name: 'checkpointManager',
        component: checkpointManagerComponents,
        props: true,
      },
    ],
  },
]

export default createRouter({
  history: createWebHashHistory(),
  routes,
})
