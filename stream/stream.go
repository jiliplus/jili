package stream

import "sync"

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

// FanIn make multi-channels to one
// NOTICE: FanIn 合并作的时候，很有可能会打乱 stream 中工作的顺序
func FanIn(
	done <-chan struct{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup

	resStream := make(chan interface{})

	// 采集器负责把通道中收到的值，放入 resStream
	collector := func(c <-chan interface{}) {
		defer wg.Done()
		for x := range c {
			select {
			case <-done:
				return
			case resStream <- x:
			}
		}
	}

	wg.Add(len(channels))
	// 每个 channel 都有一个专门的采集器
	for _, c := range channels {
		go collector(c)
	}

	// 所有的 channel 手机完毕后，关闭输出通道
	go func() {
		wg.Wait()
		close(resStream)
	}()

	return resStream
}

var (
	orDoneStub1 = func() {}
	orDoneStub2 = func(bool) {}
	orDoneStub3 = func() {}
	orDoneStub4 = func() {}
)

// OrDone allows the process of reading data from the stream to be interrupted by done.
func OrDone(done <-chan struct{}, stream <-chan interface{}) <-chan interface{} {
	resStream := make(chan interface{})
	go func() {
		defer close(resStream)
		for {
			select {
			case <-done: // 读取可抢占
				orDoneStub1()
				return
			case val, ok := <-stream:
				orDoneStub2(ok)
				if !ok {
					return
				}
				select {
				case <-done: // 写入可抢占
					orDoneStub3()
					return
				case resStream <- val:
					orDoneStub4()
				}
			}
		}
	}()
	return resStream
}
