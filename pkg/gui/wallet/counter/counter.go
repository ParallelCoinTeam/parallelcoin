package counter

import (
	"gioui.org/widget"
)

type Counter struct {
	Value           int
	OperateValue    int
	From            int
	To              int
	PageFunction    func()
	CounterInput    *widget.Editor
	CounterIncrease *widget.Clickable
	CounterDecrease *widget.Clickable
	CounterReset    *widget.Clickable
}

func (c *Counter) Increase() {
	if c.Value < c.To {
		c.Value = c.Value + c.OperateValue
	}
}

func (c *Counter) Decrease() {
	if c.Value > c.From {
		c.Value = c.Value - c.OperateValue
	}
	if c.Value < 0 {
		c.Value = 0
	}
}
func (c *Counter) Reset() {
	c.Value = 0
}
