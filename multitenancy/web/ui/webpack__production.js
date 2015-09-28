var path = require('path')
var webpack = require('webpack')
var autoprefixer = require('autoprefixer')
var HtmlWebpackPlugin = require('html-webpack-plugin')
var ExtractTextPlugin = require("extract-text-webpack-plugin")

var exclude = /node_modules/

function getLoader(test, exclude, loader) {
  return {test: test, exclude: exclude, loader: loader}
}

module.exports = {
  entry: './src/index.js',
  output: {
    path: __dirname + '/public/',
    filename: 'js/bundle.js',
    publicPath: '/'
  },
  plugins: [
    // new webpack.optimize.UglifyJsPlugin(),
    // new webpack.optimize.OccurenceOrderPlugin(),
    // new webpack.optimize.DedupePlugin(),
    new ExtractTextPlugin("css/style.css"),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      favicon: 'assets/favicon.ico',
      inject: true
    }),
    new webpack.DefinePlugin({
      __DEV__: JSON.stringify(JSON.parse('false'))
    })
  ],
  resolve: {
    root: __dirname + '/src/',
    extensions: ['', '.js', '.jsx']
  },
  module: {
    loaders: [
      getLoader(/\.jsx?$/, exclude, 'babel'),
      getLoader(
        /\.scss$/,
        exclude,
        ExtractTextPlugin.extract(
          'style',
          'css!postcss!sass'
        )
      )
    ]
  },
  postcss: {
    defaults: [autoprefixer]
  }
}


