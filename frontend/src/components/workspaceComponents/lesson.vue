<template>
  <a-card title="课程选择" :bordered="false" :loading="isProjectDataLoading">
    <a-list :data-source="data">
      <template #renderItem="{ item }">
        <a-list-item>
          <a
            class="container mx-auto bg-gray-50 h-auto w-full p-4 space-y-1 rounded-lg hover:bg-gray-200"
          >
            <div class="block text-2xl font-extrabold font-sans relative">
              <div class="inline-block text-blue-400">|</div>
              <div class="inline-block ml-2 text-gray-700">
                {{ item.title }}
              </div>
              <div class="inline-block absolute inset-y-0 right-0">></div>
            </div>
            <div
              class="block text-sm pr-8 text-opacity-5 font-mono leading-relaxed h-auto line-clamp-3"
            >
              {{ item.content }}
            </div>
          </a>
        </a-list-item>
      </template>
    </a-list>
  </a-card>
</template>

<script lang="ts" type="module">
// antd
import { RightOutlined } from '@ant-design/icons-vue'

interface DataItem {
  title: string
  content: string
}

const bodyStyle = {
  'line-height': '22px',
  height: '44px',
  overflow: 'hidden',
  '-webkit-line-clamp': '1',
  ' text-overflow': 'ellipsis',
  display: '-webkit-box',
  '-webkit-box-orient': 'vertical',
}
const data: DataItem[] = [
  {
    title: '计算机网络',
    content:
      '计算机网络是指将地理位置不同的具有独立功能的多台计算机及其外部设备，通过通信线路连接起来，在网络操作系统，网络管理软件及网络通信协议的管理和协调下，实现资源共享和信息传递的计算机系统。',
  },
  {
    title: 'Title 2',
    content: '',
  },
  {
    title: 'Title 3',
    content: '',
  },
  {
    title: 'Title 4',
    content: '',
  },
  {
    title: 'Title 5',
    content: '',
  },
  {
    title: 'Title 6',
    content: '',
  },
]

export default {
  components: {
    RightOutlined,
  },
  setup(props, ctx) {
    return {
      isProjectDataLoading: false,
      data,
      bodyStyle,
    }
  },
}
</script>

<style>
</style>