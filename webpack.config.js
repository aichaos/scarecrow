var webpack = require("webpack");

module.exports = {
	entry: "./http/src/index.js",

	output: {
		path: "./http/public/static/js",
		filename: "scarecrow.js",
		publicPath: "",
	},

	module: {
		loaders: [
			{
				test: /\.jsx?$/,
				exclude: /node_modules/,
				loader: 'babel-loader?presets[]=es2015&presets[]=react'
			}
		]
	},

	plugins: process.env.NODE_ENV === "production" ? [
		new webpack.optimize.DedupePlugin(),
		new webpack.optimize.OccurrenceOrderPlugin(),
		new webpack.optimize.UglifyJsPlugin()
	] : []
}
