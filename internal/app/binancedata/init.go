package binancedata

import (
	"fmt"
	"log"

	"github.com/adshao/go-binance"
	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"

	"github.com/aQuaYi/jili/internal/pkg/beary"
)

// TODO: 把 init 还原到 main.go
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

	// initial db
	db, err = gorm.Open("sqlite3", dbName)
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Printf("%s 数据库已经打开\n", dbName)

	// initial bearychat
	bc = beary.NewChannel()
	bc.Info("Binance Collector 启动了")

	// 设置 log 输出的时间格式带微秒
	log.SetFlags(log.Lmicroseconds)

	initialDBSaver()
}

func initialDBSaver() {
	tradesChan = make(chan []*trade, 30)
	go func() {
		data := make([]*trade, 0, 22000)
		for ts := range tradesChan {
			data = append(data, ts...)
			if len(data) >= 20000 {
				save(data)
				data = make([]*trade, 0, 22000)
			}
		}
	}()
}
