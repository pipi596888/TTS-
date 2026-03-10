import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3000,
    host: '0.0.0.0',
    proxy: {
      '/api/user': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/feedback': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/admin/system': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/admin/users': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/admin/roles': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/admin/logs': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/admin/feedback': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/admin/voice': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      '/api/voice': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      '/api/tts': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/api/works': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    },
  },
})

