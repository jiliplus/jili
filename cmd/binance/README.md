# Binance 命令

## [`collector`](./collector/README.md)

收集 Binance 交易所的所有历史交易 Tick 到 `binance.sqlite3` 数据库中。以交易对的 `symbol` 作为表名。

请加 <a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=7f61280435c41608fb8cb96cf8af7d31ef0007c44b223c9e3596ce84dec329bc"><img border="0" src="https://img.shields.io/badge/QQ%20群-23%2053%2000%2093-blue.svg" alt="jili 交流群" title="jili 交流群"></a> 获取最新数据的百度网盘下载地址。

## [`splitter`](./splitter/README.md)

把 `collector` 保存的数据库，按月分割到不同的数据库。请注意，划分标准是运行程序电脑所设置的时区。

## [`merger`](./merger/README.md)

把 `splitter` 划分的数据库，重新组合成单一的 `binance.sqlite3` 文件。