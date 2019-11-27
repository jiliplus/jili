# 交易记录

这里记录了交易所的历史成交数据。由于数据量太大，我放到了百度网盘。
请加
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=7f61280435c41608fb8cb96cf8af7d31ef0007c44b223c9e3596ce84dec329bc"><img border="0" src="https://img.shields.io/badge/QQ%20群-23%2053%2000%2093-blue.svg" alt="jili 交流群" title="jili 交流群"></a> 获取下载地址。

![jili 交流群](https://user-images.githubusercontent.com/6028869/68080839-5d677700-fe3e-11e9-9e1d-9eeb71e5832c.jpg)

## 数据格式

数据使用以下方式保存。以 `binance` 为例：

- 全部数据保存在 SQLite3 数据库文件中。
- 所有路径、数据库文件名、表名都是小写。
- 全部数据放在本目录的 `binance` 文件夹下面。
- 按照**北京时间**，数据按天放入各个数据文件中。今天的 `binance` 数据文件名是 `2019-11-08.binance.sqlite3`
- 各个交易对的数据放在相应的表中。

<!-- TODO: 提供工具，把每天的汇总成一个文件，并计算相关的表 -->
