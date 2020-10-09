package p9

import (
	"gioui.org/layout"
)

type changeStateHook func(b bool)

type _bool struct {
	value       *bool
	clk         *_clickable
	changed     bool
	changeState changeStateHook
}

func (b *_bool) GetValue() bool {
	return *b.value
}

func (b *_bool) Value(value *bool) {
	b.value = value
}

func NewBool(value *bool) *_bool {
	return &_bool{
		value:       value,
		clk:         NewClickable(),
		changed:     false,
		changeState: func(b bool){},
	}
}

func (b *_bool) SetHook(fn changeStateHook) *_bool {
	b.changeState = fn
	return b
}

// Changed reports whether value has changed since the last call to Changed.
func (b *_bool) Changed() bool {
	changed := b.changed
	b.changed = false
	return changed
}

// History returns the history of presses in the buffer
func (b *_bool) History() []press {
	return b.clk.History()
}

// Fn renders the events of the boolean widget
func (b *_bool) Fn(gtx layout.Context) layout.Dimensions {
	dims := b.clk.Fn(gtx)
	for b.clk.Clicked() {
		*b.value = !*b.value
		b.changed = true
		b.changeState(*b.value)
	}
	return dims
}
