# bucket 模块

## 特色

除了普通 `Token Bucket` 算法模块提供的 `Wait` 方法。本模块，还提供了一个 `Hurry` 方法。

`Hurry` 能使用预先保留的 token，可以避免一些优先级高的请求在等待。
