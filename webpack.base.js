const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
  devtool: "source-map",
  entry: {
    main: path.resolve(__dirname, "web/main.ts"),
    //serviceWorker: path.resolve(__dirname, "web/serviceWorker.ts")
  },
  output: {
    clean: true,
    path: path.resolve(__dirname, "bin/assets")
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        exclude: /node_modules/,
        use: ["babel-loader", "ts-loader"]
      },
      {
        test: /.less$/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: "css-loader",
            options: {
              url: false
            }
          },
          "postcss-loader",
          "less-loader"
        ]
      },
      {
        test: /.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: "css-loader",
            options: {
              url: false
            }
          },
          "postcss-loader"
        ]
      }
    ]
  },
  plugins: [
    new CopyPlugin({
      patterns: [
        {
          from: path.resolve(__dirname, "web/fonts"),
          to: path.resolve(__dirname, "bin/assets/fonts")
        }
      ]
    }),
    new HtmlWebpackPlugin({
      filename: "../templates/head.html",
      template: path.resolve(__dirname, "web/head.html"),
      chunks: ["main"],
      chunksSortMode: "manual",
      publicPath: "/"
    })
  ],
  resolve: {
    extensions: [".ts", ".tsx", ".js", ".jsx"]
  },
  watchOptions: {
    ignored: /node_modules/
  }
};
