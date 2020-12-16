package qu

type C chan struct{}

func T(i ...int) C {
	if len(i) > 0 {
		return make(C, i[0])
	} else {
		return make(C)
	}
}

func (c C) Quit() {
	close(c)
}

func (c C) Wait() <-chan struct{} {
	return c
}
