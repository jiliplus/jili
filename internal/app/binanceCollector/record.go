package binancecollector

import (
	"container/heap"
	"log"
	"time"
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
	res := make(records, 0, len(symbols))
	for _, s := range symbols {
		tp := newTrade(s)
		if db.HasTable(tp) {
			db.Last(tp)
			heap.Push(&res, newRecord(s, tp.UTC, tp.ID))
		} else {
			db.CreateTable(tp)
			heap.Push(&res, newRecord(s, 0, 0))
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

func (rs *records) update(utc, id int64) {
	(*rs)[0].utc = utc
	(*rs)[0].id = id
	heap.Fix(rs, 0)
}

func (rs *records) isUpdated() bool {
	latest := time.Unix(0, 1000000*(*rs)[0].utc)
	return time.Minute*3 > time.Since(latest)
}
