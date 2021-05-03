package gel

import (
	"image"

	"github.com/p9c/gio/io/pointer"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/text"
	"github.com/p9c/gio/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Checkable struct {
	*Window
	label              string
	color              string
	font               text.Font
	textSize           unit.Value
	iconColor          string
	size               unit.Value
	checkedStateIcon   *[]byte
	uncheckedStateIcon *[]byte
	shaper             text.Shaper
	checked            bool
}

// Checkable creates a checkbox type widget
func (w *Window) Checkable() *Checkable {
	font := "bariol regular"
	var f text.Font
	if fon, e := w.Theme.collection.Font(font); !E.Chk(e) {
		f = fon
	}
	// for i := range w.collection {
	// 	if w.collection[i].Font.Typeface == text.Typeface(font) {
	// 		f = w.collection[i].Font
	// 		break
	// 	}
	// }
	return &Checkable{
		Window:             w,
		label:              "checkable",
		color:              "Primary",
		font:               f,
		textSize:           w.TextSize.Scale(14.0 / 16.0),
		iconColor:          "Primary",
		size:               w.TextSize.Scale(1.5),
		checkedStateIcon:   &icons.ToggleCheckBox,
		uncheckedStateIcon: &icons.ToggleCheckBoxOutlineBlank,
		shaper:             w.shaper,
	}
}

// Label sets the label on the checkbox
func (c *Checkable) Label(txt string) *Checkable {
	c.label = txt
	return c
}

// Color sets the color of the checkbox label
func (c *Checkable) Color(color string) *Checkable {
	c.color = color
	return c
}

// Font sets the font used on the label
func (c *Checkable) Font(font string) *Checkable {
	if fon, e := c.Theme.collection.Font(font); !E.Chk(e) {
		c.font = fon
	}
	return c
}

// TextScale sets the size of the font relative to the base text size
func (c *Checkable) TextScale(scale float32) *Checkable {
	c.textSize = c.Theme.TextSize.Scale(scale)
	return c
}

// IconColor sets the color of the icon
func (c *Checkable) IconColor(color string) *Checkable {
	c.iconColor = color
	return c
}

// Scale sets the size of the checkbox icon relative to the base font size
func (c *Checkable) Scale(size float32) *Checkable {
	c.size = c.Theme.TextSize.Scale(size)
	return c
}

// CheckedStateIcon loads the icon for the checked state
func (c *Checkable) CheckedStateIcon(ic *[]byte) *Checkable {
	c.checkedStateIcon = ic
	return c
}

// UncheckedStateIcon loads the icon for the unchecked state
func (c *Checkable) UncheckedStateIcon(ic *[]byte) *Checkable {
	c.uncheckedStateIcon = ic
	return c
}

// Fn renders the checkbox widget
func (c *Checkable) Fn(gtx l.Context, checked bool) l.Dimensions {
	var icon *Icon
	if checked {
		icon = c.Icon().
			Color(c.iconColor).
			Src(c.checkedStateIcon)
	} else {
		icon = c.Icon().
			Color(c.iconColor).
			Src(c.uncheckedStateIcon)
	}
	icon.size = c.size
	// D.S(icon)
	dims :=
		c.Theme.Flex(). // AlignBaseline().
			Rigid(
				// c.Theme.ButtonInset(0.25,
				func(gtx l.Context) l.Dimensions {
					size := gtx.Px(c.size)
					// icon.color = c.iconColor
					// TODO: maybe make a special code for raw colors to do this kind of alpha
					//  or add a parameter to apply it
					// if gtx.Queue == nil {
					// 	icon.color = f32color.MulAlpha(c.Theme.Colors.Get(icon.color), 150)
					// }
					icon.Fn(gtx)
					return l.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				},
				// ).Fn,
			).
			Rigid(
				// c.Theme.ButtonInset(0.25,
				func(gtx l.Context) l.Dimensions {
					paint.ColorOp{Color: c.Theme.Colors.GetNRGBAFromName(c.color)}.Add(gtx.Ops)
					return c.Caption(c.label).Color(c.color).Fn(gtx)
					// return widget.Label{}.Layout(gtx, c.shaper, c.font, c.textSize, c.label)
				},
				// ).Fn,
			).
			Fn(gtx)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
