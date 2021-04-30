import { createRouter, createWebHashHistory } from 'vue-router'

// 路由
const loginComponents = () => import('@/components/login.vue') // 用户登陆界面
const dashboardComponents = () => import('@/components/dashboard.vue') // 操作界面
const workspaceComponents = () => import('@/components/workSpace.vue') // 工作间统计
const shellComponents = () => import('@/components/shell.vue') // 终端
const departmentManagerComponents = () => import('@/components/departmentManager.vue')// 班级管理
const AccountManagerComponents = () => import('@/components/accountManager.vue') // 用户管理
const deploymentManagerComponents = () => import('@/components/deploymentManager.vue')// 配置管理
const checkpointManagerComponents = () => import('@/components/checkpointManager.vue')// 检查点管理
const missionManagerComponents = () => import('@/components/missionManager.vue') // 实验管理
const examManagerComponents = () => import('@/components/examManager.vue')// 考试管理

const profileComponents = () => import('@/components/profile.vue') // 个人资料
const sessionManagerComponents = () => import('@/components/sessionManager.vue') // 实验会话
const lessonManagerComponents = () => import('@/components/lessonManager.vue')// 课程管理

// workspcae
const LessonSelector = () => import('@/components/workspaceComponents/lesson.vue')
const MissionSelector = () => import('@/components/workspaceComponents/mission.vue')
const ExamSelector = () => import('@/components/workspaceComponents/exam.vue')
const ExamMissionSelector = () => import('@/components/workspaceComponents/examMission.vue')

const shellWatcherComponents = () => import('@/components/shellWatcher.vue')// 终端监控

// 查询成绩
const stuScore = () => import('@/components/score/student.vue')
const teaScore = () => import('@/components/score/teacher.vue')
const ex = () => import('@/components/score/ex.vue')
const ms = () => import('@/components/score/ms.vue')

// 统计
const counter = () => import('@/components/counter.vue')

// 实验相关
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
      // 学生查询成绩
      {
        path: '/dashboard/score/stu',
        name: 'stuScore',
        component: stuScore,
      },
      // 老师查询成绩
      {
        path: '/dashboard/score/teacher',
        name: 'teaScore',
        component: teaScore,
      },
      {
        path: '/dashboard/score/ex',
        name: 'exScore',
        component: ex,
        props: true,
      },
      {
        path: '/dashboard/score/ms',
        name: 'msScore',
        component: ms,
        props: true,
      },
      {
        path: 'admin/counter',
        name: 'counterManager',
        component: counter,
        props: true,
      },
    ],
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
