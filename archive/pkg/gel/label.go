package gel

import (
	"fmt"
	"image"
	"image/color"
	"unicode/utf8"

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/text"
	"github.com/p9c/gio/unit"
	"golang.org/x/image/math/fixed"
)

// Label is text drawn inside an empty box
type Label struct {
	*Window
	// Face defines the text style.
	font text.Font
	// Color is the text color.
	color color.NRGBA
	// Alignment specify the text alignment.
	alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	maxLines int
	text     string
	textSize unit.Value
	shaper   text.Shaper
}

// screenPos describes a character position (in text line and column numbers,
// not pixels): Y = line number, X = rune column.
type screenPos image.Point

type segmentIterator struct {
	Lines     []text.Line
	Clip      image.Rectangle
	Alignment text.Alignment
	Width     int
	Offset    image.Point
	startSel  screenPos
	endSel    screenPos

	pos    screenPos   // current position
	line   text.Line   // current line
	layout text.Layout // current line's Layout

	// pixel positions
	off         fixed.Point26_6
	y, prevDesc fixed.Int26_6
}

func (l *segmentIterator) Next() (text.Layout, image.Point, bool, int, image.Point, bool) {
	for l.pos.Y < len(l.Lines) {
		if l.pos.X == 0 {
			l.line = l.Lines[l.pos.Y]

			// Calculate X & Y pixel coordinates of left edge of line. We need y
			// for the next line, so it's in l, but we only need x here, so it's
			// not.
			x := align(l.Alignment, l.line.Width, l.Width) + fixed.I(l.Offset.X)
			l.y += l.prevDesc + l.line.Ascent
			l.prevDesc = l.line.Descent
			// Align baseline and line start to the pixel grid.
			l.off = fixed.Point26_6{X: fixed.I(x.Floor()), Y: fixed.I(l.y.Ceil())}
			l.y = l.off.Y
			l.off.Y += fixed.I(l.Offset.Y)
			if (l.off.Y + l.line.Bounds.Min.Y).Floor() > l.Clip.Max.Y {
				break
			}

			if (l.off.Y + l.line.Bounds.Max.Y).Ceil() < l.Clip.Min.Y {
				// This line is outside/before the clip area; go on to the next line.
				l.pos.Y++
				continue
			}

			// Copy the line's Layout, since we slice it up later.
			l.layout = l.line.Layout

			// Find the left edge of the text visible in the l.Clip clipping
			// area.
			for len(l.layout.Advances) > 0 {
				_, n := utf8.DecodeRuneInString(l.layout.Text)
				adv := l.layout.Advances[0]
				if (l.off.X + adv + l.line.Bounds.Max.X - l.line.Width).Ceil() >= l.Clip.Min.X {
					break
				}
				l.off.X += adv
				l.layout.Text = l.layout.Text[n:]
				l.layout.Advances = l.layout.Advances[1:]
				l.pos.X++
			}
		}

		selected := l.inSelection()
		endx := l.off.X
		rune := 0
		nextLine := true
		retLayout := l.layout
		for n := range l.layout.Text {
			selChanged := selected != l.inSelection()
			beyondClipEdge := (endx + l.line.Bounds.Min.X).Floor() > l.Clip.Max.X
			if selChanged || beyondClipEdge {
				retLayout.Advances = l.layout.Advances[:rune]
				retLayout.Text = l.layout.Text[:n]
				if selChanged {
					// Save the rest of the line
					l.layout.Advances = l.layout.Advances[rune:]
					l.layout.Text = l.layout.Text[n:]
					nextLine = false
				}
				break
			}
			endx += l.layout.Advances[rune]
			rune++
			l.pos.X++
		}
		offFloor := image.Point{X: l.off.X.Floor(), Y: l.off.Y.Floor()}

		// Calculate the width & height if the returned text.
		//
		// If there's a better way to do this, I'm all ears.
		var d fixed.Int26_6
		for _, adv := range retLayout.Advances {
			d += adv
		}
		size := image.Point{
			X: d.Ceil(),
			Y: (l.line.Ascent + l.line.Descent).Ceil(),
		}

		if nextLine {
			l.pos.Y++
			l.pos.X = 0
		} else {
			l.off.X = endx
		}

		return retLayout, offFloor, selected, l.prevDesc.Ceil() - size.Y, size, true
	}
	return text.Layout{}, image.Point{}, false, 0, image.Point{}, false
}

func (l *segmentIterator) inSelection() bool {
	return l.startSel.LessOrEqual(l.pos) &&
		l.pos.Less(l.endSel)
}

