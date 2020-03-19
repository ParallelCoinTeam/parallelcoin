package monitor

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

func (st *State) FlexV(children ...layout.FlexChild) () {
	layout.Flex{Axis: layout.Vertical}.Layout(st.Gtx, children...)
}

func (st *State) FlexH(children ...layout.FlexChild) {
	layout.Flex{Axis: layout.Horizontal}.Layout(st.Gtx, children...)
}

func (st *State) Inset(size int, fn func()) {
	layout.UniformInset(unit.Dp(float32(size))).Layout(st.Gtx, fn)
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

func (st *State) Rectangle(width, height int, color string) {
	gelook.DuoUIdrawRectangle(st.Gtx,
		width, height, st.Theme.Colors[color],
		[4]float32{0, 0, 0, 0},
		[4]float32{0, 0, 0, 0},
	)
}

func (st *State) IconButton(icon, fg, bg string, button *gel.Button) {
	st.Theme.DuoUIbutton("", "", "",
		st.Theme.Colors[bg], "", st.Theme.Colors[fg], icon,
		st.Theme.Colors[fg], 0, 32, 41, 41,
		0, 0).IconLayout(st.Gtx, button)
}

func (st *State) TextButton(label, fontFace string, fontSize int, fg, bg string,
	button *gel.Button) {
	st.Theme.DuoUIbutton(
		st.Theme.Fonts[fontFace],
		label,
		st.Theme.Colors[fg],
		st.Theme.Colors[bg],
		st.Theme.Colors[bg],
		st.Theme.Colors[fg],
		"settingsIcon",
		st.Theme.Colors["Light"],
		fontSize, 0, 80, 32, 4, 4).
		Layout(st.Gtx, button)
}

func (st *State) Text(txt, color, face, tag string) func() {
	return func() {
		var desc gelook.DuoUIlabel
		switch tag {
		case "body1":
			desc = st.Theme.Body1(txt)
		case "body2":
			desc = st.Theme.Body2(txt)
		case "h2":
			desc = st.Theme.H2(txt)
		case "h3":
			desc = st.Theme.H3(txt)
		case "h4":
			desc = st.Theme.H4(txt)
		case "h5":
			desc = st.Theme.H5(txt)
		case "h6":
			desc = st.Theme.H6(txt)
		}
		desc.Font.Typeface = st.Theme.Fonts[face]
		desc.Color = st.Theme.Colors[color]
		desc.Layout(st.Gtx)
	}
}

func Toggle(b *bool) bool {
	*b = !*b
	return *b
}
