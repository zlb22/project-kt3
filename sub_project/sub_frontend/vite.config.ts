import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import fs from 'node:fs'
import path from 'node:path'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import Components from 'unplugin-vue-components/vite'
import AutoImport from 'unplugin-auto-import/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig({
  // Deployed under Nginx at /topic-three/
  base: '/topic-three/',
  plugins: [vue(), vueJsx(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: true,
    https: {
      key: fs.readFileSync(path.resolve(__dirname, '../../certs/server.key')),
      cert: fs.readFileSync(path.resolve(__dirname, '../../certs/server.crt')),
    },
    proxy: {
      // 配置代理（只在本地开发有效，上线无效）
      '/web/keti3': {
        target: 'https://test-api.edstars.com.cn/',
        changeOrigin: true
      }
    }
  }
})
