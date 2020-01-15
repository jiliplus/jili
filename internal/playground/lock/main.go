package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
)

// 本程序验证了 sync.Mutex 会在多个 Lock 请求中，均衡地分配机会
// 不会出现单个 goroutine 长期联系独占 Lock 的情况。
func main() {
	start := time.Now()
	var wg sync.WaitGroup
	num := 10
	wg.Add(num)
	var m sync.Mutex
	endTime := make([]time.Duration, 10)
	for i := 0; i < num; i++ {
		go func(i int) {
			count := 0
			for count < 20 {
				m.Lock()
				count++
				time.Sleep(time.Millisecond * 100)
				m.Unlock()
				// 注释掉下一行会发现，其影响很小
				runtime.Gosched()
			}
			endTime[i] = time.Since(start)
			wg.Done()
		}(i)
	}
	wg.Wait()
	sort.Slice(endTime, func(i, j int) bool {
		return endTime[i] < endTime[j]
	})
	for i := 0; i < num; i++ {
		fmt.Println(endTime[i])
	}
	fmt.Println("程序耗时：", time.Since(start))
}
