package helpers

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

func DuoUIdrawRect(gtx *layout.Context, w, h int, color color.RGBA) {
	square := f32.Rectangle{
		Max: f32.Point{
			X: float32(w),
			Y: float32(h),
		},
	}
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: square}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: image.Point{X: w, Y: h}}
}
