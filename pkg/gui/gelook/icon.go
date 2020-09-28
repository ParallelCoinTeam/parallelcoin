// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"github.com/stalker-loki/app/slog"
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/iconvg"
)

type DuoUIIcon struct {
	Color color.RGBA
	src   []byte
	size  unit.Value
	// Cached values.
	op       paint.ImageOp
	imgSize  int
	imgColor color.RGBA
}

// NewDuoUIIcon returns a new DuoUIIcon from DuoUIIconVG data.
func NewDuoUIIcon(data []byte) (ic *DuoUIIcon, err error) {
	if _, err = iconvg.DecodeMetadata(data); slog.Check(err) {
		return
	}
	return &DuoUIIcon{src: data, Color: rgb(0x000000)}, nil
}

func (ic *DuoUIIcon) Layout(gtx *layout.Context, sz unit.Value) {
	ico := ic.image(gtx.Px(sz))
	ico.Add(gtx.Ops)
	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: toPointF(ico.Size()),
		},
	}.Add(gtx.Ops)
}

func (ic *DuoUIIcon) image(sz int) paint.ImageOp {
	if sz == ic.imgSize && ic.Color == ic.imgColor {
		return ic.op
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz,
		Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = ic.Color
	if err := iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{Palette: &m.Palette}); slog.Check(err) {
	}
	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	ic.imgColor = ic.Color
	return ic.op
}
