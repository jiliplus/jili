# Clock 模块

- [总体思路](#%e6%80%bb%e4%bd%93%e6%80%9d%e8%b7%af)
- [真实的 Clock](#%e7%9c%9f%e5%ae%9e%e7%9a%84-clock)
- [虚拟的 Clock](#%e8%99%9a%e6%8b%9f%e7%9a%84-clock)
- [死锁](#%e6%ad%bb%e9%94%81)
- [小技巧](#%e5%b0%8f%e6%8a%80%e5%b7%a7)

## 总体思路

clock 模块分为真实与虚拟的两个部分。真实的部分是以 time 和 context 标准库为基础，封装成了 Clock 接口。
然后，虚拟部分也实现了 Clock 接口，只是这一部分可以人为的操控时间的改变。

## 真实的 Clock

`time.Now` 提供了计算机中的当前时间。`time` 标准库还有另外一些函数，需要用到 `time.Now` 提供的时间。比如，`time.Sleep` 是在当前时间的基础上，休眠一段时间。

真实的 Clock 就是把 `time` 和 `context` 标准库中这些以当前时间为基础进行运算的函数，封装成了 `realClock` 的方法。

代码全部放在 `real-*.go` 中。

## 虚拟的 Clock

想要虚拟出时间流逝，其实光靠 `mockClock` 是不够的。还需要有

- `mockTimer`：虚拟的计时器
- `mockTicker`：虚拟的间隔器
- `mockContext`：虚拟的上下文

随着时间的流逝，这三种对象需要在恰当的时间点进行触发。

代码全部放在 `mock-*.go` 中。

<!-- TODO: 删除死锁，因为我在程序中，已经消除这种情况了。 -->

## 死锁

在使用 `mockTicker` 的时候，

```go
for {
	// ...
	<- mockTicker.C
	now := mock.Now()
	// ...
}
```

以上代码很容易造成死锁。因为当 `mockClock` 的调整间隔大于 `mockTicker` 的周期的时候，会在 `mockClock` 的临界区内多次往 `mockTicker.C` 内发送数据。`mock.Now()` 又需要进入临界区，才能让 `<- mockTicker.C` 第二次获取数据。

正确的用法是，

```go
for {
	// ...
	now := <- mockTicker.C
	// ...
}
```

## 小技巧

设置一个 ticker 到 mock 中，这样在 Mock.Set 的时候，会以 ticker 的周期为步进，跳到 Set 的时间点上，可以让 Mock 更像 time
