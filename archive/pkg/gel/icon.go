package gel

import (
	"image"
	"image/color"
	"image/draw"
	
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/unit"
	"golang.org/x/exp/shiny/iconvg"
)

type Icon struct {
	*Window
	color string
	src   *[]byte
	size  unit.Value
	// Cached values.
	sz       int
	op       paint.ImageOp
	imgSize  int
	imgColor string
}

type IconByColor map[color.NRGBA]paint.ImageOp
type IconBySize map[float32]IconByColor
type IconCache map[*[]byte]IconBySize

// Icon returns a new Icon from iconVG data.
func (w *Window) Icon() *Icon {
	return &Icon{Window: w, size: w.TextSize, color: "DocText"}
}

// Color sets the color of the icon image. It must be called before creating the image
func (i *Icon) Color(color string) *Icon {
	i.color = color
	return i
}

// Src sets the icon source to draw from
func (i *Icon) Src(data *[]byte) *Icon {
	_, e := iconvg.DecodeMetadata(*data)
	if E.Chk(e) {
		D.Ln("no image data, crashing")
		panic(e)
		// return nil
	}
	i.src = data
	return i
}

// Scale changes the size relative to the base font size
func (i *Icon) Scale(scale float32) *Icon {
	i.size = i.Theme.TextSize.Scale(scale)
	return i
}

func (i *Icon) Size(size unit.Value) *Icon {
	i.size = size
	return i
}

// Fn renders the icon
func (i *Icon) Fn(gtx l.Context) l.Dimensions {
	ico := i.image(gtx.Px(i.size))
	if i.src == nil {
		panic("icon is nil")
	}
	ico.Add(gtx.Ops)
	paint.PaintOp{
		// Rect: f32.Rectangle{
		// 	Max: Fpt(ico.Size()),
		// },
	}.Add(gtx.Ops)
	return l.Dimensions{Size: ico.Size()}
}

func (i *Icon) image(sz int) paint.ImageOp {
	// if sz == i.imgSize && i.color == i.imgColor {
	// 	// D.Ln("reusing old icon")
	// 	return i.op
	// }
	if ico, ok := i.Theme.iconCache[i.src]; ok {
		if isz, ok := ico[i.size.V]; ok {
			if icl, ok := isz[i.Theme.Colors.GetNRGBAFromName(i.color)]; ok {
				return icl
			}
		}
	}
	m, _ := iconvg.DecodeMetadata(*i.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz,
		Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = color.RGBA(i.Theme.Colors.GetNRGBAFromName(i.color))
	if e := iconvg.Decode(&ico, *i.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	}); E.Chk(e) {
	}
	operation := paint.NewImageOp(img)
	// create the maps if they don't exist
	if _, ok := i.Theme.iconCache[i.src]; !ok {
		i.Theme.iconCache[i.src] = make(IconBySize)
	}
	if _, ok := i.Theme.iconCache[i.src][i.size.V]; !ok {
		i.Theme.iconCache[i.src][i.size.V] = make(IconByColor)
	}
	i.Theme.iconCache[i.src][i.size.V][i.Theme.Colors.GetNRGBAFromName(i.color)] = operation
	return operation
}
