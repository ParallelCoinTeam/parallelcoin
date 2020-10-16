// SPDX-License-Identifier: Unlicense OR MIT

package lyt

import (
	"fmt"
	"image"
	"strconv"

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
func Format(gtx layout.Context, format string, widgets ...layout.Widget) layout.Dimensions {
	if format == "" {
		return layout.Dimensions{}
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
	return formatExpr(gtx, &f, widgets)
}

func formatExpr(gtx layout.Context, f *formatter, widgets []layout.Widget) layout.Dimensions {
	switch peek(f) {
	case '_':
		return formatWidget(gtx, f, widgets)
	default:
		return formatLayout(gtx, f, widgets)
	}
}

func formatLayout(gtx layout.Context, f *formatter, widgets []layout.Widget) layout.Dimensions {
	name := parseName(f)
	if name == "" {
		errorf("missing layout name")
	}
	expect(f, "(")
	fexpr := func(gtx layout.Context) layout.Dimensions {
		return formatExpr(gtx, f, widgets)
	}
	align, ok := dirFor(name)
	var dims layout.Dimensions
	if ok {
		dims = align.Layout(gtx, fexpr)
		expect(f, ")")
		return dims
	}
	switch name {
	case "inset":
		in := parseInset(gtx, f, widgets)
		dims = in.Layout(gtx, fexpr)
	case "hflex":
		dims = formatFlex(gtx, layout.Horizontal, f, widgets)
	case "vflex":
		dims = formatFlex(gtx, layout.Vertical, f, widgets)
	case "stack":
		dims = formatStack(gtx, f, widgets)
	case "hmax":
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		dims = formatExpr(gtx, f, widgets)
	case "vmax":
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		dims = formatExpr(gtx, f, widgets)
	case "max":
		gtx.Constraints.Min = gtx.Constraints.Max
		dims = formatExpr(gtx, f, widgets)
	case "hmin":
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
		dims = formatExpr(gtx, f, widgets)
	case "vmin":
		gtx.Constraints.Max.Y = gtx.Constraints.Min.Y
		dims = formatExpr(gtx, f, widgets)
	case "min":
		gtx.Constraints.Max = gtx.Constraints.Min
		dims = formatExpr(gtx, f, widgets)
	case "hcap":
		w := gtx.Px(parseValue(f))
		expect(f, ",")
		cs := gtx.Constraints
		cs.Max = cs.Constrain(image.Pt(w, cs.Max.X))
		dims = formatExpr(gtx, f, widgets)
	case "vcap":
		h := gtx.Px(parseValue(f))
		expect(f, ",")
		cs := gtx.Constraints
		cs.Max = cs.Constrain(image.Pt(cs.Max.X, h))
		dims = formatExpr(gtx, f, widgets)
	default:
		errorf("invalid layout %q", name)
	}
	expect(f, ")")
	return dims
}

func formatWidget(gtx layout.Context, f *formatter, widgets []layout.Widget) layout.Dimensions {
	expect(f, "_")
	if i, max := f.current, len(widgets)-1; i > max {
		errorf("widget index %d out of bounds [0;%d]", i, max)
	}
	if f.skip == 0 {
		return widgets[f.current](gtx)
	}
	f.current++
	return layout.Dimensions{}
}

func formatStack(gtx layout.Context, f *formatter, widgets []layout.Widget) layout.Dimensions {
	st := layout.Stack{}
	backup := *f
	// Parse alignment, if present.
	name := parseName(f)
	align, ok := dirFor(name)
	if ok {
		st.Alignment = align
		expect(f, ",")
	} else {
		*f = backup
	}
	var children []layout.StackChild
loop:
	for {
		switch peek(f) {
		case ')':
			break loop
		case 'r':
			w := func(gtx layout.Context) layout.Dimensions {
				expect(f, "r(")
				dims := formatExpr(gtx, f, widgets)
				expect(f, ")")
				if peek(f) == ',' {
					expect(f, ",")
				}
				return dims
			}
			backup := *f
			f.skip++
			w(gtx)
			f.skip--
			children = append(children, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				*f = backup
				return w(gtx)
			}))
		case 'e':
			w := func(gtx layout.Context) layout.Dimensions {
				expect(f, "e(")
				dims := formatExpr(gtx, f, widgets)
				expect(f, ")")
				if peek(f) == ',' {
					expect(f, ",")
				}
				return dims
			}
			backup := *f
			f.skip++
			w(gtx)
			f.skip--
			children = append(children, layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				*f = backup
				return w(gtx)
			}))
		default:
			errorf("invalid stack child")
		}
	}
	if f.skip == 0 {
		backup := *f
		dims := st.Layout(gtx, children...)
		*f = backup
		return dims
	}
	return layout.Dimensions{}
}

func formatFlex(gtx layout.Context, axis layout.Axis, f *formatter, widgets []layout.Widget) layout.Dimensions {
	fl := layout.Flex{Axis: axis}
	backup := *f
	// Parse alignment, if present.
	name := parseName(f)
	al, ok := alignmentFor(name)
	if ok {
		fl.Alignment = al
		expect(f, ",")
	} else {
		*f = backup
	}
	var children []layout.FlexChild
loop:
	for {
		switch peek(f) {
		case ')':
			break loop
		case 'r':
			w := func(gtx layout.Context) layout.Dimensions {
				expect(f, "r(")
				dims := formatExpr(gtx, f, widgets)
				expect(f, ")")
				if peek(f) == ',' {
					expect(f, ",")
				}
				return dims
			}
			backup := *f
			f.skip++
			w(gtx)
			f.skip--
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				*f = backup
				return w(gtx)
			}))
		case 'f':
			var weight float32
			w := func(gtx layout.Context) layout.Dimensions {
				expect(f, "f(")
				weight = parseFloat(f)
				expect(f, ",")
				dims := formatExpr(gtx, f, widgets)
				expect(f, ")")
				if peek(f) == ',' {
					expect(f, ",")
				}
				return dims
			}
			backup := *f
			f.skip++
			w(gtx)
			f.skip--
			children = append(children, layout.Flexed(weight, func(gtx layout.Context) layout.Dimensions {
				*f = backup
				return w(gtx)
			}))
		default:
			errorf("invalid flex child")
		}
	}
	if f.skip == 0 {
		backup := *f
		dims := fl.Layout(gtx, children...)
		*f = backup
		return dims
	}
	return layout.Dimensions{}
}

func parseInset(gtx layout.Context, f *formatter, widgets []layout.Widget) layout.Inset {
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
