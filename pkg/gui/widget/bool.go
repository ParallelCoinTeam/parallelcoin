package widget

import (
	"gioui.org/layout"
)

type changeStateHook func(b bool)

type Bool struct {
	value       *bool
	clk         *Clickable
	changed     bool
	changeState changeStateHook
}

func (b *Bool) GetValue() *bool {
	return b.value
}

func (b *Bool) Value(value *bool) {
	b.value = value
}

func NewBool(value *bool) *Bool {
	return &Bool{
		value:       value,
		clk:         NewClickable(),
		changed:     false,
		changeState: func(b bool){},
	}
}

func (b *Bool) SetHook(fn changeStateHook) *Bool {
	b.changeState = fn
	return b
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
	for b.clk.Clicked() {
		*b.value = !*b.value
		b.changed = true
		b.changeState(*b.value)
	}
	return dims
}
