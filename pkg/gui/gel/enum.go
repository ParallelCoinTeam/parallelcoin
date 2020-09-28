package gel

import (
	"gioui.org/gesture"
	"gioui.org/layout"
)

type Enum struct {
	clicks []gesture.Click
	values []string
	value  string
}

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Value processes events and returns the last selected value, or
// the empty string.
func (e *Enum) Value(gtx *layout.Context) string {
	for i := range e.clicks {
		for _, ev := range e.clicks[i].Events(gtx) {
			switch ev.Type {
			case gesture.TypeClick:
				e.value = e.values[i]
			}
		}
	}
	return e.value
}

// Layout adds the event handler for key.
func (e *Enum) Layout(gtx *layout.Context, key string) {
	if index(e.values, key) == -1 {
		e.values = append(e.values, key)
		e.clicks = append(e.clicks, gesture.Click{})
		e.clicks[len(e.clicks)-1].Add(gtx.Ops)
	} else {
		idx := index(e.values, key)
		e.clicks[idx].Add(gtx.Ops)
	}
}

func (e *Enum) SetValue(value string) {
	e.value = value
}
