package binancecollector

import (
	"context"
	"log"
)

func request(symbol string, id int) []*trade {
	// 获取历史交易记录
	originals, err := client.NewHistoricalTradesService().Symbol(symbol).FromID(int64(id)).Limit(1000).Do(context.TODO())
	if err != nil {
		log.Fatal("client get historycal trades service err:", err)
	}
	res := make([]*trade, 0, 1000)
	for _, original := range originals {
		r := convert(original, symbol)
		res = append(res, r)
	}
	return res
}
