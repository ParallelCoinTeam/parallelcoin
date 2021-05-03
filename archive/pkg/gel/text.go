package gel

import (
	"image"
	"unicode/utf8"

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/text"
	"github.com/p9c/gio/unit"

	"golang.org/x/image/math/fixed"
)

// Text is a widget for laying out and drawing text.
type Text struct {
	*Window
	// alignment specify the text alignment.
	alignment text.Alignment
	// maxLines limits the number of lines. Zero means no limit.
	maxLines int
}

func (w *Window) Text() *Text {
	return &Text{Window: w}
}

// Alignment sets the alignment for the text
func (t *Text) Alignment(alignment text.Alignment) *Text {
	t.alignment = alignment
	return t
}

// MaxLines sets the alignment for the text
func (t *Text) MaxLines(maxLines int) *Text {
	t.maxLines = maxLines
	return t
}

type lineIterator struct {
	Lines     []text.Line
	Clip      image.Rectangle
	Alignment text.Alignment
	Width     int
	Offset    image.Point
	
	y, prevDesc fixed.Int26_6
	txtOff      int
}

func (l *lineIterator) Next() (text.Layout, image.Point, bool) {
	for len(l.Lines) > 0 {
		line := l.Lines[0]
		l.Lines = l.Lines[1:]
		x := align(l.Alignment, line.Width, l.Width) + fixed.I(l.Offset.X)
		l.y += l.prevDesc + line.Ascent
		l.prevDesc = line.Descent
		// Align baseline and line start to the pixel grid.
		off := fixed.Point26_6{X: fixed.I(x.Floor()), Y: fixed.I(l.y.Ceil())}
		l.y = off.Y
		off.Y += fixed.I(l.Offset.Y)
		if (off.Y + line.Bounds.Min.Y).Floor() > l.Clip.Max.Y {
			break
		}
		lo := line.Layout
		start := l.txtOff
		l.txtOff += len(line.Layout.Text)
		if (off.Y + line.Bounds.Max.Y).Ceil() < l.Clip.Min.Y {
			continue
		}
		for len(lo.Advances) > 0 {
			_, n := utf8.DecodeRuneInString(lo.Text)
			adv := lo.Advances[0]
			if (off.X + adv + line.Bounds.Max.X - line.Width).Ceil() >= l.Clip.Min.X {
				break
			}
			off.X += adv
			lo.Text = lo.Text[n:]
			lo.Advances = lo.Advances[1:]
			start += n
		}
		end := start
		endx := off.X
		rn := 0
		for n, r := range lo.Text {
			if (endx + line.Bounds.Min.X).Floor() > l.Clip.Max.X {
				lo.Advances = lo.Advances[:rn]
				lo.Text = lo.Text[:n]
				break
			}
			end += utf8.RuneLen(r)
			endx += lo.Advances[rn]
			rn++
		}
		offf := image.Point{X: off.X.Floor(), Y: off.Y.Floor()}
		return lo, offf, true
	}
	return text.Layout{}, image.Point{}, false
}

// func linesDimens(lines []text.Line) l.Dimensions {
// 	var width fixed.Int26_6
// 	var h int
// 	var baseline int
// 	if len(lines) > 0 {
// 		baseline = lines[0].Ascent.Ceil()
// 		var prevDesc fixed.Int26_6
// 		for _, l := range lines {
// 			h += (prevDesc + l.Ascent).Ceil()
// 			prevDesc = l.Descent
// 			if l.Width > width {
// 				width = l.Width
// 			}
// 		}
// 		h += lines[len(lines)-1].Descent.Ceil()
// 	}
// 	w := width.Ceil()
// 	return l.Dimensions{
// 		Size: image.Point{
// 			X: w,
// 			Y: h,
// 		},
// 		Baseline: h - baseline,
// 	}
// }

func (t *Text) Fn(gtx l.Context, s text.Shaper, font text.Font, size unit.Value, txt string) l.Dimensions {
	cs := gtx.Constraints
	textSize := fixed.I(gtx.Px(size))
	lines := s.LayoutString(font, textSize, cs.Max.X, txt)
	if max := t.maxLines; max > 0 && len(lines) > max {
		lines = lines[:max]
	}
	dims := linesDimens(lines)
	dims.Size = cs.Constrain(dims.Size)
	cl := textPadding(lines)
	cl.Max = cl.Max.Add(dims.Size)
	it := lineIterator{
		Lines:     lines,
		Clip:      cl,
		Alignment: t.alignment,
		Width:     dims.Size.X,
	}
	for {
		ll, off, ok := it.Next()
		if !ok {
			break
		}
		stack := op.Save(gtx.Ops)
		op.Offset(l.FPt(off)).Add(gtx.Ops)
		s.Shape(font, textSize, ll).Add(gtx.Ops)
		clip.Rect(cl.Sub(off)).Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Load()
	}
	return dims
}

// func textPadding(lines []text.Line) (padding image.Rectangle) {
// 	if len(lines) == 0 {
// 		return
// 	}
// 	first := lines[0]
// 	if d := first.Ascent + first.Bounds.Min.Y; d < 0 {
// 		padding.Min.Y = d.Ceil()
// 	}
// 	last := lines[len(lines)-1]
// 	if d := last.Bounds.Max.Y - last.Descent; d > 0 {
// 		padding.Max.Y = d.Ceil()
// 	}
// 	if d := first.Bounds.Min.X; d < 0 {
// 		padding.Min.X = d.Ceil()
// 	}
// 	if d := first.Bounds.Max.X - first.Width; d > 0 {
// 		padding.Max.X = d.Ceil()
// 	}
// 	return
// }

// func align(align text.Alignment, width fixed.Int26_6, maxWidth int) fixed.Int26_6 {
// 	mw := fixed.I(maxWidth)
// 	switch align {
// 	case text.Middle:
// 		return fixed.I(((mw - width) / 2).Floor())
// 	case text.End:
// 		return fixed.I((mw - width).Floor())
// 	case text.Start:
// 		return 0
// 	default:
// 		panic(fmt.Errorf("unknown alignment %v", align))
// 	}
// }
