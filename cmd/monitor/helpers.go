package monitor

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func (m *State) FlexVertical(children ...layout.FlexChild) () {
	layout.Flex{Axis: layout.Vertical}.Layout(m.Gtx, children...)
}

func (m *State) FlexHorizontal(children ...layout.FlexChild) {
	layout.Flex{Axis: layout.Horizontal}.Layout(m.Gtx, children...)
}

func (m *State) Inset(size int, fn func()) {
	layout.UniformInset(unit.Dp(float32(size))).Layout(m.Gtx, fn)
}

func Rigid(widget func()) layout.FlexChild {
	return layout.Rigid(widget)
}

func Flexed(weight float32, widget func()) layout.FlexChild {
	return layout.Flexed(weight, widget)
}
