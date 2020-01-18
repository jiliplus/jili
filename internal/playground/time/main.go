package main

import (
	"log"
	"time"
)

func main() {

	// passedTimer := time.NewTimer(-24 * time.Hour)
	// log.Println("ptc:", <-passedTimer.C)
	// tick := time.Tick(time.Millisecond * 500)
	// var wg sync.WaitGroup
	// wg.Add(10)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		time.Sleep(time.Second)
	// 		log.Println(<-tick)
	// 		wg.Done()
	// 	}
	// }()
	// wg.Wait()
	// t := time.NewTimer(time.Second)
	// t.Stop()
	// log.Println(time.AfterFunc(time.Second, func() {
	// 	log.Println("in AfterFunc")
	// }))
	// time.Sleep(time.Second)

	// log.Println("")
	// var wg sync.WaitGroup
	// wg.Add(1)
	// timer := time.AfterFunc(time.Second*5, func() {
	// 	log.Println("in After")
	// 	wg.Done()
	// })
	// timer.Reset(time.Second * 2)
	// go func() {
	// 	log.Println(<-timer.C)
	// 	// 两个 Done 并不会引发 deadlock
	// 	// 因为上一行永远不会执行
	// 	wg.Done()
	// }()
	// wg.Wait()

	log.Println("")
	timer := time.AfterFunc(time.Second, func() {
		log.Println("in After")
	})
	time.Sleep(time.Second * 2)
	log.Println(timer.Reset(time.Second))
	time.Sleep(time.Second * 2)
	log.Println("")
	log.Println(timer.Stop())
	log.Println(timer.Reset(time.Second))
	time.Sleep(time.Second * 2)
	log.Println("")
}
