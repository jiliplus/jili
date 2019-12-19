package binancecollector

import (
	"context"
)

// 获取历史交易记录
func request(symbol string, id int64) ([]*trade, error) {
	originals, err := client.NewHistoricalTradesService().Symbol(symbol).FromID(id).Limit(1000).Do(context.TODO())
	if err != nil {
		return nil, err
	}
	res := make([]*trade, 0, 1000)
	for _, original := range originals {
		r := convert(original, symbol)
		res = append(res, r)
	}
	return res, nil
}
