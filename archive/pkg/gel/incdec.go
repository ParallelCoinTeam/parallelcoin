package gel

import (
	"fmt"
	
	l "github.com/p9c/gio/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type IncDec struct {
	*Window
	nDigits           int
	min, max          int
	amount            int
	current           int
	changeHook        func(n int)
	inc, dec          *Clickable
	color, background string
	inactive          string
	scale             float32
}

// IncDec is a simple increment/decrement for a number setting
func (w *Window) IncDec() (out *IncDec) {
	out = &IncDec{
		Window: w,
		// nDigits:    nDigits,
		// min:        min,
		// max:        max,
		// current:    current,
		// changeHook: changeHook,
		inc:        w.Clickable(),
		dec:        w.Clickable(),
		color:      "DocText",
		background: "Transparent",
		inactive:   "Transparent",
		amount:     1,
		scale:      1,
	}
	return
}

func (in *IncDec) Scale(n float32) *IncDec {
	in.scale = n
	return in
}

func (in *IncDec) Amount(n int) *IncDec {
	in.amount = n
	return in
}

func (in *IncDec) ChangeHook(fn func(n int)) *IncDec {
	in.changeHook = fn
	return in
}

func (in *IncDec) SetCurrent(current int) *IncDec {
	in.current = current
	return in
}

func (in *IncDec) GetCurrent() int {
	return in.current
}

func (in *IncDec) Max(max int) *IncDec {
	in.max = max
	return in
}

func (in *IncDec) Min(min int) *IncDec {
	in.min = min
	return in
}

func (in *IncDec) NDigits(nDigits int) *IncDec {
	in.nDigits = nDigits
	return in
}

func (in *IncDec) Color(color string) *IncDec {
	in.color = color
	return in
}

func (in *IncDec) Background(color string) *IncDec {
	in.background = color
	return in
}
func (in *IncDec) Inactive(color string) *IncDec {
	in.inactive = color
	return in
}

func (in *IncDec) Fn(gtx l.Context) l.Dimensions {
	out := in.Theme.Flex().AlignMiddle()
	incColor, decColor := in.color, in.color
	if in.current != in.min {
		out.Rigid(
			in.Inset(
				0.25,
				in.ButtonLayout(
					in.inc.SetClick(
						func() {
							ic := in.current
							ic -= in.amount
							if ic < in.min {
								ic = in.min
							}
							in.current = ic
							in.changeHook(ic)
						},
					),
				).Background(in.background).Embed(
					in.Icon().Color(decColor).Scale(in.scale).Src(&icons.ContentRemove).Fn,
				).Fn,
			).Fn,
		)
	}
	cur := fmt.Sprintf("%"+fmt.Sprint(in.nDigits)+"d", in.current)
	out.Rigid(in.Caption(cur).Color(in.color).TextScale(in.scale).Font("go regular").Fn)
	if in.current != in.max {
		out.Rigid(
			in.Inset(
				0.25,
				in.ButtonLayout(
					in.dec.SetClick(
						func() {
							ic := in.current
							ic += in.amount
							if in.current > in.max {
								in.current = in.max
							} else {
								in.current = ic
								in.changeHook(in.current)
							}
						},
					),
				).Background(in.background).Embed(
					in.Icon().Color(incColor).Scale(in.scale).Src(&icons.ContentAdd).Fn,
				).Fn,
			).Fn,
		)
	}
	return out.Fn(gtx)
}
