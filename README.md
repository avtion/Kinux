# 基于Kubernetes的Linux实验考试平台

## 1. 项目介绍

`Kinux` 是基于 `Kubernetes` 的Linux实验考试平台，前端采用 `Vue3.0` 编写，后端采用 `Golang` 编写，项目作为本科毕业设计用于毕业论文编写，仍处于`WIP`(work in progress)阶段

## 2. 项目预览

## 3. 技术选型
🐋 Kubernetes： `K3s v1.19.4+k3s1`，使用`containerd`实现容器虚拟化  
💻 前端：`Vite2 + Vue3.0 + Antd2`，包括`xterm.js`实现页面终端  
🧠 后端： `Golang`，使用`Gin`作为Web框架    
📡 前后端交互：前端使用`axios.js`发起Ajax请求，后端基于`gorilla/websocket`的实现HTTP双工通信  
☢️ 鉴权：结合`Casbin`与`Json Web Token`  
🧫 数据库：基于`GORM V2`支持`SQLite`和`MySQL`两种关系型数据库  
🏗️ 富文本渲染：`v-md-editor v2.2.1`


## 4. 项目架构
⏳ TODO

## 5. 项目计划
⏳ TODO

## 6. 开发者
- [Avtion](https://github.com/avtion): 📧 mail@avtion.cn

## 7. 开源
请遵守`Apache2.0`协议并保留作者技术支持声明，谢谢。