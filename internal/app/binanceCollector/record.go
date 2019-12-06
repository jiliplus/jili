package binancecollector

import (
	"container/heap"
	"log"
	"time"

	"github.com/aQuaYi/jili/internal/pkg/tools"
)

// record 是 priorityQueue 中的元素
type record struct {
	symbol string
	utc    int64
	id     int64
}

func newRecord(symbol string, utc, id int64) *record {
	return &record{
		symbol: symbol,
		utc:    utc,
		id:     id,
	}
}

// records implements heap.Interface and holds entries.
type records []*record

func newRecords() *records {
	symbols := allSymbols()

	if !tools.IsExist(dbName) {
		return brandNewRecords(symbols)
	}

	return newRecordsFromDB(symbols)
}

func brandNewRecords(symbols []string) *records {
	res := make(records, 0, len(symbols))
	for _, s := range symbols {
		heap.Push(&res, newRecord(s, 0, -1))
	}
	return &res
}

func newRecordsFromDB(symbols []string) *records {
	res := make(records, 0, len(symbols))
	for _, s := range symbols {
		tp := newTrade(s)
		db.Last(tp)
		if tp.UTC == 0 {
			heap.Push(&res, newRecord(s, 0, -1))
		} else {
			heap.Push(&res, newRecord(s, tp.UTC, tp.ID))
		}
	}
	return &res
}

func (rs records) Len() int { return len(rs) }

func (rs records) Less(i, j int) bool {
	return rs[i].utc < rs[j].utc
}

func (rs records) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Push 往 rs 中放 record
func (rs *records) Push(x interface{}) {
	temp := x.(*record)
	*rs = append(*rs, temp)
}

// Pop 从 rs 中取出最优先的 record
func (rs *records) Pop() interface{} {
	temp := (*rs)[len(*rs)-1]
	*rs = (*rs)[0 : len(*rs)-1]
	return temp
}

func (rs *records) first() (symbol string, id int64) {
	symbol = (*rs)[0].symbol
	id = (*rs)[0].id
	utc := (*rs)[0].utc
	log.Printf("the first symbol: %s, ID: %d, Time: %s", symbol, id, time.Unix(0, utc*1000000))
	return
}

func (rs *records) update(time, id int64) {
	(*rs)[0].utc = time
	(*rs)[0].id = id
	heap.Fix(rs, 0)
}
