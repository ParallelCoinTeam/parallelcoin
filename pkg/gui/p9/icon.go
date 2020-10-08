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

type Ico struct {
	th    *Theme
	color color.RGBA
	src   []byte
	size  unit.Value
	// Cached values.
	sz       int
	op       paint.ImageOp
	imgSize  int
	imgColor color.RGBA
}

// Icon returns a new _icon from iconVG data.
func (th *Theme) Icon() *Ico {
	return &Ico{th: th, size: th.textSize, color: rgb(0xff000000)}
}

// Color sets the color of the icon image. It must be called before creating the image
func (i *Ico) Color(color string) *Ico {
	i.color = i.th.Colors.Get(color)
	return i
}

// Src sets the icon source to draw from
func (i *Ico) Src(data []byte) *Ico {
	_, err := iconvg.DecodeMetadata(data)
	if Check(err) {
		return nil
	}
	i.src = data
	return i
}

// Scale changes the size relative to the base font size
func (i *Ico) Scale(scale float32) *Ico {
	i.size = i.th.textSize.Scale(scale)
	return i
}

// Fn renders the icon
func (i *Ico) Fn(gtx l.Context) l.Dimensions {
	ico := i.image(gtx.Px(i.size))
	ico.Add(gtx.Ops)
	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: toPointF(ico.Size()),
		},
	}.Add(gtx.Ops)
	return l.Dimensions{Size: ico.Size()}
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func (i *Ico) image(sz int) paint.ImageOp {
	if sz == i.imgSize && i.color == i.imgColor {
		return i.op
	}
	m, _ := iconvg.DecodeMetadata(i.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz,
		Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = i.color
	iconvg.Decode(&ico, i.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	i.op = paint.NewImageOp(img)
	i.imgSize = sz
	i.imgColor = i.color
	return i.op
}
