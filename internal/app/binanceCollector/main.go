package binancecollector

import (
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"

	"github.com/aQuaYi/jili/internal/pkg/beary"
)

const (
	configFile = "binance.toml"
	dbName     = "binance.sqlite3"
)

var (
	client     *binance.Client
	db         *gorm.DB
	bc         *beary.Channel
	tradesChan chan []*trade
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

	// initial tradesChan
	tradesChan = make(chan []*trade)
}

// Run a binance client to collect historical trades
// NOTICE: 国内的 IP 无法访问 binance 的 API
func Run() {
	defer db.Close()

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

	rs := newRecords()
	ethbtc := rs.ethbtc()

	var day int

	ticker := time.NewTicker(256 * time.Millisecond)

	for !rs.isUpdated() {
		// 访问限制是，每分钟 240 次。
		// 也就是每次的间隔时间为 250 毫秒
		// 我把 ticker 设置成 256 毫秒
		go deal(rs)

		rs.first()

		mu.RLock()
		utc := ethbtc.utc
		mu.RUnlock()
		if day != dayOf(utc) {
			day = dayOf(utc)
			date := time.Unix(0, utc*1000000)
			msg := fmt.Sprintf("ETHBTC 已经收集到了 %s 的数据。", date)
			bc.Info(msg)
			log.Println(msg)
		}

		<-ticker.C
	}

}

func deal(rs *records) {
	// 由于网络出现了较大延迟
	// 导致前面还有很多没有处理完的。
	// 所以，跳过这一次
	if rs.isDelayed() {
		msg := "rs 已经 delay 了。所以跳过这一次。"
		bc.Verbose(msg)
		log.Println(msg)
		return
	}
	r := rs.pop()
	symbol, id := r.symbol, r.id
	trades, err := request2(symbol, id+1)
	if err == nil {
		tradesChan <- trades
		last := len(trades) - 1
		utc, id := trades[last].UTC, trades[last].ID
		r.utc, r.id = utc, id
	} else {
		msg := fmt.Sprintf("client get historycal trades service err: %s", err)
		bc.Fatal(msg)
		log.Println(msg)
	}
	rs.push(r)
}

func dayOf(utc int64) int {
	return time.Unix(0, utc*1000000).Day()
}

// func deal(symbol string, id int64) {
// 	var day int
// 	var utc int64
// 	var date time.Time
// 	ticker := time.NewTicker(time.Second)
// 	//
// 	for !isToday(date) {
// 		trades := request(symbol, id+1)
// 		//
// 		last := len(trades) - 1
// 		utc, id = trades[last].UTC, trades[last].ID
// 		//
// 		date = time.Unix(0, utc*1000000)
// 		if day != date.Day() {
// 			day = date.Day()
// 			msg := fmt.Sprintf("%s 收集到了 %s 的数据。", symbol, date)
// 			bc.Info(msg)
// 			log.Printf(msg)
// 		}
// 		log.Printf("%s %s", symbol, date)
// 		// 保存数据
// 		save(trades)
// 		// 每秒运行
// 		<-ticker.C
// 	}
// }
