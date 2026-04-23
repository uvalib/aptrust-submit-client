import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
   define: {
      // enable hydration mismatch details in production build
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'true'
   },
   plugins: [
      vue(),
   ],
   resolve: {
      alias: {
         '@': fileURLToPath(new URL('./src', import.meta.url))
      },
   },
   server: { // this is used in dev mode only
      port: 8080,
      proxy: {
         '/api': {
            target: process.env.APTSUBMIT_CONFIG, // export APTSUBMIT_CONFIG=http://localhost:8085
            changeOrigin: true,
            logLevel: 'debug'
         },
         '/healthcheck': {
            target: process.env.APTSUBMIT_CONFIG,
            changeOrigin: true,
            logLevel: 'debug'
         },
         '/version': {
            target: process.env.APTSUBMIT_CONFIG,
            changeOrigin: true,
            logLevel: 'debug'
         },
         '/authenticate': {
            target: process.env.APTSUBMIT_CONFIG,
            changeOrigin: true,
            logLevel: 'debug'
         },
         '/config': {
            target: process.env.APTSUBMIT_CONFIG,
            changeOrigin: true,
            logLevel: 'debug'
         },
      },
   },
   css: {
      preprocessorOptions: {
         scss: {
            api: "modern-compiler",
            additionalData: `@use "@/assets/theme/colors.scss" as *;`
         },
      }
   },
})
