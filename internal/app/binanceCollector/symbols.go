package binancecollector

import (
	"context"
	"fmt"
	"log"
)

func getSymbols() []string {
	res := make([]string, 0, 1024)

	info, err := client.NewExchangeInfoService().Do(context.TODO())
	if err != nil {
		log.Fatal("Binance NewExchangeInfoService err:", err)
	}
	for _, s := range info.Symbols {
		res = append(res, s.Symbol)
	}

	fmt.Println("symbol 的数量是 ", len(info.Symbols))
	return res
}
