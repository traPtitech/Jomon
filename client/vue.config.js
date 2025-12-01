const crypto = require("crypto");

// Workaround for OpenSSL 3 bug
const crypto_orig_createHash = crypto.createHash;
crypto.createHash = algorithm =>
  crypto_orig_createHash(algorithm == "md4" ? "sha256" : algorithm);

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
        target: "http://127.0.0.1:1323",
        changeOrigin: true
      }
    }
  },
  configureWebpack: {
    output: {
      hashFunction: "xxhash64"
    },
    cache: false
  },
  parallel: false,
  chainWebpack: config => {
    config.module
      .rule("js")
      .use("babel-loader")
      .tap(options => {
        return { ...options, cacheCompression: false };
      });
  },
  productionSourceMap: false
};
