import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { chunkSplitPlugin } from 'vite-plugin-chunk-split'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: true,
    port: 8002,
    watch: {
      usePolling: false,
    }
  },
  build: {
    chunkSizeWarningLimit: 1600,
    minify: "esbuild"
  },
  plugins: [
    svelte(), chunkSplitPlugin()
  ]
})
