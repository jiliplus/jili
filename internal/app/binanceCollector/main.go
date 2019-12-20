package binancecollector

import (
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

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

// Run a binance client to collect historical trades
func Run() {
	defer db.Close()

	bc.Warning("NOTICE: 国内的 IP 无法访问 binance 的 API")

	rs := newRecords(allSymbols(), 12)

	// 访问限制是，每分钟 240 次。
	// 也就是每次的间隔时间为 250 毫秒
	ticker := time.NewTicker(250 * time.Millisecond)

	check := checkFunc()

	for !rs.isDone() {
		check(rs.first())
		// 由于网络出现了较大延迟
		// 导致前面还有很多没有处理完的。
		// 所以，跳过这一次
		if rs.isOverload() {
			log.Println("request 已经 overload 了，跳过这一次。")
		} else {
			go do(rs)
		}
		<-ticker.C
	}

	msg := "全部 symbol 都已经收集完毕"
	bc.Info(msg)
	log.Println(msg)
}

func do(rs *records) {
	r := rs.pop()
	symbol, utc, id := r.symbol, r.utc, r.id
	trades, err := request(symbol, id+1)
	if err != nil {
		msg := fmt.Sprintf("client get historycal trades service err: %s", err)
		bc.Fatal(msg)
		log.Println(msg)
		rs.push(r)
		return
	}
	// 无论获取了多少数据，都可以发送到 tradesChan
	tradesChan <- trades
	//
	size := len(trades)
	if size < 1000 {
		rs.decrement()
		msg := fmt.Sprintf("%s 从 %s 和 %d 获取的数据长度为 %d, 决定不再放回去。rs.remain = %d", symbol, localTime(utc), id, size, rs.getRemain())
		bc.Fatal(msg)
		log.Println(msg)
		return
	}
	last := size - 1
	r.utc, r.id = trades[last].UTC, trades[last].ID
	rs.push(r)
}

func localTime(UTCInMillionSecond int64) time.Time {
	return time.Unix(0, UTCInMillionSecond*1000000)
}

func checkFunc() func(r *symbolRecord) {
	count := int64(0)
	seconds := int64(24 * 60 * 60) // seconds of one day
	return func(r *symbolRecord) {
		if r == nil {
			return
		}
		days := r.utc / seconds
		if count < days {
			count = days
			msg := fmt.Sprintf("%s 已经收集到 %s, ID = %d\n", r.symbol, localTime(r.utc), r.id)
			log.Println(msg)
			bc.Info(msg)
		}
	}
}
