var jsLoaders = ['babel?stage=0&optional=runtime'];

if (process.env.NODE_ENV !== 'production') {
  jsLoaders.unshift('react-hot');
}

module.exports = {
  devtool: 'source-map',
  entry: {
    browser: ['./app/index-browser.js'],
    server: ['./app/index-server.js'],
  },
  output: {
    path: './public',
    filename: 'bundle.[name].js',
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        loaders: jsLoaders,
        exclude: /node_modules/,
      },
      {
        test: /\.css$/,
        loaders: ['style', 'css?modules'],
      },
    ],
  },
};
