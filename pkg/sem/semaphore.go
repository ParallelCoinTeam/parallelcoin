package sem

type semaphore struct{}

type T chan semaphore

// New creates a new semaphore of a given capacity.
func New(limit int) T {
	if limit < 1 {
		panic("a semaphore requires at least one buffer per worker")
	}
	t := make(T, limit)
	return t
}

// Acquire causes threads waiting on release to stop working.
// Note that because the channel is buffered,
// this can execute as many times as number of buffers,
// thus blocking subsequent callers until one of the buffers is released
func (t T) Acquire() {
	//log.L.Debug("acquiring semaphore")
	t <- semaphore{}
}

// Release can be waited on or selected to get a stop work signal.
// This essentially empties a slot in the semaphore which allows another
// thread to acquire it
func (t T) Release() T {
	//log.L.Debug("releasing semaphore")
	return t
}
