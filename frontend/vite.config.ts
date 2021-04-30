import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import styleImport from 'vite-plugin-style-import';

// https://vitejs.dev/config/
export default defineConfig({
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@api': resolve(__dirname, './src/apis'),
      '@image': resolve(__dirname, './src/assets/images'),
    },
  },
  plugins: [
    vue(),
    styleImport({
      libs: [
        {
          libraryName: 'ant-design-vue',
          esModule: true,
          resolveStyle: (name) => {
            return `ant-design-vue/es/${name}/style/index`;
          },
        },
      ],
    }),
  ],
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:9001/',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      '/ws/': {
        target: 'ws://127.0.0.1:9001/v1/ws/',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/ws/, ''),
        ws: true
      },
    }
  }
})
