package binancedata

import (
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance"
	"github.com/jinzhu/gorm"

	// 引入 sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jujili/jili/internal/pkg/beary"
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

// Collect a binance client to collect historical trades
func Collect() {
	defer db.Close()

	bc.Warning("NOTICE: 国内的 IP 无法访问 binance 的 API")

	rs := newRecords(getSymbols(), 8)

	// 访问限制是，每分钟 240 次。
	// 也就是每次的间隔时间为 250 毫秒
	ticker := time.NewTicker(275 * time.Millisecond)

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
	if r == nil {
		return
	}
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
	// 1000 是 binance 能返回的最大数据量
	if size < 1000 {
		// 决定不在放回，所以总的任务量需要减少
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

func localTime(UTCInMillisecond int64) time.Time {
	return time.Unix(0, UTCInMillisecond*1000000)
}

// symbolRecord 是按照时间顺序收集的
// 所以，需要 check 函数来检查是否已经收到了新的日期。
func checkFunc() func(r *symbolRecord) {
	count := int64(0)
	// milliseconds of one day
	// binance api 返回的时间，都是 unix 的毫秒。
	msOfOneDay := int64(24 * 60 * 60 * 1000) // 一天包含的毫秒数
	return func(r *symbolRecord) {
		if r == nil {
			return
		}
		days := r.utc / msOfOneDay // record 的 utc 日期所代表的天数
		if count < days {
			count = days
			msg := fmt.Sprintf("%s 已经收集到 %s, ID = %d", r.symbol, localTime(r.utc), r.id)
			log.Println(msg)
			bc.Info(msg)
		}
	}
}
