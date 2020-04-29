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

	ast.Equal(now, bar.Begin)
	ast.Equal(open, bar.Open)
	ast.Equal(high, bar.High)
	ast.Equal(low, bar.Low)
	ast.Equal(close, bar.Close)
}
