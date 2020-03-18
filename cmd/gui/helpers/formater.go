// SPDX-License-Identifier: Unlicense OR MIT
package helpers

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
)

type formatter struct {
	current int
	orig    string
	expr    string
	skip    int
}

type formatError string

// Format lays out widgets according to a format string, similar to
// how fmt.Printf interpolates a string.
//
// The format string is an epxression where layouts are similar to
// function calls, and the underscore denotes a widget from the
// arguments. The ith _ invokes the ith widget from the arguments.
//
// If the layout format is invalid, Format panics with an error where
// a cross, ✗, marks the error position.
//
// For example,
//
//   Format(gtx, "inset(8dp, _)", w)
//
// is equivalent to
//
//   layout.UniformInset(unit.Dp(8)).Layout(gtx, w)
//
// Available layouts:
//
// inset(insets, widget) applies inset to widget. Insets are either:
// one value for uniform insets; two values for top/bottom and
// right/left insets; three values for top, bottom and right/left
// insets; or four values for top, right, bottom, left insets.
//
// direction(widget) aligns a widget. Direction is one of north, northeast,
// east, southeast, south, southwest, west, northwest, center.
//
// hmax/vmax/max(widget) forces the horizontal, vertical or both
// constraints to their maximum before laying out widget.
//
// hmin/vmin/min(widget) forces the horizontal, vertical or both
// constraints to their minimum before laying out widget.
//
// hcap/vcap(size, widget) caps the maximum horizontal or vertical
// constraints to size.
//
// hflex/vflex(alignment, children...) lays out children with a
// horizontal or vertical layout.Flex. Each rigid child must be on the form
// r(widget), and each flex child on the form f(<weight>, widget).
// If alignment is specified, it must be one of: start, middle, end,
// baseline. The default alignment is start.
//
// stack(alignment, children) lays out children with a layout.Stack. Each
// Rigid child must be on the form r(widget), and each expand child
// on the form e(widget).
// If alignment is specified it must be one of the directions listed
// above.
func RGB(c uint32) color.RGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func FormatTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

func FormatAmount(a float64) string {
	return fmt.Sprintf("%0.8f", a)
	// return fmt.Sprintf("%d", t.amount)
}
func Alpha(a float64, col color.RGBA) color.RGBA {
	col.A = byte(float64(col.A) * a)
	col.R = byte(float64(col.R) * a)
	col.G = byte(float64(col.G) * a)
	col.B = byte(float64(col.B) * a)
	return col
}

func Format(gtx *layout.Context, format string, widgets ...layout.Widget) {
	if format == "" {
		return
	}
	f := formatter{
		orig: format,
		expr: format,
	}
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(formatError); !ok {
				panic(err)
			}
			pos := len(f.orig) - len(f.expr)
			msg := f.orig[:pos] + "✗" + f.orig[pos:]
			panic(fmt.Errorf("Format: %s:%d: %s", msg, pos, err))
		}
	}()
	formatExpr(gtx, &f, widgets)
}

func formatExpr(gtx *layout.Context, f *formatter, widgets []layout.Widget) {
	switch peek(f) {
	case '_':
		formatWidget(gtx, f, widgets)
	default:
		formatLayout(gtx, f, widgets)
	}
}

func formatLayout(gtx *layout.Context, f *formatter, widgets []layout.Widget) {
	name := parseName(f)
	if name == "" {
		errorf("missing layout name")
	}
	expect(f, "(")
	fexpr := func() {
		formatExpr(gtx, f, widgets)
	}
	// align, ok := dirFor(name)
	// if ok {
	//	layout.Align(align).Layout(gtx, fexpr)
	//	expect(f, ")")
	//	return
	// }
	switch name {
	case "inset":
		in := parseInset(gtx, f, widgets)
		in.Layout(gtx, fexpr)
	case "hflex":
		formatFlex(gtx, layout.Horizontal, f, widgets)
	case "vflex":
		formatFlex(gtx, layout.Vertical, f, widgets)
	case "stack":
		formatStack(gtx, f, widgets)
	case "hmax":
		cs := gtx.Constraints
		cs.Width.Min = cs.Width.Max
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	case "vmax":
		cs := gtx.Constraints
		cs.Height.Min = cs.Height.Max
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	case "max":
		cs := gtx.Constraints
		cs.Width.Min = cs.Width.Max
		cs.Height.Min = cs.Height.Max
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	case "hmin":
		cs := gtx.Constraints
		cs.Width.Max = cs.Width.Min
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	case "vmin":
		cs := gtx.Constraints
		cs.Height.Max = cs.Height.Min
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	case "min":
		cs := gtx.Constraints
		cs.Width.Max = cs.Width.Min
		cs.Height.Max = cs.Height.Min
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	case "hcap":
		w := parseValue(f)
		expect(f, ",")
		cs := gtx.Constraints
		cs.Width.Max = cs.Width.Constrain(gtx.Px(w))
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	case "vcap":
		h := parseValue(f)
		expect(f, ",")
		cs := gtx.Constraints
		cs.Height.Max = cs.Height.Constrain(gtx.Px(h))
		ctxLayout(gtx, cs, func() {
			formatExpr(gtx, f, widgets)
		})
	default:
		errorf("invalid layout %q", name)
	}
	expect(f, ")")
}

