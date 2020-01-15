package main

import (
	"log"
	"sync"
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

	log.Println("")
	var wg sync.WaitGroup
	wg.Add(2)
	timer := time.AfterFunc(time.Second*5, func() {
		log.Println("in After")
		wg.Done()
	})
	timer.Reset(time.Second * 2)
	go func() {
		log.Println(<-timer.C)
		wg.Done()
	}()
	wg.Wait()
}
