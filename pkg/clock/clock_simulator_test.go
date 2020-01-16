package clock

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Simulator_setNowTo(t *testing.T) {
	Convey("假设存在模拟器 s", t, func() {
		now := time.Now()
		s := &Simulator{
			now: now,
		}
		Convey("如果想把 s.now 设置成过去的时间点", func() {
			last := s.now
			s.setNowTo(s.now.Add(-time.Second))
			Convey("s.now 还是会等于原来的值", func() {
				So(s.now, ShouldEqual, last)
			})
		})
		Convey("如果想把 s.now 设置成以后的时间点", func() {
			last := s.now
			now := s.now.Add(time.Second)
			s.setNowTo(now)
			Convey("s.now 会被设置成新值", func() {
				So(last.Before(s.now), ShouldBeTrue)
				So(s.now, ShouldEqual, now)
			})
		})
	})
}
