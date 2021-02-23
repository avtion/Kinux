import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@api': resolve(__dirname, './src/apis'),
      '@image': resolve(__dirname, './src/assets/images'),
    },
  },
  plugins: [vue()],
})
