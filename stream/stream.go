package stream

// Or return a DONE channel
// The DONE channel will be closed if any one of dones is closed.
func Or(dones ...<-chan struct{}) <-chan struct{} {
	if len(dones) == 0 {
		panic("Or 没有输入参数")
	}
	return or(dones)
}

func or(dones []<-chan struct{}) <-chan struct{} {
	if len(dones) == 1 {
		return dones[0]
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-dones[0]:
		case <-or(dones[1:]):
		}
	}()
	return done
}

// Repeat will repeat call fn() and send result to returned channel
func Repeat(
	done <-chan struct{},
	fn func() interface{},
) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

// FanOut make multi-workers to parallely do the work
// NOTICE: FanOut 分配工作的时候，很有可能会打乱 stream 中工作的顺序
func FanOut(
	done <-chan struct{},
	worker func(<-chan struct{}, <-chan interface{}) <-chan interface{},
	stream <-chan interface{},
	num int,
) []<-chan interface{} {
	res := make([]<-chan interface{}, 0, num)
	for i := 0; i < num; i++ {
		res = append(res, worker(done, stream))
	}
	return res
}
