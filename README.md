# 基于Kubernetes的Linux实验考试平台

## 1. 项目介绍

`Kinux` 是基于 `Kubernetes` 的Linux实验考试平台，前端采用 `Vue3` 编写，后端采用 `Golang` 编写，项目作为本科毕业设计用于毕业论文的编写，仅供参考。

## 2. 技术选型

🐋 Kubernetes： `K3s v1.19.4+k3s1`，使用`containerd`实现容器虚拟化

💻 前端：`Vite2 + Vue3.0 + Antd2`，包括`xterm.js`实现页面终端

🧠 后端： `Golang`，使用`Gin`作为Web框架

📡 前后端交互：前端使用`axios.js`发起Ajax请求，后端基于`gorilla/websocket`的实现HTTP双工通信

☢️ 鉴权：基于`Casbin`结合`Json Web Token`  的`RBAC`鉴权方式

🧫 数据库：基于`GORM V2`支持`SQLite`和`MySQL`两种关系型数据库

🏗️ 富文本渲染：`v-md-editor v2.2.1`

## 3. 项目特性

✅ 用户登陆  
✅ 导航栏路由  
✅ 实验项目  
✅ 实验终端  
✅ 富文本支持  
✅ 个人资料  
❌ 页面传输文件至容器  
✅ 考试项目   
✅ 随机头像生成  
✅ 修改密码  
✅ 班级管理  
✅ 用户管理   
✅ 实验管理  
✅ 考试管理  
✅ 数据展示  
✅ 实验会话  
❌ 在线镜像Image编辑  
❓ 后台限时容器回收   
✅ Websocket双向通信流

## 4. 开发者

- [Avtion](https://github.com/avtion): 📧 [mail@avtion.cn](mailto:mail@avtion.cn)

## 5. 开源

请遵守`Apache2.0`协议并保留作者技术支持声明，谢谢。