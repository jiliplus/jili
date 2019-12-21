package binancecollector

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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

func (t trade) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s ，", t.Symbol))
	sb.WriteString(fmt.Sprintf("ID: %d", t.ID))
	sb.WriteString(fmt.Sprintf("价格: %12f", t.Price))
	sb.WriteString(fmt.Sprintf("数量: %12f", t.Quantity))
	sb.WriteString(fmt.Sprintf("时间: %s", localTime(t.UTC)))
	sb.WriteString(fmt.Sprintf("Buyer Maker: %t", t.IsBuyerMaker))
	sb.WriteString(fmt.Sprintf("Best Match: %t", t.IsBestMatch))
	return sb.String()
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

func (t *trade) monthDBName() string {
	date := localTime(t.UTC)
	return fmt.Sprintf("../data/%d%02d.binance.sqlite3", date.Year(), date.Month())
}
