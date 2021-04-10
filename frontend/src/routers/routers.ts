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
import lessonManagerComponents from '@/components/lessonManager.vue' // 课程管理

// workspcae
import LessonSelector from '@/components/workspaceComponents/lesson.vue'
import MissionSelector from '@/components/workspaceComponents/mission.vue'
import ExamSelector from '@/components/workspaceComponents/exam.vue'
import ExamMissionSelector from '@/components/workspaceComponents/examMission.vue'

import shellWatcherComponents from '@/components/shellWatcher.vue' // 终端监控

const workspaceChild = [
  {
    path: '/dashboard/lesson',
    name: 'lessonSelector',
    component: LessonSelector,
  },
  {
    path: '/dashboard/mission/:lesson',
    name: 'missionSelector',
    component: MissionSelector,
  },
  {
    path: '/dashboard/exam',
    name: 'examSelector',
    component: ExamSelector,
  },
  {
    path: '/dashboard/exam/mission/:exam',
    name: 'examMissionSelector',
    component: ExamMissionSelector,
  },
]

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
        redirect: { name: 'lessonSelector' },
        children: workspaceChild,
      },
      {
        path: 'shell',
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
      {
        path: 'admin/lesson',
        name: 'lessonManager',
        component: lessonManagerComponents,
        props: true,
      },
      {
        path: 'admin/watcher',
        name: 'shellWatcher',
        component: shellWatcherComponents,
        props: true,
      },
    ],
  },
]

export default createRouter({
  history: createWebHashHistory(),
  routes,
})
