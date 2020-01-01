package binancedata

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

// MDB = month database
var MDB = make(map[string]*gorm.DB, 1024)

// var mDBmu sync.Mutex
var wg sync.WaitGroup

// Split db according with month
func Split() {

	log.Println("In Split now.")

	symbols := getSymbols()

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

		// db.Table(symbol).Count(&count)
		msg = fmt.Sprintf("%s 表，一共有 %d 条数据\n", symbol, count(symbol))
		bc.Info(msg)
		log.Println(msg)

		tradesChan := saver2(symbol)
		source(tradesChan, symbol)
		wg.Wait()
	}
}

func newTmp() []*trade {
	// capacity 代表了一次写入数据库的最大数量
	capacity := 10 * 10000
	return make([]*trade, 0, capacity)
}

func saver2(symbol string) chan<- *trade {
	tradesChan := make(chan *trade, 100)
	month := time.Month(0)
	tmp := newTmp()
	go func() {
		wg.Add(1)
		for t := range tradesChan {
			date := localTime(t.UTC)
			if month != date.Month() || len(tmp) == cap(tmp) {
				month = date.Month()
				save2disk(tmp)
				tmp = newTmp()
			}
			tmp = append(tmp, t)
		}
		save2disk(tmp)
		wg.Done()
	}()
	return tradesChan
}

// TODO: 删除此处内容
func saver(symbol string) chan<- []*trade {
	tradesChan := make(chan []*trade, 100)
	month := time.Month(0)
	tmp := newTmp()
	go func() {
		wg.Add(1)
		for ts := range tradesChan {
			for _, t := range ts {
				date := localTime(t.UTC)
				if month != date.Month() || len(tmp) == cap(tmp) {
					month = date.Month()
					save2disk(tmp)
					tmp = newTmp()
				}
				tmp = append(tmp, t)
			}
		}
		save2disk(tmp)
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

func count(symbol string) uint {
	var result = struct {
		Count uint
	}{}
	// db.Raw("select count(map) as maps from (select distinct map from matches)").Scan(&result)
	sql := fmt.Sprintf("SELECT COUNT(ROWID) AS count FROM %s", symbol)
	db.Raw(sql).Scan(&result)
	return result.Count
}
