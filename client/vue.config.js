module.exports = {
  transpileDependencies: ["vuetify"],
  devServer: {
    proxy: {
      "/api": {
        target: "http://localhost:1323/", // local api server
        changeOrigin: true
      }
    }
  }
};
