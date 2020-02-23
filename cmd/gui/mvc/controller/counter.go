package controller

type DuoUIcounter struct {
	Value        int
	OperateValue int
}

func (c *DuoUIcounter) Increase() {
	c.Value = c.Value + c.OperateValue
}

func (c *DuoUIcounter) Decrease() {
	c.Value = c.Value - c.OperateValue
}
func (c *DuoUIcounter) Reset() {
	c.Value = 0
}
