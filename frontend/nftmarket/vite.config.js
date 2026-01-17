import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vitejs.dev/config/
export default defineConfig({
  // 🚀 关键修复点：设置公共基础路径为相对路径
  // 使用 './' 确保 index.html 引用资源时是相对于自身目录的，
  // 而不是绝对根路径 /src/main.js
  base: './',

  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
})