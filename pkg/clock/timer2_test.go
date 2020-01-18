package clock

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

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
		expectTime := now.Add(duration)
		isFired := false
		// TODO: 删除此处内容
		c := s.AfterFunc(duration, func() {

		})
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
