const path = require('path');

module.exports = {
    "mode": "development",
    "entry": "./src/index.ts",
    // "output": {
        // "path": __dirname+'/dist',
        // "filename": "[name].[chunkhash:8].js"
    // },
    "devtool": "source-map",
    devServer: {
        contentBase: './dist',
        hot: true
    },
    "module": {
        "rules": [
            {
                "enforce": "pre",
                "test": /\.(js|ts)x?$/,
                "exclude": /node_modules/,
                "use": "eslint-loader"
            },
            {
                "test": /\.tsx?$/,
                "exclude": /node_modules/,
                "use": "ts-loader"
            },
            {
                "test": /\.scss$/,
                "use": [
                    "style-loader",
                    "css-loader",
                    "sass-loader"
                ]
            }
        ]
    }
};
