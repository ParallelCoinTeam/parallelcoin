package gel

type DuoUIcounter struct {
	Value           int
	OperateValue    int
	From            int
	To              int
	CounterInput    *Editor
	CounterIncrease *Button
	CounterDecrease *Button
	CounterReset    *Button
}

func (c *DuoUIcounter) Increase() {
	if c.Value < c.To {
		c.Value = c.Value + c.OperateValue
	}
}

func (c *DuoUIcounter) Decrease() {
	if c.Value > c.From {
		c.Value = c.Value - c.OperateValue
	}
	if c.Value < 0 {
		c.Value = 0
	}
}
func (c *DuoUIcounter) Reset() {
	c.Value = 0
}
