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
	client *binance.Client
	db     *gorm.DB
	bc     *beary.Channel
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

	// initail bearychat
	bc = beary.NewChannel()
	bc.Info("Binance Collector 启动了")

	// 设置 log 输出的时间格式带微秒
	log.SetFlags(log.Lmicroseconds)
}

// Run a binance client to collect historical trades
// NOTICE: 国内的 IP 无法访问 binance 的 API
func Run() {
	defer db.Close()

	rs := newRecords()

	var day int

	done := false

	for !done {
		data := make([]*trade, 0, 20*1000)
		// 一口气获取 20 次，然后统一保存
		for i := 0; i < 20; i++ {
			symbol, utc, id := rs.first()

			if day != dayOf(utc) {
				day = dayOf(utc)
				date := time.Unix(0, utc*1000000)
				msg := fmt.Sprintf("已经收集到了 %s 的数据。", date)
				bc.Info(msg)
			}

			data = append(data, request(symbol, id+1)...)

			last := len(data) - 1
			utc, id = data[last].UTC, data[last].ID

			rs.update(utc, id)

			date := time.Unix(0, utc*1000000)
			log.Printf("%s %s", symbol, date)

		}

		done = rs.isUpdated()

		go save(data)

	}

}

func dayOf(utc int64) int {
	return time.Unix(0, utc*1000000).Day()
}

func deal(symbol string, id int64) {
	var day int
	var utc int64
	var date time.Time
	ticker := time.NewTicker(time.Second)
	//
	for !isToday(date) {
		trades := request(symbol, id+1)
		//
		last := len(trades) - 1
		utc, id = trades[last].UTC, trades[last].ID
		//
		date = time.Unix(0, utc*1000000)
		if day != date.Day() {
			day = date.Day()
			msg := fmt.Sprintf("%s 收集到了 %s 的数据。", symbol, date)
			bc.Info(msg)
			log.Printf(msg)
		}
		log.Printf("%s %s", symbol, date)
		// 保存数据
		save(trades)
		// 每秒运行
		<-ticker.C
	}
}

func isToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() &&
		date.Month() == now.Month() &&
		date.Day() == now.Day()
}
