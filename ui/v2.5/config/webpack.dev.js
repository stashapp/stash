// development config
require('dotenv').config();
const merge = require('webpack-merge');
const commonConfig = require('./webpack.common');

module.exports = merge(commonConfig, {
  mode: 'development',
  entry: [
    './src/index.tsx' // the entry point of our app
  ],
  devServer: {
    host: '0.0.0.0',
    hot: true, // enable HMR on the server host: '0.0.0.0',
    port: process.env.PORT,
    historyApiFallback: true,
    stats: {
        assets: true,
        builtAt: true,
        modules: false,
        children: false
    }
  },
  devtool: 'eval-source-map',
});
