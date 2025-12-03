/// <reference types="vitest" />
import vue from "@vitejs/plugin-vue";
import path from "path";
import { fileURLToPath } from "url";
import { defineConfig } from "vite";
import vuetify from "vite-plugin-vuetify";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), vuetify({ autoImport: true })],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src")
    },
    extensions: [".mjs", ".js", ".ts", ".jsx", ".tsx", ".json", ".vue"]
  },
  css: {
    preprocessorOptions: {
      scss: {
        api: "modern-compiler",
        silenceDeprecations: ["legacy-js-api"],
        additionalData: `@use "@/styles/_variables.scss" as *;`
      }
    }
  },
  server: {
    host: true,
    port: 8080,
    proxy: {
      "/api": {
        target: "http://localhost:1323",
        changeOrigin: true
      }
    }
  },
  // @ts-expect-error: test option is not in UserConfigExport
  test: {
    globals: true,
    environment: "happy-dom",
    setupFiles: ["./tests/setup.ts"],
    server: {
      deps: {
        inline: ["vuetify"]
      }
    },
    exclude: [
      "**/node_modules/**",
      "**/dist/**",
      "**/cypress/**",
      "**/.{idea,git,cache,output,temp}/**",
      "**/{karma,rollup,webpack,vite,vitest,jest,ava,babel,nyc,cypress,tsup,build}.config.*",
      "tests/e2e/**"
    ]
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vuetify: ["vuetify"],
          vendor: ["vue", "vue-router", "pinia"]
        }
      }
    }
  }
});
