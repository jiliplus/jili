package bar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newBar(t *testing.T) {
	ast := assert.New(t)
	//
	now := time.Now()
	open, high, low, close := 1., 8., 0., 4.
	ticks := []float64{open, high, low, close}

	bar := newBar(now, ticks)

	ast.Equal(now, bar.begin)
	ast.Equal(open, bar.open)
	ast.Equal(high, bar.high)
	ast.Equal(low, bar.low)
	ast.Equal(close, bar.close)
}
