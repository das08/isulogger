const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    proxy: {
      '^/api': {
        target: 'http://isulogger-server:8082',
        pathRewrite: { "^/api/": "/" }
      },
    },
    allowedHosts: ["all"],
  }
})
