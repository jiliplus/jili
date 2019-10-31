package jili

import (
	"time"

	"../internal/tools"
)

type bar struct {
	begin                  time.Time
	open, high, low, close float64
}

func newBar(begin time.Time, ticks []float64) bar {
	open, high, low := ticks[0], ticks[0], ticks[0]
	close := ticks[len(ticks)-1]
	for _, t := range ticks {
		high = tools.maxFloat64(high, t)
		low = tools.minFloat64(low, t)
	}
	return bar{
		begin: begin,
		open:  open,
		high:  high,
		low:   low,
		close: close,
	}
}
