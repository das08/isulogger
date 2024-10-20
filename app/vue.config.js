const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    proxy: {
      '^/api': {
        target: 'http://localhost:8082',
        pathRewrite: { "^/api/": "/" }
      },
    },
    allowedHosts: ["all"],
  }
})
