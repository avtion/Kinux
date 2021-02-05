import { createRouter, createWebHashHistory } from "vue-router"

// 路由
import loginComponents from "@/components/login.vue" // 用户登陆界面
import dashboardComponents from "@/components/dashboard.vue" // 操作界面
import workspaceComponents from "@/components/workSpace.vue" // 工作间统计

const routes = [
  {
    path: "/",
    name: "login",
    component: loginComponents,
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: dashboardComponents,
    children: [
      {
        path: "",
        name: "workspace",
        component: workspaceComponents,
      },
    ],
  },
]

export default createRouter({
  history: createWebHashHistory(),
  routes,
})
