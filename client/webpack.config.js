// const {clientPort} = require('../config')
const path = require('path')
const webpack = require('webpack')
const clientPort = 8000

const config = {
	devtool: 'source-map',
	entry: {
		app: path.resolve(__dirname, 'app/index.js'),
		vendor: ['react', 'react-dom']
	},
	output: {
		path: path.resolve(__dirname, 'public'),
		filename: 'js/[name].js',
		publicPath: '/'
	},
	resolve: {
		alias: {
			app: path.resolve(__dirname, 'app')
		},
		extensions: ['.js']
	},
	module: {
		rules: [
			{
				test: /\.js$/,
				exclude: [
					path.resolve(__dirname, 'node_modules')
				],
				use: [{
					loader: 'babel-loader',
					options: {presets: ['es2015', 'stage-0', 'react']}
				}]
			}
		],
	},
	plugins: [
		new webpack.ProvidePlugin({
			React: 'react',
			ReactDOM: 'react-dom',
			moment: 'moment'
		}),
		new webpack.optimize.CommonsChunkPlugin({
			name: 'vendor'
		}),
		new webpack.DefinePlugin({
			'process.env': {
				NODE_ENV: JSON.stringify(process.env.NODE_ENV)
			}
		}),
	],
}

if (process.env.NODE_ENV === 'production') {
	config.plugins.push(
		new webpack.optimize.UglifyJsPlugin({
			compress: {
				screw_ie8: true,
				warnings: false
			}
		})
	)
} else {
	config.devServer = {
		contentBase: path.join(__dirname, 'public'),
		inline: true,
		host: '0.0.0.0',
		port: clientPort
	}
}

module.exports = config
