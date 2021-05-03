package gel

import (
	"image"

	"github.com/p9c/gio/gesture"
	"github.com/p9c/gio/io/pointer"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
)

type Enum struct {
	value   string
	changed bool
	clicks  []gesture.Click
	values  []string
	hook    func(value string)
}

func (w *Window) Enum() *Enum {
	return &Enum{hook: func(string) {}}
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
func (e *Enum) Fn(gtx l.Context, key string) l.Dimensions {
	defer op.Save(gtx.Ops).Load()
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
				if niew := e.values[idx]; niew != e.value {
					e.value = niew
					e.changed = true
					e.hook(e.value)
				}
			}
		}
		clk.Add(gtx.Ops)
	}

	return l.Dimensions{Size: gtx.Constraints.Min}
}
