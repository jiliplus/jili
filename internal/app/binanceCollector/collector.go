package binancecollector

import (
	"context"
	"fmt"
	"log"

	"github.com/adshao/go-binance"
	"github.com/pelletier/go-toml"
)

const (
	configFile = "binance.toml"
)

func run() {
	tree, err := toml.LoadFile(configFile)
	if err != nil {
		msg := fmt.Sprintf("无法导入 %s，%s", configFile, err)
		panic(msg)
	}
	client := binance.NewClient(
		tree.GetPosition("APIkey").String(),
		tree.GetPosition("SecretKey").String())
	res, err := client.NewHistoricalTradesService().Symbol("BNBETH").FromID(0).Limit(1000).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res)
	}
}
