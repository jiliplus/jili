package binancecollector

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance"
	"github.com/pelletier/go-toml"
)

const (
	configFile = "binance.toml"
)

// Run a binance client to collect historical trades
func Run() {
	config, err := toml.LoadFile(configFile)
	if err != nil {
		msg := fmt.Sprintf("无法导入 %s，%s", configFile, err)
		panic(msg)
	}
	a, s := config.Get("APIKey").(string), config.Get("SecretKey").(string)
	fmt.Printf("APIKey   : %s\n", a)
	fmt.Printf("SecretKey: %s\n", s)
	client := binance.NewClient(a, s)
	// client.BaseURL = "api.binance.co"
	res, err := client.NewHistoricalTradesService().Symbol("ETHBTC").FromID(0).Limit(1000).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	r := res[0]
	fmt.Printf("%d,%d,%s\n", r.ID, r.Time, time.Unix(0, r.Time*1000000))
}
