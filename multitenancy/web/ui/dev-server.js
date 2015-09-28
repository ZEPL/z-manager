import config           from "./webpack__development.js"
import webpack          from "webpack"
import WebpackDevServer from "webpack-dev-server"

const devPort = "9090"
const devHost = "localhost"
const devUrl  = `http://${devHost}:${devPort}/`
const builder = webpack(config(devUrl))

new WebpackDevServer(builder, {
  publicPath: devUrl,
  quiet: false,
  noInfo: false,
  hot: true,
  historyApiFallback: true,
  stats: {
    assets: false,
    colors: true,
    version: false,
    hash: false,
    timings: false,
    chunks: false,
    chunkModules: false
  }
}).listen(devPort, devHost, (err, result) => {
  if (err) {
    console.log(err)
  }
  
  console.log("Webpack server is listening at " + devUrl)
})

