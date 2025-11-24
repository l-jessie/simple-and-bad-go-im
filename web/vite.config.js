import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      // Proxy /v1/api requests to our Go backend
      '/v1/api': {
        target: 'http://localhost:8070',
        changeOrigin: true,
        // WebSocket proxying must be explicitly enabled.
        ws: true,
      },
    },
  },
})
