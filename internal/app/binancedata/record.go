package binancedata

import (
	"container/heap"
	"log"
	"sync"
)

type records struct {
	queue  *symbolQueue
	remain int
	load   int
	sync.RWMutex
}

// 添加 load 参数是因为
// 我的梯子太烂了，访问 binance 的延迟很高，而且偶尔还会无法访问。
// load 代表了正在试图访问 binance 的 goroutine 数量。
// 超过限制后，需要停止访问。
// 不然的话，积压的请求一起涌向 binance 服务器的时候，会被封 IP
// 建议把 load 设置成 12
func newRecords(symbols []string, load int) *records {
	return &records{
		queue:  newSymbolQueue(symbols),
		remain: len(symbols),
		load:   load,
	}
}

func (rs *records) isOverload() bool {
	rs.RLock()
	defer rs.RUnlock()
	return rs.remain-len(*rs.queue) > rs.load
}

func (rs *records) getRemain() int {
	rs.RLock()
	defer rs.RUnlock()
	return rs.remain
}

func (rs *records) decrement() {
	rs.Lock()
	defer rs.Unlock()
	rs.remain--
}

// symbolRecord 是 priorityQueue 中的元素
type symbolRecord struct {
	symbol string
	utc    int64
	id     int64
}

func newSymbolRecord(symbol string, utc, id int64) *symbolRecord {
	return &symbolRecord{
		symbol: symbol,
		utc:    utc,
		id:     id,
	}
}

func (rs *records) first() *symbolRecord {
	rs.RLock()
	defer rs.RUnlock()
	if len(*rs.queue) == 0 {
		return nil
	}
	r := (*rs.queue)[0]
	log.Printf("symbol: %s,\t ID: %12d, Time: %s", r.symbol, r.id, localTime(r.utc))
	return r
}

func (rs *records) pop() *symbolRecord {
	rs.Lock()
	res := heap.Pop(rs.queue).(*symbolRecord)
	rs.Unlock()
	return res
}

func (rs *records) push(r *symbolRecord) {
	rs.Lock()
	heap.Push(rs.queue, r)
	rs.Unlock()
}

func (rs *records) isDone() bool {
	rs.RLock()
	defer rs.RUnlock()
	return rs.remain == 0
}

// symbolQueue implements heap.Interface and holds entries.
type symbolQueue []*symbolRecord

func newSymbolQueue(symbols []string) *symbolQueue {
	res := make(symbolQueue, 0, len(symbols))
	for _, s := range symbols {
		tp := newTrade(s)
		if db.HasTable(tp) {
			db.Last(tp)
			log.Printf("已经从 %s 的表中获取了 UTC = %s， ID = %d\n", s, localTime(tp.UTC), tp.ID)
		} else {
			db.CreateTable(tp)
			log.Printf("已经创建 %s 的表。\n", s)
		}
		heap.Push(&res, newSymbolRecord(s, tp.UTC, tp.ID))
	}
	return &res
}

func (rs symbolQueue) Len() int { return len(rs) }

func (rs symbolQueue) Less(i, j int) bool {
	return rs[i].utc < rs[j].utc
}

func (rs symbolQueue) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Push 往 rs 中放 record
func (rs *symbolQueue) Push(x interface{}) {
	temp := x.(*symbolRecord)
	*rs = append(*rs, temp)
}

// Pop 从 rs 中取出最优先的 record
func (rs *symbolQueue) Pop() interface{} {
	temp := (*rs)[len(*rs)-1]
	*rs = (*rs)[0 : len(*rs)-1]
	return temp
}
