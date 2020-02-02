# Binance 命令

- [`collector`](#collector)
- [`splitter`](#splitter)
- [`merger`](#merger)
- [按表分割数据库](#%e6%8c%89%e8%a1%a8%e5%88%86%e5%89%b2%e6%95%b0%e6%8d%ae%e5%ba%93)
	- [数据库准备](#%e6%95%b0%e6%8d%ae%e5%ba%93%e5%87%86%e5%a4%87)
	- [操作命令](#%e6%93%8d%e4%bd%9c%e5%91%bd%e4%bb%a4)

## [`collector`](./collector/README.md)

收集 Binance 交易所的所有历史交易 Tick 到 `binance.sqlite3` 数据库中。以交易对的 `symbol` 作为表名。

请加 <a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=7f61280435c41608fb8cb96cf8af7d31ef0007c44b223c9e3596ce84dec329bc"><img border="0" src="https://img.shields.io/badge/QQ%20群-23%2053%2000%2093-blue.svg" alt="jili 交流群" title="jili 交流群"></a> 获取最新数据的百度网盘下载地址。

## [`splitter`](./splitter/README.md)

把 `collector` 保存的数据库，按月分割到不同的数据库。请注意，划分标准是运行程序电脑所设置的时区。

## [`merger`](./merger/README.md)

把 `splitter` 划分的数据库，重新组合成单一的 `binance.sqlite3` 文件。

## 按表分割数据库

把 `binance.sqlite3` 中表 `BTCUSDT` 完整地复制到 `btcusdt.sqlite3` 中表 `TICK` 的方法。

### 数据库准备

1. <https://sqlitebrowser.org> 下载并安装 `DB Browser for SQLite` 软件。
1. 使用 `DB Browser for SQLite` 打开 `binance.sqlite3` 文件。
1. `数据库结构` → `名称` → `表`，选中 `BTCUSDT` 表，右键，点击 `复制 Create 语句`。
1. `DB Browser for SQLite` → `文件` → `新建数据库`，命名为 `btcusdt.sqlite3`，save。
1. `执行 SQL` → `SQL 1`，粘贴刚刚复制的 Create 语句。
1. 把语句中的表名 `BTCUSDT` 修改成为 `TICK`，执行语句。

此时，已经准备好了 `btcusdt.sqlite3` 数据库文件及其 `TICK` 表。

### 操作命令

由于 `binance.sqlite3` 和 `btcusdt.sqlite3` 放在同一个目录下，下方只出现了文件名。
也可以使用绝对路径替换文件名。
在命令行中进入 `btcusdt.sqlite3` 所在目录，并输入命令：

> sqlite3 btcusdt.sqlite3
> sqlite> ATTACH 'binance.sqlite3' AS SRC;
> sqlite> .database
> main: .../btcusdt.sqlite3
> SRC: .../binance.sqlite3
> sqlite> INSERT INTO tick SELECT * FROM SRC.btcusdt;
> sqlite> .quit
> $

操作完毕
