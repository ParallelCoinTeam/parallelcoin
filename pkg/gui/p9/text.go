package p9

import (
	"fmt"
	"image"
	"unicode/utf8"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"golang.org/x/image/math/fixed"
)

// _text is a widget for laying out and drawing text.
type _text struct {
	// alignment specify the text alignment.
	alignment text.Alignment
	// maxLines limits the number of lines. Zero means no limit.
	maxLines int
}

func (th *Theme) Text() *_text {
	return &_text{}
}

type _lineIterator struct {
	lines     []text.Line
	clip      image.Rectangle
	alignment text.Alignment
	width     int
	offset    image.Point

	y, prevDesc fixed.Int26_6
	txtOff      int
}

func (th *Theme) LineIterator() *_lineIterator {
	return &_lineIterator{}
}

func (l *_lineIterator) Lines(lines []text.Line) *_lineIterator {
	l.lines = lines
	return l
}

func (l *_lineIterator) Clip(clip image.Rectangle) *_lineIterator {
	l.clip = clip
	return l
}

func (l *_lineIterator) Alignment(alignment text.Alignment) *_lineIterator {
	l.alignment = alignment
	return l
}

func (l *_lineIterator) Width(width int) *_lineIterator {
	l.width = width
	return l
}

func (l *_lineIterator) Offset(offset image.Point) *_lineIterator {
	l.offset = offset
	return l
}

const inf = 1e6

func (l *_lineIterator) Next() (start, end int, glyph []text.Glyph, offF f32.Point, is bool) {
	for len(l.lines) > 0 {
		line := l.lines[0]
		l.lines = l.lines[1:]
		x := align(l.alignment, line.Width, l.width) + fixed.I(l.offset.X)
		l.y += l.prevDesc + line.Ascent
		l.prevDesc = line.Descent
		// Align baseline and line start to the pixel grid.
		off := fixed.Point26_6{X: fixed.I(x.Floor()), Y: fixed.I(l.y.Ceil())}
		l.y = off.Y
		off.Y += fixed.I(l.offset.Y)
		if (off.Y + line.Bounds.Min.Y).Floor() > l.clip.Max.Y {
			break
		}
		lineLayout := line.Layout
		start = l.txtOff
		l.txtOff += line.Len
		if (off.Y + line.Bounds.Max.Y).Ceil() < l.clip.Min.Y {
			continue
		}
		for len(lineLayout) > 0 {
			g := lineLayout[0]
			adv := g.Advance
			if (off.X + adv + line.Bounds.Max.X - line.Width).Ceil() >= l.clip.Min.X {
				break
			}
			off.X += adv
			lineLayout = lineLayout[1:]
			start += utf8.RuneLen(g.Rune)
		}
		end = start
		endX := off.X
		for i, g := range lineLayout {
			if (endX + line.Bounds.Min.X).Floor() > l.clip.Max.X {
				lineLayout = lineLayout[:i]
				break
			}
			end += utf8.RuneLen(g.Rune)
			endX += g.Advance
		}
		offF = f32.Point{X: float32(off.X) / 64, Y: float32(off.Y) / 64}
		is = true
		return
	}
	return
}

func (l _text) Fn(gtx layout.Context, s text.Shaper, font text.Font, size unit.Value, txt string) layout.Dimensions {
	cs := gtx.Constraints
	textSize := fixed.I(gtx.Px(size))
	lines := s.LayoutString(font, textSize, cs.Max.X, txt)
	if max := l.maxLines; max > 0 && len(lines) > max {
		lines = lines[:max]
	}
	dims := linesDimensions(lines)
	dims.Size = cs.Constrain(dims.Size)
	clip := textPadding(lines)
	clip.Max = clip.Max.Add(dims.Size)
	it := _lineIterator{
		lines:     lines,
		clip:      clip,
		alignment: l.alignment,
		width:     dims.Size.X,
	}
	for {
		start, end, l, off, ok := it.Next()
		if !ok {
			break
		}
		lClip := layout.FRect(clip).Sub(off)
		stack := op.Push(gtx.Ops)
		op.Offset(off).Add(gtx.Ops)
		str := txt[start:end]
		s.ShapeString(font, textSize, str, l).Add(gtx.Ops)
		paint.PaintOp{Rect: lClip}.Add(gtx.Ops)
		stack.Pop()
	}
	return dims
}

func textPadding(lines []text.Line) (padding image.Rectangle) {
	if len(lines) == 0 {
		return
	}
	first := lines[0]
	if d := first.Ascent + first.Bounds.Min.Y; d < 0 {
		padding.Min.Y = d.Ceil()
	}
	last := lines[len(lines)-1]
	if d := last.Bounds.Max.Y - last.Descent; d > 0 {
		padding.Max.Y = d.Ceil()
	}
	if d := first.Bounds.Min.X; d < 0 {
		padding.Min.X = d.Ceil()
	}
	if d := first.Bounds.Max.X - first.Width; d > 0 {
		padding.Max.X = d.Ceil()
	}
	return
}

func linesDimensions(lines []text.Line) layout.Dimensions {
	var width fixed.Int26_6
	var h int
	var baseline int
	if len(lines) > 0 {
		baseline = lines[0].Ascent.Ceil()
		var prevDesc fixed.Int26_6
		for _, l := range lines {
			h += (prevDesc + l.Ascent).Ceil()
			prevDesc = l.Descent
			if l.Width > width {
				width = l.Width
			}
		}
		h += lines[len(lines)-1].Descent.Ceil()
	}
	w := width.Ceil()
	return layout.Dimensions{
		Size: image.Point{
			X: w,
			Y: h,
		},
		Baseline: h - baseline,
	}
}

func align(align text.Alignment, width fixed.Int26_6, maxWidth int) fixed.Int26_6 {
	mw := fixed.I(maxWidth)
	switch align {
	case text.Middle:
		return fixed.I(((mw - width) / 2).Floor())
	case text.End:
		return fixed.I((mw - width).Floor())
	case text.Start:
		return 0
	default:
		panic(fmt.Errorf("unknown alignment %v", align))
	}
}
