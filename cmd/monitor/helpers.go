package monitor

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

func (m *State) FlexV(children ...layout.FlexChild) () {
	layout.Flex{Axis: layout.Vertical}.Layout(m.Gtx, children...)
}

func (m *State) FlexH(children ...layout.FlexChild) {
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

func Spacer() layout.FlexChild {
	return Flexed(1, func() {})
}

func (m *State) Rectangle(width, height int, color string) {
	gelook.DuoUIdrawRectangle(m.Gtx,
		width, height, m.Theme.Colors[color],
		[4]float32{0, 0, 0, 0},
		[4]float32{0, 0, 0, 0},
	)
}

func (m *State) IconButton(icon, fg, bg string, button *gel.Button) {
	m.Theme.DuoUIbutton("", "", "",
		m.Theme.Colors[bg], "", m.Theme.Colors[fg], icon,
		m.Theme.Colors[fg], 0, 32, 41, 41,
		0, 0).IconLayout(m.Gtx, button)
}

func (m *State) TextButton(label, fontFace string, fontSize int, fg, bg string,
	button *gel.Button) {
	m.Theme.DuoUIbutton(
		m.Theme.Fonts[fontFace],
		label,
		m.Theme.Colors[fg],
		m.Theme.Colors[bg],
		m.Theme.Colors[bg],
		m.Theme.Colors[fg],
		"settingsIcon",
		m.Theme.Colors["Light"],
		fontSize, 0, 80, 32, 4, 4).
		Layout(m.Gtx, button)
}

func Toggle(b *bool) bool {
	*b = !*b
	return *b
}
