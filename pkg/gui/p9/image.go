// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Image is a widget that displays an image.
type Image struct {
	// src is the image to display.
	src paint.ImageOp
	// scale is the ratio of image pixels to dps. If scale is zero Image falls back to a scale that match a standard 72
	// DPI.
	scale float32
}

func (th *Theme) Image() *Image {
	return &Image{}
}

func (i *Image) Src(img paint.ImageOp) *Image {
	i.src = img
	return i
}

func (i *Image) Scale(scale float32) *Image {
	i.scale = scale
	return i
}

func (im *Image) Fn(gtx layout.Context) layout.Dimensions {
	scale := im.scale
	if scale == 0 {
		scale = 160.0 / 72.0
	}
	// size := im.src.Rect.Size()
	size := im.src.Size()
	wf, hf := float32(size.X), float32(size.Y)
	w, h := gtx.Px(unit.Dp(wf*scale)), gtx.Px(unit.Dp(hf*scale))
	cs := gtx.Constraints
	d := cs.Constrain(image.Pt(w, h))
	stack := op.Push(gtx.Ops)
	clip.Rect(image.Rectangle{Max: d}).Add(gtx.Ops)
	im.src.Add(gtx.Ops)
// Rect: f32.Rectangle{Max: f32.Point{X: float32(w), Y: float32(h)}}
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Pop()
	return layout.Dimensions{Size: d}
}
