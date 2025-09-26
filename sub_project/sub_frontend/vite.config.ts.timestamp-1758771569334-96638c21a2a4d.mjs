// vite.config.ts
import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "file:///home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/node_modules/vite/dist/node/index.js";
import fs from "node:fs";
import path from "node:path";
import vue from "file:///home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import Components from "file:///home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/node_modules/unplugin-vue-components/dist/vite.js";
import AutoImport from "file:///home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/node_modules/unplugin-auto-import/dist/vite.js";
import { ElementPlusResolver } from "file:///home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/node_modules/unplugin-vue-components/dist/resolvers.js";
var __vite_injected_original_dirname = "/home/ubuntu/zlb/project-kt3/sub_project/sub_frontend";
var __vite_injected_original_import_meta_url = "file:///home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/vite.config.ts";
var vite_config_default = defineConfig({
  // Deployed under Nginx at /topic-three/online-experiment/
  base: "/topic-three/online-experiment/",
  plugins: [
    vue(),
    vueJsx(),
    AutoImport({
      resolvers: [ElementPlusResolver()]
    }),
    Components({
      resolvers: [ElementPlusResolver()]
    })
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url))
    }
  },
  server: {
    host: true,
    https: {
      key: fs.readFileSync(path.resolve(__vite_injected_original_dirname, "../../certs/server.key")),
      cert: fs.readFileSync(path.resolve(__vite_injected_original_dirname, "../../certs/server.crt"))
    },
    proxy: {
      // 配置代理（只在本地开发有效，上线无效）
      "/web/keti3": {
        target: "https://test-api.edstars.com.cn/",
        changeOrigin: true
      }
    }
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvaG9tZS91YnVudHUvemxiL3Byb2plY3Qta3QzL3N1Yl9wcm9qZWN0L3N1Yl9mcm9udGVuZFwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiL2hvbWUvdWJ1bnR1L3psYi9wcm9qZWN0LWt0My9zdWJfcHJvamVjdC9zdWJfZnJvbnRlbmQvdml0ZS5jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL2hvbWUvdWJ1bnR1L3psYi9wcm9qZWN0LWt0My9zdWJfcHJvamVjdC9zdWJfZnJvbnRlbmQvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgeyBmaWxlVVJMVG9QYXRoLCBVUkwgfSBmcm9tICdub2RlOnVybCdcblxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSdcbmltcG9ydCBmcyBmcm9tICdub2RlOmZzJ1xuaW1wb3J0IHBhdGggZnJvbSAnbm9kZTpwYXRoJ1xuaW1wb3J0IHZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnXG5pbXBvcnQgdnVlSnN4IGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZS1qc3gnXG5pbXBvcnQgQ29tcG9uZW50cyBmcm9tICd1bnBsdWdpbi12dWUtY29tcG9uZW50cy92aXRlJ1xuaW1wb3J0IEF1dG9JbXBvcnQgZnJvbSAndW5wbHVnaW4tYXV0by1pbXBvcnQvdml0ZSdcbmltcG9ydCB7IEVsZW1lbnRQbHVzUmVzb2x2ZXIgfSBmcm9tICd1bnBsdWdpbi12dWUtY29tcG9uZW50cy9yZXNvbHZlcnMnXG5cbi8vIGh0dHBzOi8vdml0ZWpzLmRldi9jb25maWcvXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoe1xuICAvLyBEZXBsb3llZCB1bmRlciBOZ2lueCBhdCAvdG9waWMtdGhyZWUvb25saW5lLWV4cGVyaW1lbnQvXG4gIGJhc2U6ICcvdG9waWMtdGhyZWUvb25saW5lLWV4cGVyaW1lbnQvJyxcbiAgcGx1Z2luczogW3Z1ZSgpLCB2dWVKc3goKSxcbiAgICBBdXRvSW1wb3J0KHtcbiAgICAgIHJlc29sdmVyczogW0VsZW1lbnRQbHVzUmVzb2x2ZXIoKV0sXG4gICAgfSksXG4gICAgQ29tcG9uZW50cyh7XG4gICAgICByZXNvbHZlcnM6IFtFbGVtZW50UGx1c1Jlc29sdmVyKCldLFxuICAgIH0pLFxuICBdLFxuICByZXNvbHZlOiB7XG4gICAgYWxpYXM6IHtcbiAgICAgICdAJzogZmlsZVVSTFRvUGF0aChuZXcgVVJMKCcuL3NyYycsIGltcG9ydC5tZXRhLnVybCkpXG4gICAgfVxuICB9LFxuICBzZXJ2ZXI6IHtcbiAgICBob3N0OiB0cnVlLFxuICAgIGh0dHBzOiB7XG4gICAgICBrZXk6IGZzLnJlYWRGaWxlU3luYyhwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnLi4vLi4vY2VydHMvc2VydmVyLmtleScpKSxcbiAgICAgIGNlcnQ6IGZzLnJlYWRGaWxlU3luYyhwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnLi4vLi4vY2VydHMvc2VydmVyLmNydCcpKSxcbiAgICB9LFxuICAgIHByb3h5OiB7XG4gICAgICAvLyBcdTkxNERcdTdGNkVcdTRFRTNcdTc0MDZcdUZGMDhcdTUzRUFcdTU3MjhcdTY3MkNcdTU3MzBcdTVGMDBcdTUzRDFcdTY3MDlcdTY1NDhcdUZGMENcdTRFMEFcdTdFQkZcdTY1RTBcdTY1NDhcdUZGMDlcbiAgICAgICcvd2ViL2tldGkzJzoge1xuICAgICAgICB0YXJnZXQ6ICdodHRwczovL3Rlc3QtYXBpLmVkc3RhcnMuY29tLmNuLycsXG4gICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZVxuICAgICAgfVxuICAgIH1cbiAgfVxufSlcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBaVYsU0FBUyxlQUFlLFdBQVc7QUFFcFgsU0FBUyxvQkFBb0I7QUFDN0IsT0FBTyxRQUFRO0FBQ2YsT0FBTyxVQUFVO0FBQ2pCLE9BQU8sU0FBUztBQUNoQixPQUFPLFlBQVk7QUFDbkIsT0FBTyxnQkFBZ0I7QUFDdkIsT0FBTyxnQkFBZ0I7QUFDdkIsU0FBUywyQkFBMkI7QUFUcEMsSUFBTSxtQ0FBbUM7QUFBeUssSUFBTSwyQ0FBMkM7QUFZblEsSUFBTyxzQkFBUSxhQUFhO0FBQUE7QUFBQSxFQUUxQixNQUFNO0FBQUEsRUFDTixTQUFTO0FBQUEsSUFBQyxJQUFJO0FBQUEsSUFBRyxPQUFPO0FBQUEsSUFDdEIsV0FBVztBQUFBLE1BQ1QsV0FBVyxDQUFDLG9CQUFvQixDQUFDO0FBQUEsSUFDbkMsQ0FBQztBQUFBLElBQ0QsV0FBVztBQUFBLE1BQ1QsV0FBVyxDQUFDLG9CQUFvQixDQUFDO0FBQUEsSUFDbkMsQ0FBQztBQUFBLEVBQ0g7QUFBQSxFQUNBLFNBQVM7QUFBQSxJQUNQLE9BQU87QUFBQSxNQUNMLEtBQUssY0FBYyxJQUFJLElBQUksU0FBUyx3Q0FBZSxDQUFDO0FBQUEsSUFDdEQ7QUFBQSxFQUNGO0FBQUEsRUFDQSxRQUFRO0FBQUEsSUFDTixNQUFNO0FBQUEsSUFDTixPQUFPO0FBQUEsTUFDTCxLQUFLLEdBQUcsYUFBYSxLQUFLLFFBQVEsa0NBQVcsd0JBQXdCLENBQUM7QUFBQSxNQUN0RSxNQUFNLEdBQUcsYUFBYSxLQUFLLFFBQVEsa0NBQVcsd0JBQXdCLENBQUM7QUFBQSxJQUN6RTtBQUFBLElBQ0EsT0FBTztBQUFBO0FBQUEsTUFFTCxjQUFjO0FBQUEsUUFDWixRQUFRO0FBQUEsUUFDUixjQUFjO0FBQUEsTUFDaEI7QUFBQSxJQUNGO0FBQUEsRUFDRjtBQUNGLENBQUM7IiwKICAibmFtZXMiOiBbXQp9Cg==
