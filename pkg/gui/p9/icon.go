// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/iconvg"
)

type _icon struct {
	color color.RGBA
	src   []byte
	size  unit.Value
	// Cached values.
	op       paint.ImageOp
	imgSize  int
	imgColor color.RGBA
}

// Icon returns a new _icon from iconVG data.
func (th *Theme) Icon(data []byte) (*_icon, error) {
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil, err
	}
	return &_icon{src: data, color: rgb(0x000000)}, nil
}

func (ic *_icon) Fn(gtx *l.Context, sz unit.Value) {
	ico := ic.image(gtx.Px(sz))
	ico.Add(gtx.Ops)
	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: toPointF(ico.Size()),
		},
	}.Add(gtx.Ops)
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func (ic *_icon) image(sz int) paint.ImageOp {
	if sz == ic.imgSize && ic.color == ic.imgColor {
		return ic.op
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz,
		Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = ic.color
	iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	ic.imgColor = ic.color
	return ic.op
}
