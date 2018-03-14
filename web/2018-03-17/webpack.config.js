const path = require('path');

module.exports = {
  // モードの設定、v4系以降はmodeを指定しないと、webpack実行時に警告が出る
  mode: 'production',
  // エントリーポイントの設定
  entry: ['babel-polyfill', './src/main.js'],
  // 出力の設定
  output: {
    // 出力するファイル名
    filename: 'bundle.js',
    // 出力先のパス（v2系以降は絶対パスを指定する必要がある）
    path: path.resolve('public')
  },
  module: {
    rules: [{
      test: /\.js$/,
      use: [{
        loader: "babel-loader",
        options:{
          presets: [['env', {'modules': false}]]
        }
      }]
    }]
  }
};
