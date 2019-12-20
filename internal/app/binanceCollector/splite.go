package binancecollector

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

// MDB = month database
var MDB = make(map[string]*gorm.DB, 1024)
var mDBmu sync.Mutex
var wg sync.WaitGroup

// Split db according with month
func Split() {
	log.Println("In Split now.")

	symbols := allSymbols()

	for _, symbol := range symbols {
		if !db.HasTable(symbol) {
			msg := fmt.Sprintf("新出现了交易对 %s", symbol)
			log.Println(msg)
			bc.Info(msg)
			continue
		}
		//
		msg := fmt.Sprintf("split %s now", symbol)
		bc.Info(msg)
		log.Println(msg)

		var count int
		db.Table(symbol).Count(&count)
		log.Printf("%s 表，一共有 %d 条数据\n", symbol, count)

		tradesChan := saver(symbol)

		// limit 代表了复制数据到内存中的数量
		// 取决于电脑内存的大小，我的电脑是 8 GB 的内存。
		// 设置成这么大可以不动用 Swap
		limit := 500 * 10000
		for offset := 0; offset <= count; offset += limit {
			var trades []*trade
			db.Table(symbol).Offset(offset).Limit(limit).Scan(&trades)
			tradesChan <- trades
		}
		close(tradesChan)
	}
	wg.Wait()
}

func newTmp() []*trade {
	// capacity 代表了一次写入数据库的最大数量
	capacity := 100 * 10000
	return make([]*trade, 0, capacity)
}

func saver(symbol string) chan<- []*trade {
	tradesChan := make(chan []*trade, 10)
	month := time.Month(0)
	tmp := newTmp()
	go func() {
		wg.Add(1)
		for ts := range tradesChan {
			for _, t := range ts {
				t.Symbol = symbol
				date := localTime(t.UTC)
				if month != date.Month() || len(tmp) == cap(tmp) {
					month = date.Month()
					saveDay(tmp)
					tmp = newTmp()
				}
				tmp = append(tmp, t)
			}
		}
		saveDay(tmp)
		wg.Done()
	}()
	return tradesChan
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
