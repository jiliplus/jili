# 交易所

`exchange` 模块封装了各大交易所的 API。

 | 交易所         | 行情接口 | 交易接口 | 版本号 |
 | -------------- | -------- | -------- | ------ |
 | binance-cn.com | [ ]      | [ ]      |        |

## `exchange` 的功能

<!-- TODO: -->

- [ ] 使用 `.toml` 文件配置交易所。 请确保 `*.toml` 已经在你的 `.gitignore` 文件中。
- [ ] 符合交易所要求的限速设置。采用两级令牌桶，重要的交易接口使用上层桶，确保资金安全。次要的行情接口使用下层桶，为交易让路。
- [ ] 为了保证限速的全局性，各个交易所采用 `singleton` 设计模式。