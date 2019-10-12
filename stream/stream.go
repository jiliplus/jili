package jili

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
