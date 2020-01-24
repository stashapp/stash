// development config
require('dotenv').config();
const merge = require('webpack-merge');
const commonConfig = require('./webpack.common');

module.exports = merge(commonConfig, {
  mode: 'development',
  entry: [
    './src/index.tsx' // the entry point of our app
  ],
  output: {
    filename: 'static/js/bundle.js',
    chunkFilename: 'static/js/[name].chunk.js',
  },
  optimization: {
      // Automatically split vendor and commons
      // https://twitter.com/wSokra/status/969633336732905474
      // https://medium.com/webpack/webpack-4-code-splitting-chunk-graph-and-the-splitchunks-optimization-be739a861366
      splitChunks: {
        chunks: 'all',
        name: false,
      },
      // Keep the runtime chunk separated to enable long term caching
      // https://twitter.com/wSokra/status/969679223278505985
      // https://github.com/facebook/create-react-app/issues/5358
      runtimeChunk: {
        name: entrypoint => `runtime-${entrypoint.name}`,
      },
  },
  devServer: {
    compress: true,
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
