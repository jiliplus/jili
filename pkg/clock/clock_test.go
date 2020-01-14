package clock

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type suit func(c Clock)

func realAndMockClock(t *testing.T, test suit) {
	Convey("测试 real clock", t, func() {
		c := NewRealClock()
		test(c)
	})
	Convey("测试 mock clock", t, func() {
		now := time.Now()
		c := NewMockClock(now)
		test(c)
	})
}
