package bucket

import (
	"testing"
	"time"

	"github.com/prashantv/gostub"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Bucket_Hurry(t *testing.T) {
	Convey("新生成一个符合 Bucket 接口的变量", t, func() {
		nowTime := time.Now()
		var sleepDur time.Duration
		stub := gostub.StubFunc(&now, nowTime)
		stub.Stub(&sleep, func(d time.Duration) {
			sleepDur = d
		})
		defer stub.Reset()
		b := New(time.Second*4, 4, 2)
		Convey("b.Hurry(1) 后，不会 sleep", func() {
			b.Hurry(1)
			So(sleepDur, ShouldEqual, 0)
		})
		Convey("b.Hurry(2) 后，不会 sleep", func() {
			b.Hurry(2)
			So(sleepDur, ShouldEqual, 0)
		})
		Convey("b.Hurry(3) 后，不会 sleep", func() {
			b.Hurry(3)
			So(sleepDur, ShouldEqual, 0)
		})
		Convey("b.Hurry(4) 后，不会 sleep", func() {
			b.Hurry(4)
			So(sleepDur, ShouldEqual, 0)
		})
		Convey("b.Hurry(5) 后，sleep 1 秒", func() {
			b.Hurry(5)
			So(sleepDur, ShouldEqual, time.Second)
		})
	})
}

func Test_Bucket_Wait(t *testing.T) {
	Convey("新生成一个符合 Bucket 接口的变量", t, func() {
		nowTime := time.Now()
		var sleepDur time.Duration
		stub := gostub.StubFunc(&now, nowTime)
		stub.Stub(&sleep, func(d time.Duration) {
			sleepDur = d
		})
		defer stub.Reset()
		b := New(time.Second*4, 4, 2)
		Convey("b.Wait(1) 后，不会 sleep", func() {
			b.Wait(1)
			So(sleepDur, ShouldEqual, 0)
		})
		Convey("b.Wait(2) 后，不会 sleep", func() {
			b.Wait(2)
			So(sleepDur, ShouldEqual, 0)
		})
		Convey("b.Wait(3) 后，sleep 1 秒", func() {
			b.Wait(3)
			So(sleepDur, ShouldEqual, time.Second)
		})
	})
}

func Test_Bucket_quickReturn(t *testing.T) {
	Convey("新生成一个符合 Bucket 接口的变量", t, func() {
		b := New(time.Second, 2, 1)
		Convey("如果 Hurry 的 count 不是正数。", func() {
			isQuickReturned := false
			stubs := gostub.Stub(&hurryQuickReturn, func() {
				isQuickReturned = true
			})
			defer stubs.Reset()
			Convey("b.Hurry(-1) 前，isQuickReturn == false", func() {
				So(isQuickReturned, ShouldBeFalse)
			})
			b.Hurry(-1)
			Convey("b.Hurry(-1) 后，isQuickReturn == true", func() {
				So(isQuickReturned, ShouldBeTrue)
			})
		})
		Convey("如果 Wait 的 count 不是正数。", func() {
			isQuickReturned := false
			stubs := gostub.Stub(&waitQuickReturn, func() {
				isQuickReturned = true
			})
			defer stubs.Reset()
			Convey("b.Wait(-1) 前，isQuickReturn == false", func() {
				So(isQuickReturned, ShouldBeFalse)
			})
			b.Wait(-1)
			Convey("b.Wait(-1) 后，isQuickReturn == true", func() {
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

func Test_bucket_new(t *testing.T) {
	Convey("想要利用 reserved 和 capacity 新创建一个 bucket", t, func() {
		Convey("当 reserved < 0", func() {
			reserved := int64(-1)
			Convey("无论 capacity 为何值, 会 panic", func() {
				capacity := int64(2)
				So(func() {
					newBucket(time.Second, capacity, reserved)
				},
					ShouldPanicWith,
					"bucket's parameter should 0 <= reserved < capacity")
			})
		})
		Convey("当 reserved = 0", func() {
			reserved := int64(0)
			Convey("reserved < capacity 时, 不会 panic", func() {
				capacity := reserved + 1
				So(func() {
					newBucket(time.Second, capacity, reserved)
				}, ShouldNotPanic)
			})
			Convey("reserved >= capacity 时, 会 panic", func() {
				capacity := reserved - 1
				So(func() {
					newBucket(time.Second, capacity, reserved)
				},
					ShouldPanicWith,
					"bucket's parameter should 0 <= reserved < capacity",
				)
				capacity = reserved
				So(func() {
					newBucket(time.Second, capacity, reserved)
				},
					ShouldPanicWith,
					"bucket's parameter should 0 <= reserved < capacity",
				)
			})
		})
		Convey("当 reserved > 0", func() {
			reserved := int64(1)
			Convey("reserved < capacity 时, 不会 panic", func() {
				capacity := reserved + 1
				So(func() {
					newBucket(time.Second, capacity, reserved)
				}, ShouldNotPanic)
			})
			Convey("reserved >= capacity 时, 会 panic", func() {
				capacity := reserved - 1
				So(func() {
					newBucket(time.Second, capacity, reserved)
				},
					ShouldPanicWith,
					"bucket's parameter should 0 <= reserved < capacity",
				)
				capacity = reserved
				So(func() {
					newBucket(time.Second, capacity, reserved)
				},
					ShouldPanicWith,
					"bucket's parameter should 0 <= reserved < capacity",
				)
			})
		})
	})
	Convey("在新创建了一个 bucket 后", t, func() {
		dur := time.Second
		capacity, reserved := int64(2), int64(1)
		nowTime := time.Now()
		stub := gostub.Stub(&now, func() time.Time {
			return nowTime
		})
		defer stub.Reset()
		b := newBucket(dur, capacity, reserved)
		d := gcd(capacity, int64(dur))
		Convey("bucket 的各项属性应该符合预期。", func() {
			So(b.start, ShouldEqual, nowTime)
			So(b.reserved, ShouldEqual, reserved)
			So(b.normal, ShouldEqual, capacity-reserved)
			So(b.interval, ShouldEqual, dur/time.Duration(d))
			So(b.quantum, ShouldEqual, capacity/d)
			So(b.tick, ShouldEqual, 0)
			So(b.hToken, ShouldEqual, reserved)
			So(b.wToken, ShouldEqual, capacity-reserved)
		})
	})
}

func Test_bucket_update(t *testing.T) {
	Convey("把新 bucket 的 hToken 和 wToken 清空", t, func() {
		capacity, reserved := int64(4), int64(2)
		b := newBucket(time.Second, capacity, reserved)
		b.hToken, b.wToken = 0, 0
		Convey("1 个 interval 后", func() {
			b.update(b.start.Add(b.interval * 1))
			Convey("hToken = 1, wToken = 0", func() {
				So(b.hToken, ShouldEqual, 1)
				So(b.wToken, ShouldEqual, 0)
			})
		})
		Convey("2 个 interval 后", func() {
			b.update(b.start.Add(b.interval * 2))
			Convey("hToken = 2, wToken = 0", func() {
				So(b.hToken, ShouldEqual, 2)
				So(b.wToken, ShouldEqual, 0)
			})
		})
		Convey("3 个 interval 后", func() {
			b.update(b.start.Add(b.interval * 3))
			Convey("hToken = 2, wToken = 1", func() {
				So(b.hToken, ShouldEqual, 2)
				So(b.wToken, ShouldEqual, 1)
			})
		})
		Convey("4 个 interval 后", func() {
			b.update(b.start.Add(b.interval * 4))
			Convey("hToken = 2, wToken = 2", func() {
				So(b.hToken, ShouldEqual, 2)
				So(b.wToken, ShouldEqual, 2)
			})
		})
		Convey("5 个 interval 后", func() {
			b.update(b.start.Add(b.interval * 5))
			Convey("hToken = 2, wToken = 2", func() {
				So(b.hToken, ShouldEqual, 2)
				So(b.wToken, ShouldEqual, 2)
			})
		})
	})
}

func Test_bucket_hTake(t *testing.T) {
	Convey("新建了一个 bucket", t, func() {
		capacity, reserved := int64(4), int64(2)
		b := newBucket(time.Second, capacity, reserved)
		Convey("hTake(1)", func() {
			remain := b.hTake(1)
			Convey("remain = 0, hToken = 1, wToken = 2", func() {
				So(remain, ShouldEqual, 0)
				So(b.hToken, ShouldEqual, 1)
				So(b.wToken, ShouldEqual, 2)
			})
		})
		Convey("hTake(2)", func() {
			remain := b.hTake(2)
			Convey("remain = 0, hToken = 0, wToken = 2", func() {
				So(remain, ShouldEqual, 0)
				So(b.hToken, ShouldEqual, 0)
				So(b.wToken, ShouldEqual, 2)
			})
		})
		Convey("hTake(3)", func() {
			remain := b.hTake(3)
			Convey("remain = 0, hToken = 0, wToken = 1", func() {
				So(remain, ShouldEqual, 0)
				So(b.hToken, ShouldEqual, 0)
				So(b.wToken, ShouldEqual, 1)
			})
		})
		Convey("hTake(4)", func() {
			remain := b.hTake(4)
			Convey("remain = 0, hToken = 0, wToken = 0", func() {
				So(remain, ShouldEqual, 0)
				So(b.hToken, ShouldEqual, 0)
				So(b.wToken, ShouldEqual, 0)
			})
		})
		Convey("hTake(5)", func() {
			remain := b.hTake(5)
			Convey("remain = 1, hToken = 0, wToken = -1", func() {
				So(remain, ShouldEqual, 1)
				So(b.hToken, ShouldEqual, 0)
				So(b.wToken, ShouldEqual, -1)
			})
		})
	})
}

func Test_bucket_wTake(t *testing.T) {
	Convey("新建了一个 bucket", t, func() {
		capacity, reserved := int64(4), int64(2)
		b := newBucket(time.Second, capacity, reserved)
		Convey("wTake(1)", func() {
			remain := b.wTake(1)
			Convey("remain = 0, hToken = 2, wToken = 1", func() {
				So(remain, ShouldEqual, 0)
				So(b.hToken, ShouldEqual, 2)
				So(b.wToken, ShouldEqual, 1)
			})
		})
		Convey("wTake(2)", func() {
			remain := b.wTake(2)
			Convey("remain = 0, hToken = 2, wToken = 0", func() {
				So(remain, ShouldEqual, 0)
				So(b.hToken, ShouldEqual, 2)
				So(b.wToken, ShouldEqual, 0)
			})
		})
		Convey("wTake(3)", func() {
			remain := b.wTake(3)
			Convey("remain = 1, hToken = 2, wToken = -1", func() {
				So(remain, ShouldEqual, 1)
				So(b.hToken, ShouldEqual, 2)
				So(b.wToken, ShouldEqual, -1)
			})
		})
	})
}

func Test_bucket_needTime(t *testing.T) {
	Convey("新建了一个 bucket", t, func() {
		capacity, reserved := int64(4), int64(2)
		b := newBucket(time.Second, capacity, reserved)
		Convey("needTime(0)", func() {
			dur := b.needTime(0, now())
			Convey("dur = 0", func() {
				So(dur, ShouldEqual, 0)
			})
		})
		Convey("needTime(4)", func() {
			dur := b.needTime(4, b.tick2Time())
			Convey("dur = time.Second", func() {
				So(dur, ShouldEqual, time.Second)
			})
		})
		Convey("needTime(8)", func() {
			dur := b.needTime(8, b.tick2Time())
			Convey("dur = time.Second", func() {
				So(dur, ShouldEqual, 2*time.Second)
			})
		})
	})
}
