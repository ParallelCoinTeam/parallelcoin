package sem

type semaphore struct{}

type T chan semaphore

func (t T) Ready() {
	t <- semaphore{}
}

func (t T) Yield() {
	<-t
}
