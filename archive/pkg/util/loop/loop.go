package loop

type to struct {
	Start int
	End   int
}

func (t *to) Do(fn func(i int)) {
	for ; t.Start != t.End; t.Start++ {
		fn(t.Start)
	}
}

func To(end int, fn func(i int)) {
	(&to{End: end}).Do(fn)
}
