package bucket

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_take(t *testing.T) {
	Convey("想要从 bucket 中拿走 token", t, func() {
		b := newBucket(time.Minute, 60)
		Convey("如果 count <=0,", func() {
			waitTime0 := b.take(time.Now(), 0)
			waitTime1 := b.take(time.Now(), -1)
			Convey("那么，不需要等待", func() {
				So(waitTime0, ShouldEqual, 0)
				So(waitTime1, ShouldEqual, 0)
			})
		})
		Convey("如果 count = b.available + 1", func() {
			waitTime := b.take(time.Now(), b.available+1)
			Convey("那么，需要等待将近一秒钟", func() {
				So(waitTime, ShouldBeBetweenOrEqual, time.Millisecond*980, time.Second)
			})
		})
	})
}

func Test_newBucket(t *testing.T) {
	Convey("想要生成 *bucket", t, func() {
		Convey("newBucket(time.Minute, 60)", func() {
			b := newBucket(time.Minute, 60)
			Convey("b.interval 应该是 一秒钟", func() {
				interval := time.Duration(b.interval)
				So(interval, ShouldEqual, time.Second)
			})
			Convey("b.quantum 应该是 1", func() {
				So(b.quantum, ShouldEqual, 1)
			})
		})
	})
}

func Test_gcd(t *testing.T) {
	Convey("想要求得 gcd", t, func() {
		Convey("gcd(8,2)", func() {
			So(gcd(8, 2), ShouldEqual, 2)
		})
		Convey("gcd(3,9)", func() {
			So(gcd(3, 9), ShouldEqual, 3)
		})
		Convey("gcd(9,3)", func() {
			So(gcd(9, 3), ShouldEqual, 3)
		})
		Convey("gcd(3,7)", func() {
			So(gcd(3, 7), ShouldEqual, 1)
		})
		Convey("gcd(7,3)", func() {
			So(gcd(7, 3), ShouldEqual, 1)
		})
	})
}
