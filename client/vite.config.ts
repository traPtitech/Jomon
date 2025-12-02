import vue from "@vitejs/plugin-vue";
import path from "path";
import { defineConfig } from "vite";
import envCompatible from "vite-plugin-env-compatible";
import vuetify from "vite-plugin-vuetify";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), vuetify({ autoImport: true }), envCompatible()],
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
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: ["./tests/setup.ts"],
    server: {
      deps: {
        inline: ["vuetify"]
      }
    }
  }
});
