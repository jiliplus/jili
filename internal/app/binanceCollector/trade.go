package binancecollector

import "github.com/adshao/go-binance"

// v1 data struct
type trade struct {
	ID           int64
	Price        string
	Quantity     string
	UTC          int64
	IsBuyerMaker bool
	IsBestMatch  bool
	Symbol       string `gorm:"-"` // 本字段不会保存到数据库
}

func convert(t *binance.Trade, symbol string) *trade {
	return &trade{
		ID:           t.ID,
		Price:        t.Price,
		Quantity:     t.Quantity,
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
