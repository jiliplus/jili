package tools

import (
	"time"
)

// LocalTime 把毫秒格式的 UTC 时间，转换成当地时间
func LocalTime(UTCMillisecond int64) time.Time {
	return time.Unix(0, UTCMillisecond*1000000)
}
