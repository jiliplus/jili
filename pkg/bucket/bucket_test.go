package bucket

import (
	"testing"
	"time"

	"github.com/prashantv/gostub"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Bucket_quickReturn(t *testing.T) {
	Convey("新生成一个符合 Bucket 接口的变量", t, func() {
		b := New(time.Second, 2, 1)
		Convey("如果 Hurry 的 count 不是正数。", func() {
			isQuickReturned := false
			stubs := gostub.Stub(&hurryQuickReturn, func() {
				isQuickReturned = true
			})
			defer stubs.Reset()
			Convey("b.Hurry(-1) 后，isQuickReturn == true", func() {
				Convey("没有使用 Bucket 时，isQuickReturn == false", func() {
					So(isQuickReturned, ShouldBeFalse)
				})
				b.Hurry(-1)
				So(isQuickReturned, ShouldBeTrue)
			})
			Convey("b.Hurry(0) 后，isQuickReturn == true", func() {
				Convey("没有使用 Bucket 时，isQuickReturn == false", func() {
					So(isQuickReturned, ShouldBeFalse)
				})
				b.Hurry(0)
				So(isQuickReturned, ShouldBeTrue)
			})
		})
		Convey("如果 Wait 的 count 不是正数。", func() {
			isQuickReturned := false
			stubs := gostub.Stub(&waitQuickReturn, func() {
				isQuickReturned = true
			})
			defer stubs.Reset()
			Convey("b.Wait(-1) 后，isQuickReturn == true", func() {
				Convey("没有使用 Bucket 时，isQuickReturn == false", func() {
					So(isQuickReturned, ShouldBeFalse)
				})
				b.Wait(-1)
				So(isQuickReturned, ShouldBeTrue)
			})
			Convey("b.Wait(0) 后，isQuickReturn == true", func() {
				Convey("没有使用 Bucket 时，isQuickReturn == false", func() {
					So(isQuickReturned, ShouldBeFalse)
				})
				b.Wait(0)
				So(isQuickReturned, ShouldBeTrue)
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
