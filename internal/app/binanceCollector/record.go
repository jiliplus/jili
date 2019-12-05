package binancecollector

import (
	"container/heap"

	"github.com/aQuaYi/jili/internal/pkg/tools"
)

// record 是 priorityQueue 中的元素
type record struct {
	symbol string
	time   int
	id     int
	// index 是 record 在 heap 中的索引号
	index int
}

func newRecord(symbol string, time, id int) *record {
	return &record{
		symbol: symbol,
		time:   time,
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

	return
}

func brandNewRecords(symbols []string) *records {
	res := make(records, 0, len(symbols))
	for _, s := range symbols {
		heap.Push(&res, newRecord(s, 0, -1))
	}
	return &res
}

func (rs records) Len() int { return len(rs) }

func (rs records) Less(i, j int) bool {
	return rs[i].time < rs[j].time
}

func (rs records) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
	rs[i].index = i
	rs[j].index = j
}

// Push 往 rs 中放 record
func (rs *records) Push(x interface{}) {
	temp := x.(*record)
	temp.index = len(*rs)
	*rs = append(*rs, temp)
}

// Pop 从 rs 中取出最优先的 record
func (rs *records) Pop() interface{} {
	temp := (*rs)[len(*rs)-1]
	temp.index = -1 // for safety
	*rs = (*rs)[0 : len(*rs)-1]
	return temp
}

func (rs *records) first() (symbol string, id int) {
	symbol = (*rs)[0].symbol
	id = (*rs)[0].id
	return
}

func (rs *records) update(time, id int) {
	(*rs)[0].time = time
	(*rs)[0].id = id
	heap.Fix(rs, 0)
}