func formatWidget(gtx *layout.Context, f *formatter, widgets []layout.Widget) {
	expect(f, "_")
	if i, max := f.current, len(widgets)-1; i > max {
		errorf("widget index %d out of bounds [0;%d]", i, max)
	}
	if f.skip == 0 {
		widgets[f.current]()
	}
	f.current++
}

func formatStack(gtx *layout.Context, f *formatter, widgets []layout.Widget) {
	st := layout.Stack{}
	backup := *f
	// Parse alignment, if present.
	name := parseName(f)
	align, ok := dirFor(name)
	if ok {
		st.Alignment = align
		expect(f, ",")
		backup = *f
	} else {
		*f = backup
	}
	var children []layout.StackChild
	commaOK := false
	// First, lay out rigid children.
loop:
	for {
		switch peek(f) {
		case ')':
			break loop
		case ',':
			if !commaOK {
				errorf("unexpected ,")
			}
			commaOK = false
			expect(f, ",")
		case 'r':
			expect(f, "r(")
			// children = append(children, st.Rigid(gtx, func() {
			//	formatExpr(gtx, f, widgets)
			// }))
			expect(f, ")")
			commaOK = true
		case 'e':
			expect(f, "e(")
			f.skip++
			formatExpr(gtx, f, widgets)
			children = append(children, layout.StackChild{})
			f.skip--
			expect(f, ")")
			commaOK = true
		default:
			errorf("invalid flex child")
		}
	}
	// Then, lay out expanded children.
	*f = backup
	child := 0
	for {
		switch peek(f) {
		case ')':
			if f.skip == 0 {
				st.Layout(gtx, children...)
			}
			return
		case ',':
			expect(f, ",")
		case 'r':
			expect(f, "r(")
			f.skip++
			formatExpr(gtx, f, widgets)
			f.skip--
			expect(f, ")")
			child++
		case 'e':
			expect(f, "e(")
			// children[child] = st.Expand(gtx, func() {
			//	formatExpr(gtx, f, widgets)
			// })
			expect(f, ")")
			child++
		default:
			errorf("invalid flex child")
		}
	}
}

func formatFlex(gtx *layout.Context, axis layout.Axis, f *formatter, widgets []layout.Widget) {
	fl := layout.Flex{Axis: axis}
	backup := *f
	// Parse alignment, if present.
	name := parseName(f)
	al, ok := alignmentFor(name)
	if ok {
		fl.Alignment = al
		expect(f, ",")
		backup = *f
	} else {
		*f = backup
	}
	var children []layout.FlexChild
	commaOK := false
	// First, lay out rigid children.
loop:
	for {
		switch peek(f) {
		case ')':
			break loop
		case ',':
			if !commaOK {
				errorf("unexpected ,")
			}
			expect(f, ",")
		case 'r':
			expect(f, "r(")
			// children = append(children, fl.Rigid(gtx, func() {
			//	formatExpr(gtx, f, widgets)
			// }))
			expect(f, ")")
			commaOK = true
		case 'f':
			expect(f, "f(")
			parseFloat(f)
			expect(f, ",")
			f.skip++
			formatExpr(gtx, f, widgets)
			children = append(children, layout.FlexChild{})
			f.skip--
			expect(f, ")")
			commaOK = true
		default:
			errorf("invalid flex child")
		}
	}
	// Then, lay out flexible children.
	*f = backup
	child := 0
	for {
		switch peek(f) {
		case ')':
			if f.skip == 0 {
				fl.Layout(gtx, children...)
			}
			return
		case ',':
			expect(f, ",")
		case 'r':
			expect(f, "r(")
			f.skip++
			formatExpr(gtx, f, widgets)
			f.skip--
			expect(f, ")")
			child++
		case 'f':
			expect(f, "f(")
			// weight := parseFloat(f)
			expect(f, ",")
			// children[child] = fl.Flex(gtx, weight, func() {
			//	formatExpr(gtx, f, widgets)
			// })
			expect(f, ")")
			child++
		default:
			errorf("invalid flex child")
		}
	}
}

