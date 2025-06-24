const { defineConfig } = require("@vue/cli-service");
const https = require("https");

const keepAliveAgent = new https.Agent({ keepAlive: true });

module.exports = defineConfig({
  css: {
    loaderOptions: {
      scss: {
        additionalData: '@import "~@/styles/index.scss";'
      }
    }
  },
  transpileDependencies: ["vuetify"],
  devServer: {
    proxy: {
      "/api": {
        target: "https://jomon-dev.tokyotech.org",
        changeOrigin: true,
        agent: keepAliveAgent
      }
    }
  },
  productionSourceMap: false
});
