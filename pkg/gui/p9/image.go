// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// _image is a widget that displays an image.
type _image struct {
	// src is the image to display.
	src paint.ImageOp
	// scale is the ratio of image pixels to dps. If scale is zero _image falls back to a scale that match a standard 72
	// DPI.
	scale float32
}

func (th *Theme) Image() *_image {
	return &_image{}
}

func (i *_image) Src(img paint.ImageOp) *_image {
	i.src = img
	return i
}

func (i *_image) Scale(scale float32) *_image {
	i.scale = scale
	return i
}

func (im *_image) Fn(gtx layout.Context) layout.Dimensions {
	scale := im.scale
	if scale == 0 {
		scale = 160.0 / 72.0
	}
	size := im.src.Rect.Size()
	wf, hf := float32(size.X), float32(size.Y)
	w, h := gtx.Px(unit.Dp(wf*scale)), gtx.Px(unit.Dp(hf*scale))
	cs := gtx.Constraints
	d := cs.Constrain(image.Pt(w, h))
	stack := op.Push(gtx.Ops)
	clip.Rect(image.Rectangle{Max: d}).Add(gtx.Ops)
	im.src.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(w), Y: float32(h)}}}.Add(gtx.Ops)
	stack.Pop()
	return layout.Dimensions{Size: d}
}