func parseInset(gtx *layout.Context, f *formatter, widgets []layout.Widget) layout.Inset {
	v1 := parseValue(f)
	if peek(f) == ',' {
		expect(f, ",")
		return layout.UniformInset(v1)
	}
	v2 := parseValue(f)
	if peek(f) == ',' {
		expect(f, ",")
		return layout.Inset{
			Top:    v1,
			Right:  v2,
			Bottom: v1,
			Left:   v2,
		}
	}
	v3 := parseValue(f)
	if peek(f) == ',' {
		expect(f, ",")
		return layout.Inset{
			Top:    v1,
			Right:  v2,
			Bottom: v3,
			Left:   v2,
		}
	}
	v4 := parseValue(f)
	expect(f, ",")
	return layout.Inset{
		Top:    v1,
		Right:  v2,
		Bottom: v3,
		Left:   v4,
	}
}

func parseValue(f *formatter) unit.Value {
	i := parseFloat(f)
	if len(f.expr) < 2 {
		errorf("missing unit")
	}
	u := f.expr[:2]
	var v unit.Value
	switch u {
	case "dp":
		v = unit.Dp(i)
	case "sp":
		v = unit.Sp(i)
	case "px":
		v = unit.Px(i)
	default:
		errorf("unknown unit")
	}
	f.expr = f.expr[len(u):]
	return v
}

func parseName(f *formatter) string {
	skipWhitespace(f)
	i := 0
loop:
	for ; i < len(f.expr); i++ {
		c := f.expr[i]
		switch {
		case c == '(' || c == ',' || c == ')':
			break loop
		case c < 'a' || 'z' < c:
			errorf("invalid character '%c' in layout name", c)
		}
	}
	fname := f.expr[:i]
	f.expr = f.expr[i:]
	return fname
}

func parseFloat(f *formatter) float32 {
	skipWhitespace(f)
	i := 0
	for ; i < len(f.expr); i++ {
		c := f.expr[i]
		if (c < '0' || c > '9') && c != '.' {
			break
		}
	}
	expr := f.expr[:i]
	v, err := strconv.ParseFloat(expr, 32)
	if err != nil {
		errorf("invalid number %q", expr)
	}
	f.expr = f.expr[i:]
	return float32(v)
}

func parseInt(f *formatter) int {
	skipWhitespace(f)
	i := 0
	for ; i < len(f.expr); i++ {
		c := f.expr[i]
		if c < '0' || c > '9' {
			break
		}
	}
	expr := f.expr[:i]
	v, err := strconv.Atoi(expr)
	if err != nil {
		errorf("invalid number %q", expr)
	}
	f.expr = f.expr[i:]
	return v
}

func peek(f *formatter) rune {
	skipWhitespace(f)
	if len(f.expr) == 0 {
		errorf("unexpected end")
	}
	return rune(f.expr[0])
}

func expect(f *formatter, str string) {
	skipWhitespace(f)
	n := len(str)
	if len(f.expr) < n || f.expr[:n] != str {
		errorf("expected %q", str)
	}
	f.expr = f.expr[n:]
}

func skipWhitespace(f *formatter) {
	for len(f.expr) > 0 {
		switch f.expr[0] {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			f.expr = f.expr[1:]
		default:
			return
		}
	}
}

func alignmentFor(name string) (layout.Alignment, bool) {
	var a layout.Alignment
	switch name {
	case "start":
		a = layout.Start
	case "middle":
		a = layout.Middle
	case "end":
		a = layout.End
	case "baseline":
		a = layout.Baseline
	default:
		return 0, false
	}
	return a, true
}

func dirFor(name string) (layout.Direction, bool) {
	var d layout.Direction
	switch name {
	case "center":
		d = layout.Center
	case "northwest":
		d = layout.NW
	case "north":
		d = layout.N
	case "northeast":
		d = layout.NE
	case "east":
		d = layout.E
	case "southeast":
		d = layout.SE
	case "south":
		d = layout.S
	case "southwest":
		d = layout.SW
	case "west":
		d = layout.W
	default:
		return 0, false
	}
	return d, true
}

func errorf(f string, args ...interface{}) {
	panic(formatError(fmt.Sprintf(f, args...)))
}

func (e formatError) Error() string {
	return string(e)
}

// layout a widget with a set of constraints and return its
// dimensions. The widget dimensions are constrained abd the previous
// constraints are restored after layout.
func ctxLayout(gtx *layout.Context, cs layout.Constraints, w layout.Widget) layout.Dimensions {
	saved := gtx.Constraints
	gtx.Constraints = cs
	gtx.Dimensions = layout.Dimensions{}
	w()
	gtx.Dimensions.Size = cs.Constrain(gtx.Dimensions.Size)
	gtx.Constraints = saved
	return gtx.Dimensions
}
