package p9

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

type _enum struct {
	value   string
	changed bool
	clicks  []gesture.Click
	values  []string
	hook    func(value string)
}

func (th *Theme) Enum() *_enum {
	return &_enum{hook: func(string) {}}
}

func (e *_enum) Value() string {
	return e.value
}

func (e *_enum) SetValue(value string) *_enum {
	e.value = value
	return e
}

func (e *_enum) SetOnChange(hook func(value string)) *_enum {
	e.hook = hook
	return e
}

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Changed reports whether Value has changed by user interaction since the last call to Changed.
func (e *_enum) Changed() bool {
	changed := e.changed
	e.changed = false
	return changed
}

// Fn adds the event handler for key.
func (e *_enum) Fn(gtx layout.Context, key string) layout.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)

	if index(e.values, key) == -1 {
		e.values = append(e.values, key)
		e.clicks = append(e.clicks, gesture.Click{})
		e.clicks[len(e.clicks)-1].Add(gtx.Ops)
	} else {
		idx := index(e.values, key)
		clk := &e.clicks[idx]
		for _, ev := range clk.Events(gtx) {
			switch ev.Type {
			case gesture.TypeClick:
				if new := e.values[idx]; new != e.value {
					e.value = new
					e.changed = true
					e.hook(e.value)
				}
			}
		}
		clk.Add(gtx.Ops)
	}

	return layout.Dimensions{Size: gtx.Constraints.Min}
}
