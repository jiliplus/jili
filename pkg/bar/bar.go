package bar

import (
	"time"

	"github.com/jujili/exchange/name"
)

// Bar 实现了 k 线的相关方法
type Bar struct {
	Begin                  time.Time
	Open, High, Low, Close float64
	Volume                 float64
	// 以下属性可以通过其他方式获取
	// 为了节约空间，不要保存在数据库中
	Symbol   string
	Exchange name.Exchange
	Interval time.Duration
}

func newBar(begin time.Time, ticks []float64) Bar {
	open, high, low := ticks[0], ticks[0], ticks[0]
	close := ticks[len(ticks)-1]
	for _, t := range ticks {
		high = maxFloat64(high, t)
		low = minFloat64(low, t)
	}
	return Bar{
		Begin: begin,
		Open:  open,
		High:  high,
		Low:   low,
		Close: close,
	}
}
