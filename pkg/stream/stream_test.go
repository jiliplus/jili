package stream

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/aQuaYi/stub"
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
		So(time.Now(), ShouldHappenWithin, 120*time.Millisecond, begin)
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

// streamMaker 返回的 channel 会输出 [0,count) 的数据流
var streamMaker = func() <-chan interface{} {
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
	stream := streamMaker()
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
				stream = streamMaker()
			})
		}
	})
}

func TestFanIn(t *testing.T) {
	stream := streamMaker()
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
				stream = streamMaker()
			})
		}
	})
}

func TestOrDone(t *testing.T) {
	Convey("存在一个输出流 stream，", t, func() {
		stream := streamMaker()
		Convey("利用 OrDone 可以按照顺序读取 stream 中的内容。", func() {
			expected := 0
			var val interface{}
			for val = range OrDone(nil, stream) {
				So(val, ShouldResemble, expected)
				expected++
			}
			So(expected, ShouldEqual, count)
		})
	})

	Convey("当输入流 stream 阻塞时，", t, func() {
		done := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(1)
		//
		defer stub.Var(&orDoneStub1, func() {
			wg.Done()
		}).Restore()
		//
		resStream := OrDone(done, nil)
		Convey("输出流应该是阻塞的。", func() {
			isBlocked := false
			select {
			case <-resStream:
			default:
				isBlocked = true
			}
			So(isBlocked, ShouldBeTrue)
		})
		Convey("可以通过 done 来抢占", func() {
			close(done)
			wg.Wait()
			//
			isBlocked, isClosed := false, false
			select {
			case _, ok := <-resStream:
				if !ok {
					isClosed = true
				}
			default:
				isBlocked = true
			}
			So(isBlocked, ShouldBeFalse)
			So(isClosed, ShouldBeTrue)
		})
	})

	Convey("当输出流 resStream 阻塞时，", t, func() {
		done := make(chan struct{})
		stream := streamMaker()
		//
		var wg2, wg3 sync.WaitGroup
		hasReceived := false
		hasCancelled := false
		wg2.Add(1)
		wg3.Add(1)
		stub2 := stub.Var(&orDoneStub2, func(ok bool) {
			hasReceived = ok
			wg2.Done()
		})
		stub3 := stub.Var(&orDoneStub3, func() {
			hasCancelled = true
			wg3.Done()
		})
		defer stub2.Restore()
		defer stub3.Restore()
		//
		resStream := OrDone(done, stream)
		Convey("可以通过 done 来抢占", func() {
			wg2.Wait()
			close(done)
			wg3.Wait()
			//
			isBlocked, isClosed := false, false
			select {
			case _, ok := <-resStream:
				if !ok {
					isClosed = true
				}
			default:
				isBlocked = true
			}
			So(hasReceived, ShouldBeTrue)
			So(hasCancelled, ShouldBeTrue)
			So(isBlocked, ShouldBeFalse)
			So(isClosed, ShouldBeTrue)
			So(<-stream, ShouldEqual, 1)
		})
	})
}

func TestDuplicate(t *testing.T) {
	Convey("如果存在一个 stream，", t, func() {
		stream := streamMaker()
		Convey("复制以后，得到的两个通道，能够收到一样的内容。", func() {
			out1, out2 := Duplicate(nil, stream)
			for v1 := range out1 {
				v2 := <-out2
				So(v1, ShouldResemble, v2)
			}
		})
	})
}

func TestBridge(t *testing.T) {
	Convey("如果 genVals()，会返回 <-chan <-chan interface{}", t, func() {
		count := 10
		genVals := func() <-chan <-chan interface{} {
			resStream := make(chan (<-chan interface{}))
			go func() {
				defer close(resStream)
				for i := 0; i < count; i++ {
					// FIXME: 把 stream 改成带缓冲的话，会出现 Data Race
					stream := make(chan interface{})
					resStream <- stream
					stream <- i
					close(stream)
				}
			}()
			return resStream
		}

		Convey("Bridge 就可以一次把他们读取出来。", func() {
			i := 0
			for v := range Bridge(nil, genVals()) {
				So(v, ShouldEqual, i)
				i++
			}
			So(i, ShouldEqual, count)
		})
	})

}
