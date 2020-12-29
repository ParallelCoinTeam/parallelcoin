package p9

import (
	"image"
	
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

type Enum struct {
	th      *Theme
	value   string
	changed bool
	clicks  []gesture.Click
	values  []string
	hook    func(value string)
}

func (th *Theme) Enum() *Enum {
	return &Enum{th: th, hook: func(string) {}}
}

func (e *Enum) Value() string {
	return e.value
}

func (e *Enum) SetValue(value string) *Enum {
	e.value = value
	return e
}

func (e *Enum) SetOnChange(hook func(value string)) *Enum {
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
func (e *Enum) Changed() bool {
	changed := e.changed
	e.changed = false
	return changed
}

// Fn adds the event handler for key.
func (e *Enum) Fn(gtx layout.Context, key string) layout.Dimensions {
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
				if newValue := e.values[idx]; newValue != e.value {
					e.value = newValue
					e.changed = true
					e.th.BackgroundProcessingQueue <- func() { e.hook(e.value) }
				}
			}
		}
		clk.Add(gtx.Ops)
	}
	
	return layout.Dimensions{Size: gtx.Constraints.Min}
}
