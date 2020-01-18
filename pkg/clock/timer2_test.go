package clock

import (
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO: 搜索锁的 test 文件夹中的 s.now 不能直接访问 s 的属性
// 会有数据竞争问题

func Test_Simulator_Sleep(t *testing.T) {
	Convey("新建一个 Simulator s", t, func() {
		now := time.Now()
		duration := time.Second
		s := NewSimulator(now)
		Convey("并发地 s.Sleep(duration) 与 s.Move()", func() {
			go func() {
				// NOTICE: time.Sleep 是为了保证后面的程序被阻塞到 s.Sleep
				// 但也不是 100% 的保证
				time.Sleep(duration)
				s.Move()
			}()
			s.Sleep(duration)
			// 其实运行到这一步，已经说明了，可以从 Sleep 中恢复。
			// 打到了 Sleep 的目的
			wakeUp := s.Now()
			Convey("Sleep 的完成时间，就是当前时间", func() {
				So(wakeUp, ShouldEqual, s.Now())
			})
		})
	})
}

func Test_Simulator_After(t *testing.T) {
	Convey("新建一个 Simulator s", t, func() {
		now := time.Now()
		duration := time.Second
		s := NewSimulator(now)
		expectTime := now.Add(duration)
		c := s.After(duration)
		for d := time.Duration(2); d < 5; d++ {
			dur := d * time.Second
			Convey(fmt.Sprintf("把 s 设置为 %s 后", dur), func() {
				s.AddOrPanic(dur)
				actualTime := <-c
				Convey("返回的时间，还是应该时 expectTime", func() {
					So(actualTime, ShouldEqual, expectTime)
				})
			})
		}
	})
}

func Test_Simulator_AfterFunc(t *testing.T) {
	Convey("新建一个 Simulator s", t, func() {
		now := time.Now()
		duration := time.Second
		s := NewSimulator(now)
		count := 0
		// 由于 AfterFunc 是并发执行的，需要进行同步
		// 同步的方法是
		// 在 timer 的 set 前，添加 wg.Add(1)
		// 在执行过的 timer.Reset 前，添加 wg.Add(1)
		// 在 s.Add 后添加 wg.Wait()
		var wg sync.WaitGroup
		wg.Add(1)
		timer := s.AfterFunc(duration, func() {
			count++
			wg.Done()
		})

		Convey("到达预定时间前", func() {
			s.Add(duration / 2)
			Convey("Func 不会被触发", func() {
				So(count, ShouldEqual, 0)
			})
			Convey("reset 时，还是活的", func() {
				isActive := timer.Reset(duration)
				So(isActive, ShouldBeTrue)
			})
			Convey("stop 时，还是活的", func() {
				isActive := timer.Stop()
				So(isActive, ShouldBeTrue)
			})
		})

		Convey("到达预定时间后", func() {
			s.Add(duration)
			wg.Wait()
			Convey("Func 会被触发一次", func() {
				So(count, ShouldEqual, 1)
			})
			Convey("过段时间，重置一下", func() {
				s.Add(duration)
				// 已经执行过了，所以需要 wg.Add(1)
				wg.Add(1)
				isActive := timer.Reset(duration)
				So(isActive, ShouldBeFalse)
				Convey("到达预定时间后", func() {
					s.Add(duration)
					wg.Wait()
					Convey("Func 会第 2 次被触发", func() {
						So(count, ShouldEqual, 2)
					})
				})
			})
		})

	})
}
