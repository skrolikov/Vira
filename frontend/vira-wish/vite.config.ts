import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5174,
    host: '0.0.0.0', // нужно для docker-сервера
    allowedHosts: ['vira-wish.loc'],
    proxy: {
      '/api': 'http://gateway:8080', // обращаться к API через Docker-сервис
    },
  },
})
