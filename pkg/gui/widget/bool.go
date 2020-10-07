package widget

import (
	"gioui.org/layout"
)

type Bool struct {
	value       *bool
	clk         *Clickable
	changed     bool
	changeState func(b, cs bool)
}

func (b *Bool) GetValue() *bool {
	return b.value
}

func (b *Bool) Value(value *bool) {
	b.value = value
}

func NewBool(value *bool, changeState func(b, cs bool)) *Bool {
	return &Bool{
		value:       value,
		clk:         NewClickable(),
		changed:     false,
		changeState: changeState,
	}
}

// Changed reports whether value has changed since the last call to Changed.
func (b *Bool) Changed() bool {
	changed := b.changed
	b.changed = false
	return changed
}

// History returns the history of presses in the buffer
func (b *Bool) History() []Press {
	return b.clk.History()
}

// Fn renders the events of the boolean widget
func (b *Bool) Fn(gtx layout.Context) layout.Dimensions {
	dims := b.clk.Fn(gtx)
	old := *b.value
	for b.clk.Clicked() {
		*b.value = !*b.value
		b.changed = true
	}
	// send the signal on the channel of the eventual changeState if it changed
	if b.changed {
		b.changeState(*b.value, old != *b.value)
	}
	return dims
}
