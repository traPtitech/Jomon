const https = require("https");
const keepAliveAgent = new https.Agent({ keepAlive: true });
const { DEV_SERVER_PROXY_HOST } = require("./dev.config");

module.exports = {
  transpileDependencies: ["vuetify"],
  devServer: {
    proxy: {
      "/api": {
        target: DEV_SERVER_PROXY_HOST,
        changeOrigin: true,
        agent: keepAliveAgent
      }
    }
  },
  productionSourceMap: false
};