func (p1 screenPos) LessOrEqual(p2 screenPos) bool {
	return p1.Y < p2.Y || (p1.Y == p2.Y && p1.X <= p2.X)
}

func (p1 screenPos) Less(p2 screenPos) bool {
	return p1.Y < p2.Y || (p1.Y == p2.Y && p1.X < p2.X)
}

// Fn renders the label as specified
func (l *Label) Fn(gtx l.Context) l.Dimensions {
	cs := gtx.Constraints
	textSize := fixed.I(gtx.Px(l.textSize))
	lines := l.shaper.LayoutString(l.font, textSize, cs.Max.X, l.text)
	if max := l.maxLines; max > 0 && len(lines) > max {
		lines = lines[:max]
	}
	dims := linesDimens(lines)
	dims.Size = cs.Constrain(dims.Size)
	cl := textPadding(lines)
	cl.Max = cl.Max.Add(dims.Size)
	it := segmentIterator{
		Lines:     lines,
		Clip:      cl,
		Alignment: l.alignment,
		Width:     dims.Size.X,
	}
	for {
		lb, off, _, _, _, ok := it.Next()
		if !ok {
			break
		}
		stack := op.Save(gtx.Ops)
		op.Offset(Fpt(off)).Add(gtx.Ops)
		l.shaper.Shape(l.font, textSize, lb).Add(gtx.Ops)
		clip.Rect(cl.Sub(off)).Add(gtx.Ops)
		paint.ColorOp{Color: l.color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Load()
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

func linesDimens(lines []text.Line) l.Dimensions {
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
	return l.Dimensions{
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

// ScaleType is a map of the set of label txsizes
type ScaleType map[string]float32

// Scales is the ratios against
//
// TODO: shouldn't that 16.0 be the text size in the theme?
var Scales = ScaleType{
	"H1":      96.0 / 16.0,
	"H2":      60.0 / 16.0,
	"H3":      48.0 / 16.0,
	"H4":      34.0 / 16.0,
	"H5":      24.0 / 16.0,
	"H6":      20.0 / 16.0,
	"Body1":   1,
	"Body2":   14.0 / 16.0,
	"Caption": 12.0 / 16.0,
}

// Label creates a label that prints a block of text
func (w *Window) Label() (l *Label) {
	var f text.Font
	var e error
	var fon text.Font
	if fon, e = w.Theme.collection.Font("plan9"); !E.Chk(e) {
		f = fon
	}
	return &Label{
		Window:   w,
		text:     "",
		font:     f,
		color:    w.Colors.GetNRGBAFromName("DocText"),
		textSize: unit.Sp(1),
		shaper:   w.shaper,
	}
}

// Text sets the text to render in the label
func (l *Label) Text(text string) *Label {
	l.text = text
	return l
}

// TextScale sets the size of the text relative to the base font size
func (l *Label) TextScale(scale float32) *Label {
	l.textSize = l.Theme.TextSize.Scale(scale)
	return l
}

// MaxLines sets the maximum number of lines to render
func (l *Label) MaxLines(maxLines int) *Label {
	l.maxLines = maxLines
	return l
}

// Alignment sets the text alignment, left, right or centered
func (l *Label) Alignment(alignment text.Alignment) *Label {
	l.alignment = alignment
	return l
}

// Color sets the color of the label font
func (l *Label) Color(color string) *Label {
	l.color = l.Theme.Colors.GetNRGBAFromName(color)
	return l
}

// Font sets the font out of the available font collection
func (l *Label) Font(font string) *Label {
	var e error
	var fon text.Font
	if fon, e = l.Theme.collection.Font(font); !E.Chk(e) {
		l.font = fon
	}
	return l
}

// H1 header 1
func (w *Window) H1(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H1"]).Font("plan9").Text(txt)
	return
}

// H2 header 2
func (w *Window) H2(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H2"]).Font("plan9").Text(txt)
	return
}

// H3 header 3
func (w *Window) H3(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H3"]).Font("plan9").Text(txt)
	return
}

// H4 header 4
func (w *Window) H4(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H4"]).Font("plan9").Text(txt)
	return
}

// H5 header 5
func (w *Window) H5(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H5"]).Font("plan9").Text(txt)
	return
}

// H6 header 6
func (w *Window) H6(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H6"]).Font("plan9").Text(txt)
	return
}

// Body1 normal body text 1
func (w *Window) Body1(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["Body1"]).Font("bariol regular").Text(txt)
	return
}

// Body2 normal body text 2
func (w *Window) Body2(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["Body2"]).Font("bariol regular").Text(txt)
	return
}

// Caption caption text
func (w *Window) Caption(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["Caption"]).Font("bariol regular").Text(txt)
	return
}
