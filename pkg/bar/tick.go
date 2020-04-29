package bar

import (
	"time"

	"github.com/jujili/exchange/name"
)

// Tick 包含了 tick data
type Tick struct {
	Exchange name.Exchange
	Symbol   string
	Date     time.Time
	Price    float64
	Volume   float64
	Type     string
}
