package stream

import (
	"context"
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

	sig := func(t time.Duration) <-chan struct{} {
		done := make(chan struct{})
		go func() {
			defer close(done)
			time.Sleep(t)
		}()
		return done
	}

	Convey("Or 的返回值会和最先关闭的输入参数一起关闭", t, func() {
		begin := time.Now()
		<-Or(
			sig(100*time.Millisecond),
			sig(100*time.Second),
			sig(100*time.Minute),
			sig(100*time.Hour),
		)
		So(time.Now(), ShouldHappenWithin, 110*time.Millisecond, begin)
	})

}

func TestRepeat(t *testing.T) {
	Convey("如果有一个对 i 进行累加的函数", t, func() {
		i := 0
		fn := func() func() interface{} {
			return func() interface{} {
				i++
				return i
			}
		}()
		ctx, cancel := context.WithCancel(context.Background())

		Convey("放入 repeat 后", func() {

			valStream := Repeat(ctx.Done(), fn)

			Convey("收到的值，应该递增", func() {
				for i := 1; i < 10; i++ {
					val := (<-valStream).(int)
					So(val, ShouldEqual, i)
				}
			})

			Convey("关闭后，", func() {
				cancel()

				val, ok := <-valStream
				Convey("立即获取的话，由于 select 的运行机制，还是有可能获取到 1", func() {
					if ok {
						So(val, ShouldEqual, 1)
					} else {
						So(val, ShouldBeNil)
					}
				})

				val, ok = <-valStream
				Convey("再次获取，一定是获取到默认值", func() {
					So(val, ShouldBeNil)
					So(ok, ShouldBeFalse)
				})

			})
		})
	})
}

// worker 把收到的值转发出去
var worker = func(done <-chan struct{}, stream <-chan interface{}) <-chan interface{} {
	resStream := make(chan interface{})
	go func() {
		defer close(resStream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-stream:
				if ok {
					resStream <- val
				} else {
					return
				}
			}
		}
	}()
	return resStream
}

var count = 100

// stream 产生了 [0,count) 的数据流
var streamFn = func() <-chan interface{} {
	res := make(chan interface{})
	go func() {
		defer close(res)
		for i := 0; i < count; i++ {
			res <- i
		}
	}()
	return res
}

func TestFanOut(t *testing.T) {
	stream := streamFn()
	Convey("如果想要多个 worker 分担 stream 中来的工作，", t, func() {
		for num := 1; num < 12; num++ {
			Convey(fmt.Sprintf("当有 %d 个 worker 的时候", num), func() {
				record := make([]int, count)

				outs := FanOut(nil, worker, stream, num)
				for _, ch := range outs {
					for index := range ch {
						record[index.(int)]++
					}
				}

				Convey("每个记录的值，都应该是 1", func() {
					for i := 0; i < count; i++ {
						So(record[i], ShouldEqual, 1)
					}
				})
			})
			Reset(func() {
				stream = streamFn()
			})
		}
	})
}

func TestFanIn(t *testing.T) {
	stream := streamFn()
	Convey("如果想要多个 worker 分担 stream 中来的工作，", t, func() {
		for num := 1; num < 12; num++ {
			Convey(fmt.Sprintf("当有 %d 个 worker 的时候", num), func() {
				record := make([]int, count)
				outs := FanOut(nil, worker, stream, num)
				for index := range FanIn(nil, outs...) {
					record[index.(int)]++
				}
				Convey("每个记录的值，都应该是 1", func() {
					for i := 0; i < count; i++ {
						So(record[i], ShouldEqual, 1)
					}
				})
			})
			Reset(func() {
				stream = streamFn()
			})
		}
	})
}
