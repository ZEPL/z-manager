var path = require('path');
var webpack = require('webpack');
var autoprefixer = require('autoprefixer');
var HtmlWebpackPlugin = require('html-webpack-plugin')

var exclude = /node_modules/;

function getLoader(test, exclude, loader) {
  return {test: test, exclude: exclude, loader: loader}
}

module.exports = function(publicPath) {return {
  devtool: 'eval',
  entry: [
    'webpack-dev-server/client?' + publicPath,
    'webpack/hot/only-dev-server',
    './src/index.js'
  ],
  output: {
    path: '/',
    publicPath: publicPath
  },
  plugins: [
    new webpack.NoErrorsPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      favicon: 'assets/favicon.ico',
      inject: 'body'
    }),
    new webpack.DefinePlugin({
      __DEV__: JSON.stringify(JSON.parse('true'))
    })
  ],
  resolve: {
    root: __dirname + '/src/',
    extensions: ['', '.js', '.jsx']
  },
  module: {
    loaders: [
      getLoader(/\.jsx?$/, exclude, 'react-hot!babel'),
      getLoader(
        /\.scss$/,
        exclude,
        'style!css!postcss!sass'
      ),
    ]
  },
  postcss: {
    defaults: [autoprefixer]
  }
}};


