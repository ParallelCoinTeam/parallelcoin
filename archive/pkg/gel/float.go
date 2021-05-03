// SPDX-License-Identifier: Unlicense OR MIT

package gel

import (
	"image"

	"github.com/p9c/gio/gesture"
	"github.com/p9c/gio/io/pointer"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
)

// Float is for selecting a value in a range.
type Float struct {
	value      float32
	drag       gesture.Drag
	pos        float32 // position normalized to [0, 1]
	length     float32
	changed    bool
	changeHook func(float32)
}

func (th *Theme) Float() *Float {
	return &Float{changeHook: func(float32) {}}
}

func (f *Float) SetValue(value float32) *Float {
	f.value = value
	return f
}
func (f *Float) Value() float32 {
	return f.value
}

func (f *Float) SetHook(fn func(fl float32)) *Float {
	f.changeHook = fn
	return f
}

// Fn processes events.
func (f *Float) Fn(gtx l.Context, pointerMargin int, min, max float32) l.Dimensions {
	size := gtx.Constraints.Min
	f.length = float32(size.X)
	var de *pointer.Event
	for _, ev := range f.drag.Events(gtx.Metric, gtx, gesture.Horizontal) {
		if ev.Type == pointer.Press || ev.Type == pointer.Drag {
			de = &ev
		}
		if ev.Type == pointer.Release {
			f.changeHook(f.value)
		}
	}
	value := f.value
	if de != nil {
		f.pos = de.Position.X / f.length
		value = min + (max-min)*f.pos
	} else if min != max {
		f.pos = value/(max-min) - min
	}
	// Unconditionally call setValue in case min, max, or value changed.
	f.setValue(value, min, max)

	if f.pos < 0 {
		f.pos = 0
	} else if f.pos > 1 {
		f.pos = 1
	}

	defer op.Save(gtx.Ops).Load()
	rect := image.Rectangle{Max: size}
	rect.Min.X -= pointerMargin
	rect.Max.X += pointerMargin
	pointer.Rect(rect).Add(gtx.Ops)
	f.drag.Add(gtx.Ops)

	return l.Dimensions{Size: size}
}

func (f *Float) setValue(value, min, max float32) {
	if min > max {
		min, max = max, min
	}
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	if f.value != value {
		f.value = value
		f.changed = true
	}
}

// Pos reports the selected position.
func (f *Float) Pos() float32 {
	return f.pos * f.length
}

// Changed reports whether the value has changed since the last call to Changed.
func (f *Float) Changed() bool {
	changed := f.changed
	f.changed = false
	return changed
}
