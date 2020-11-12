package p9

type Pool struct {
	th              *Theme
	bools           []*Bool
	boolsInUse      int
	lists           []*List
	listsInUse      int
	checkables      []*Checkable
	checkablesInUse int
	clickables      []*Clickable
	clickablesInUse int
	editors          []*Editor
	editorsInUse     int
	incDecs         []*IncDec
	incDecsInUse    int
}

func (th *Theme) NewPool() *Pool {
	return &Pool{th: th}
}

func (p *Pool) Reset() {
	p.boolsInUse = 0
	p.listsInUse = 0
	p.checkablesInUse = 0
	p.clickablesInUse = 0
	p.editorsInUse = 0
	p.incDecsInUse = 0
}
