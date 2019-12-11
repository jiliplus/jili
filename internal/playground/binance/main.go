package main

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

var (
	client *binance.Client
)

func init() {
	// initial client
	config, err := toml.LoadFile(configFile)
	if err != nil {
		msg := fmt.Sprintf("无法导入 %s，%s", configFile, err)
		panic(msg)
	}
	a, s := config.Get("APIKey").(string), config.Get("SecretKey").(string)
	fmt.Printf("APIKey   : %s\n", a)
	fmt.Printf("SecretKey: %s\n", s)
	client = binance.NewClient(a, s)
	fmt.Println("client 初始化完毕")

	// 设置 log 输出的时间格式带微秒
	log.SetFlags(log.Lmicroseconds)
}

// NOTICE: 国内的 IP 无法访问 binance 的 API

func main() {
	eis, err := client.NewExchangeInfoService().Do(context.TODO())
	if err != nil {
		log.Println("exchange info service err:", err)
	}
	fmt.Println(eis.RateLimits)
}
