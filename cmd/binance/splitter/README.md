# Binance 数据库分割

## 使用方法

```toml
APIKey = ""
SecretKey = ""
```

1. 填写并保存以上内容到 `binance.toml` 文件中。
1. 创建 `../data` 目录
1. 运行 `main.go` 程序。

会把 `binance.sqlite3` 按照电脑设置的时区，分割到月