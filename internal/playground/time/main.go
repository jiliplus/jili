package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	tick := time.Tick(time.Millisecond * 500)
	var wg sync.WaitGroup
	wg.Add(10)
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println(<-tick)
			wg.Done()
		}
	}()
	wg.Wait()
	t := time.NewTimer(time.Second)
	t.Stop()
}
