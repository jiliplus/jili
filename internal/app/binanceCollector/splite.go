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
		}
		//
		msg := fmt.Sprintf("split %s now", symbol)
		bc.Info(msg)
		log.Println(msg)
		//
		tp := newTrade(symbol)
		db.Last(tp)
		maxID := tp.ID
		log.Printf("%s 的 max ID = %d", symbol, maxID)

		sql := fmt.Sprintf("SELECT * FROM %s", symbol)
		rows, err := db.Raw(sql).Rows()
		if err != nil {
			panic("Split Raw Rows Err:" + err.Error())
		}

		tradesChan := saver()
		for rows.Next() {
			var t trade
			if err := rows.Scan(&t.ID, &t.Price, &t.Quantity, &t.UTC, &t.IsBuyerMaker, &t.IsBestMatch); err != nil {
				panic("rows.Scan err:" + err.Error())
			}
			t.Symbol = symbol
			tradesChan <- &t
		}
		rows.Close()
		// if rows.Close().Error() != "" {
		// panic("rows.Close().Error():" + rows.Close().Error())
		// }
		close(tradesChan)
	}
	wg.Wait()
}

func newTmp() []*trade {
	return make([]*trade, 0, 20*10000)
}

func saver() chan<- *trade {
	tradesChan := make(chan *trade, 2000)
	month := time.Month(0)
	tmp := newTmp()
	go func() {
		wg.Add(1)
		for t := range tradesChan {
			date := localTime(t.UTC)
			if month != date.Month() || len(tmp) == 200000 {
				month = date.Month()
				saveDay(tmp)
				tmp = newTmp()
			}
			tmp = append(tmp, t)
		}
		saveDay(tmp)
		wg.Done()
	}()
	return tradesChan
}
