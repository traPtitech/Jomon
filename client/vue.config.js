const https = require("https");

const keepAliveAgent = new https.Agent({ keepAlive: true });

module.exports = {
  css: {
    loaderOptions: {
      scss: {
        additionalData: '@import "src/styles/index.scss";'
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
};
