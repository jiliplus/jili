package jili

import (
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOr(t *testing.T) {
	Convey("如果 or 没有输入参数, 会 panic", t, func() {
		So(func() { Or() }, ShouldPanicWith, "Or 没有输入参数")
	})

	Convey("当有多个输入参数时，", t, func() {
		done := make(chan struct{})
		dones := []<-chan struct{}{done}
		for i := 1; i < 10; i++ {
			Convey(fmt.Sprintf("如果 or 了 %d 个 chan, ", len(dones)), func() {
				orDone := Or(dones...)
				var begin, end time.Time

				Convey("在关闭了输入参数中最后一个通道以后，", func() {
					var wg sync.WaitGroup
					wg.Add(2)

					go func() {
						<-orDone
						end = time.Now()
						wg.Done()
					}()

					go func() {
						begin = time.Now()
						close(done)
						wg.Done()
					}()

					wg.Wait()
					Convey("返回的结果才会关闭。", func() {
						So(begin, ShouldHappenBefore, end)
					})
				})
			})

			done = make(chan struct{})
			dones = append(dones, done)
		}
	})

}
