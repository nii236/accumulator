const webpack = require("webpack")
const path = require("path")
const HtmlWebpackPlugin = require("html-webpack-plugin")

const config = {
	entry: "./src/index.tsx",
	output: {
		path: path.resolve(__dirname, "dist"),
		publicPath: "/",
		filename: "[name].[contenthash].js",
	},
	module: {
		rules: [
			{
				test: /\.(js|jsx)$/,
				use: "babel-loader",
				exclude: /node_modules/,
			},
			{
				test: /\.(ts|tsx)?$/,
				loader: "ts-loader",
				exclude: /node_modules/,
			},
			{
				test: /\.(sass|scss|css)$/,
				use: ["style-loader", "css-loader", "sass-loader"],
			},
			{
				test: /\.svg$/,
				use: "file-loader",
			},
			{
				test: /\.png$/,
				use: [
					{
						loader: "url-loader",
						options: {
							mimetype: "image/png",
						},
					},
				],
			},
		],
	},
	resolve: {
		extensions: [".js", ".jsx", ".tsx", ".ts"],
	},
	plugins: [
		new webpack.ContextReplacementPlugin(/moment[\/\\]locale$/, /en/),
		new HtmlWebpackPlugin({
			template: "src/index.html",
			inject: true,
			appMountId: "app",
		}),
	],
	optimization: {
		runtimeChunk: "single",
		splitChunks: {
			cacheGroups: {
				vendor: {
					test: /[\\\/]node_modules[\\\/]/,
					name: "vendors",
					chunks: "all",
				},
			},
		},
	},
}

module.exports = config
