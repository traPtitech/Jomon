const https = require("https");
const path = require("path");
const { DEV_SERVER_PROXY_HOST } = require("./dev.config");

const keepAliveAgent = new https.Agent({ keepAlive: true });

module.exports = {
  configureWebpack: {
    resolve: {
      alias: {
        "/@": path.resolve(__dirname, "src").replace(/\\/g, "/")
      },
      extensions: [".js", ".vue", ".json"]
    }
  },
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
