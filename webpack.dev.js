const config = require("./webpack.base.js");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

config.mode = "development";
config.output.filename = pathData => {
  return pathData.chunk.name.includes("serviceWorker") ? "js/serviceWorker.js" : "js/[name].development.js";
};

config.plugins.push(
  new MiniCssExtractPlugin({
    filename: "css/[name].development.css",
    chunkFilename: "css/[id].development.css"
  })
);

module.exports = config;
