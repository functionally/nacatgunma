const path = require('path')

module.exports = {
  entry           : './src/controller.js'
, mode            : 'development'
, output          : {
    filename      : 'controller.js'
  , path          : path.resolve(__dirname, '.')
  , libraryTarget : "var"
  , library       : "Controller"
  },
  devServer: {
    static: ".",
    host: "127.0.0.1",
    port: 3000
  }
}
