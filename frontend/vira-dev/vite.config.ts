import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    host: '0.0.0.0', // нужно для docker-сервера
    allowedHosts: ['vira-dev.loc'],
    proxy: {
      '/api': 'http://gateway:8080', // обращаться к API через Docker-сервис
    },
  },
})