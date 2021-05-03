package gel

import (
	"fmt"

	l "github.com/p9c/gio/layout"
)

type IntSlider struct {
	*Window
	min, max   *Clickable
	floater    *Float
	hook       func(int)
	value      int
	minV, maxV float32
}

func (w *Window) IntSlider() *IntSlider {
	return &IntSlider{
		Window:  w,
		min:     w.Clickable(),
		max:     w.Clickable(),
		floater: w.Float(),
		hook:    func(int) {},
	}
}

func (i *IntSlider) Min(min float32) *IntSlider {
	i.minV = min
	return i
}

func (i *IntSlider) Max(max float32) *IntSlider {
	i.maxV = max
	return i
}

func (i *IntSlider) Hook(fn func(v int)) *IntSlider {
	i.hook = fn
	return i
}

func (i *IntSlider) Value(value int) *IntSlider {
	i.value = value
	i.floater.SetValue(float32(value))
	return i
}

func (i *IntSlider) GetValue() int {
	return int(i.floater.Value() + 0.5)
}

func (i *IntSlider) Fn(gtx l.Context) l.Dimensions {
	return i.Flex().Rigid(
		i.Button(
			i.min.SetClick(func() {
				i.floater.SetValue(i.minV)
				i.hook(int(i.minV))
			})).
			Inset(0.25).
			Color("Primary").
			Background("Transparent").
			Font("bariol regular").
			Text("0").
			Fn,
	).Flexed(1,
		i.Inset(0.25,
			i.Slider().
				Float(i.floater.SetHook(func(fl float32) {
					iFl := int(fl + 0.5)
					i.value = iFl
					i.floater.SetValue(float32(iFl))
					i.hook(iFl)
				})).
				Min(i.minV).Max(i.maxV).
				Fn,
		).Fn,
	).Rigid(
		i.Button(
			i.max.SetClick(func() {
				i.floater.SetValue(i.maxV)
				i.hook(int(i.maxV))
			})).
			Inset(0.25).
			Color("Primary").
			Background("Transparent").
			Font("bariol regular").
			Text(fmt.Sprint(int(i.maxV))).
			Fn,
	).Fn(gtx)
}
