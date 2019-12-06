package binancecollector

import (
	"log"
	"strconv"

	"github.com/adshao/go-binance"
)

// v1 data struct
type trade struct {
	ID           int64
	Price        float64
	Quantity     float64
	UTC          int64
	IsBuyerMaker bool
	IsBestMatch  bool
	Symbol       string `gorm:"-"` // 本字段不会保存到数据库
}

func convert(t *binance.Trade, symbol string) *trade {
	p, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		log.Fatal("convert t.Price err:", err)
	}
	q, err := strconv.ParseFloat(t.Quantity, 64)
	if err != nil {
		log.Fatal("convert t.Quantity err:", err)
	}
	return &trade{
		ID:           t.ID,
		Price:        p,
		Quantity:     q,
		UTC:          t.Time,
		IsBuyerMaker: t.IsBuyerMaker,
		IsBestMatch:  t.IsBestMatch,
		Symbol:       symbol,
	}
}

func newTrade(symbol string) *trade {
	return &trade{
		Symbol: symbol,
	}
}

func (t *trade) TableName() string {
	return t.Symbol
}
