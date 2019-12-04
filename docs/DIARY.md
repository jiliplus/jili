# 开发日记

在这里记录我开发这个项目的点点滴滴。

## 2019-12-04

已经过了快一个月，日记都要变成月报了。

这段时间 jili 项目进度很慢，主要有两个原因，一是，身体有 3 周不适，无法编程；二是，一直没有想好要做成什么样子。

就在刚才，我想清楚了。在本地的电脑上收集历史数据，然后，利用工具软件来转换总的数据和每天的数据。

## 2019-11-07

获取了 "ETHBTC" 的第一条交易记录的时间

> 0,1500004804757,2017-07-14 12:00:04.757 +0800 CST

## 2019-11-06

在粗略地浏览了 [GoEx](https://github.com/nntaoli-project/GoEx) 和 [go-binance](https://github.com/adshao/go-binance) 后。我还是决定自己写交易所 `API` 的封装。

在目前的 3 大交易所中，只有 `binance` 提供历史交易数据。所以，从它开始写起。

## 2019-11-03

这个坑其实已经开了很久了，现在终于想要把它填满了。

前几天，利用 [Standard Go Project Layout](https://github.com/golang-standards/project-layout) 和 [golang开发目录结构](https://segmentfault.com/a/1190000012926524) 把项目的目录整理了一下，看起来果然整洁多了。
