package binancecollector

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance"
)

// 获取历史交易记录
func request(symbol string, id int64) []*trade {
	var originals []*binance.Trade
	var err error
	originals, err = client.NewHistoricalTradesService().Symbol(symbol).FromID(id).Limit(1000).Do(context.TODO())
	for err != nil {
		msg := fmt.Sprintf("client get historycal trades service err: %s", err)
		bc.Fatal(msg)
		time.Sleep(time.Minute * 3)
		originals, err = client.NewHistoricalTradesService().Symbol(symbol).FromID(id).Limit(1000).Do(context.TODO())
	}
	res := make([]*trade, 0, 1000)
	for _, original := range originals {
		r := convert(original, symbol)
		res = append(res, r)
	}
	return res
}
