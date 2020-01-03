package bucket

import (
	"sync"
	"time"
)

var now = time.Now
var sleep = time.Sleep

type bucket struct {
	// start 保存了 bucket 创建的时间
	start time.Time
	// bucket 的最大容量
	capacity int64
	// 每次添加的数量
	quantum int64
	// tick 的时长
	interval int64
	// mutex 保护以下属性
	sync.Mutex
	// available tokens
	available int64
	// latestTick for available tokens
	latestTick int64
}

// newBucket return bucket point
func newBucket(duration time.Duration, capacity int64) *bucket {
	if capacity <= 0 {
		panic("bucket's capacity should > 0")
	}
	//   rate
	// = capacity ÷ duration
	// = quantum ÷ interval
	// 由于此 bucket 的应用场景中，
	// duration 远大于 capacity 且
	// 极大概率 duration%capacity 等于 0
	// 所以，采用此方案求取 interval 和 quantum
	// 在其他场景下，还是要使用
	// https://github.com/juju/ratelimit/blob/f60b32039441cd828005f82f3a54aafd00bc9882/ratelimit.go#L104
	// 中使用的方法。
	d := gcd(int64(duration), capacity)
	interval, quantum := int64(duration)/d, capacity/d
	return &bucket{
		start:      now(),
		capacity:   capacity,
		quantum:    quantum,
		interval:   interval,
		available:  capacity,
		latestTick: 0,
	}
}

func gcd(m, n int64) int64 {
	if n == 0 {
		return m
	}
	return gcd(n, m%n)
}

func (b *bucket) take(now time.Time, count int64) (waitTime time.Duration) {
	if count <= 0 {
		return 0
	}
	tick := b.tickOf(now)
	b.updateAvailable(tick)
	remain := b.available - count
	if remain >= 0 {
		b.available = remain
		return 0
	}
	// +(b.quantum-1) 是为了到达 endTick 时，
	// 一定有足够的 token
	endTick := tick + (-remain+(b.quantum-1))/b.quantum
	endTime := b.start.Add(time.Duration(endTick * b.interval))
	waitTime = endTime.Sub(now)
	return
}

func (b *bucket) tickOf(now time.Time) int64 {
	return int64(now.Sub(b.start)) / b.interval
}

func (b *bucket) updateAvailable(newTick int64) {
	lastTick := b.latestTick
	b.latestTick = newTick
	b.available += (newTick - lastTick) * b.quantum
	if b.available > b.capacity {
		b.available = b.capacity
	}
	return
}
