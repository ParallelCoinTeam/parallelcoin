// SPDX-License-Identifier: Unlicense OR MIT

package gui

import (
	"bufio"
	"bytes"
	"image"
	"io"
	"math"
	"runtime"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	
	"gioui.org/io/clipboard"
	"gioui.org/op/clip"
	
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	
	"golang.org/x/image/math/fixed"
)

func (w *Window) Editor() *Editor {
	e := &Editor{
		submitHook: func(string) {},
		changeHook: func(string) {},
		focusHook:  func(bool) {},
	}
	return e
}

// Editor implements an editable and scrollable text area.
type Editor struct {
	alignment text.Alignment
	// singleLine force the text to stay on a single line. singleLine also sets the scrolling direction to horizontal.
	singleLine bool
	// submit enabled translation of carriage return keys to SubmitEvents. If not enabled, carriage returns are inserted
	// as newlines in the text.
	submit bool
	// mask replaces the visual display of each rune in the contents with the given rune. Newline characters are not
	// masked. When non-zero, the unmasked contents are accessed by Len, Text, and SetText.
	mask rune
	
	eventKey     int
	font         text.Font
	shaper       text.Shaper
	textSize     fixed.Int26_6
	blinkStart   time.Time
	focused      bool
	rr           editBuffer
	maskReader   maskReader
	lastMask     rune
	maxWidth     int
	viewSize     image.Point
	valid        bool
	lines        []text.Line
	shapes       []line
	dims         layout.Dimensions
	requestFocus bool
	Caret        struct {
		on     bool
		scroll bool
		
		// xoff is the offset to the current caret position when moving between lines.
		xoff fixed.Int26_6
		
		// Line is the caret line position as an index into lines.
		Line int
		// Col is the caret column measured in runes.
		Col int
		// (x, y) are the caret coordinates.
		x fixed.Int26_6
		y int
	}
	scroller  gesture.Scroll
	scrollOff image.Point
	
	clicker gesture.Click
	
	// events is the list of events not yet processed.
	events []EditorEvent
	// prevEvents is the number of events from the previous frame.
	prevEvents int
	submitHook func(string)
	changeHook func(string)
	focusHook  func(bool)
}

type maskReader struct {
	// rr is the underlying reader.
	rr      io.RuneReader
	maskBuf [utf8.UTFMax]byte
	// mask is the utf-8 encoded mask rune.
	mask []byte
	// overflow contains excess mask bytes left over after the last Read call.
	overflow []byte
}

func (m *maskReader) Reset(r io.RuneReader, mr rune) {
	m.rr = r
	n := utf8.EncodeRune(m.maskBuf[:], mr)
	m.mask = m.maskBuf[:n]
}

// Read reads from the underlying reader and replaces every rune with the mask rune.
func (m *maskReader) Read(b []byte) (n int, err error) {
	for len(b) > 0 {
		var replacement []byte
		if len(m.overflow) > 0 {
			replacement = m.overflow
		} else {
			var r rune
			r, _, err = m.rr.ReadRune()
			if err != nil {
				break
			}
			if r == '\n' {
				replacement = []byte{'\n'}
			} else {
				replacement = m.mask
			}
		}
		nn := copy(b, replacement)
		m.overflow = replacement[nn:]
		n += nn
		b = b[nn:]
	}
	return n, err
}

type EditorEvent interface {
	isEditorEvent()
}


// A ChangeEvent is generated for every user change to the text.
type ChangeEvent struct{}

// A SubmitEvent is generated when submit is set and a carriage return key is pressed.
type SubmitEvent struct {
	Text string
}

type line struct {
	offset image.Point
	clip   op.CallOp
}

const (
	blinksPerSecond  = 1
	maxBlinkDuration = 10 * time.Second
)

// Events returns available editor events.
func (e *Editor) Events() []EditorEvent {
	events := e.events
	e.events = nil
	e.prevEvents = 0
	return events
}

func (e *Editor) processEvents(gtx layout.Context) {
	// Flush events from before the previous Open.
	n := copy(e.events, e.events[e.prevEvents:])
	e.events = e.events[:n]
	e.prevEvents = n
	if e.shaper == nil {
		// Can't process events without a shaper.
		return
	}
	e.processPointer(gtx)
	e.processKey(gtx)
}


func (e *Editor) Alignment(alignment text.Alignment) *Editor {
	e.alignment = alignment
	return e
}

func (e *Editor) SingleLine() *Editor {
	e.singleLine = true
	return e
}

func (e *Editor) Submit(submit bool) *Editor {
	e.submit = submit
	return e
}

func (e *Editor) Mask(mask rune) *Editor {
	e.mask = mask
	return e
}

func (e *Editor) SetSubmit(submitFn func(txt string)) *Editor {
	e.submitHook = submitFn
	return e
}

func (e *Editor) SetChange(changeFn func(txt string)) *Editor {
	e.changeHook = changeFn
	return e
}

func (e *Editor) SetFocus(focusFn func(is bool)) *Editor {
	e.focusHook = focusFn
	return e
}
func (e *Editor) makeValid() {
	if e.valid {
		return
	}
	e.lines, e.dims = e.layoutText(e.shaper)
	line, col, x, y := e.layoutCaret()
	e.Caret.Line = line
	e.Caret.Col = col
	e.Caret.x = x
	e.Caret.y = y
	e.valid = true
}

func (e *Editor) processPointer(gtx layout.Context) {
	sbounds := e.scrollBounds()
	var smin, smax int
	var axis gesture.Axis
	if e.singleLine {
		axis = gesture.Horizontal
		smin, smax = sbounds.Min.X, sbounds.Max.X
	} else {
		axis = gesture.Vertical
		smin, smax = sbounds.Min.Y, sbounds.Max.Y
	}
	sdist := e.scroller.Scroll(gtx.Metric, gtx, gtx.Now, axis)
	var soff int
	if e.singleLine {
		e.scrollRel(sdist, 0)
		soff = e.scrollOff.X
	} else {
		e.scrollRel(0, sdist)
		soff = e.scrollOff.Y
	}
	for _, evt := range e.clicker.Events(gtx) {
		switch {
		case evt.Type == gesture.TypePress && evt.Source == pointer.Mouse,
			evt.Type == gesture.TypeClick && evt.Source == pointer.Touch:
			e.blinkStart = gtx.Now
			e.moveCoord(image.Point{
				X: int(math.Round(float64(evt.Position.X))),
				Y: int(math.Round(float64(evt.Position.Y))),
			})
			e.requestFocus = true
			if e.scroller.State() != gesture.StateFlinging {
				e.Caret.scroll = true
			}
		}
	}
	if (sdist > 0 && soff >= smax) || (sdist < 0 && soff <= smin) {
		e.scroller.Stop()
	}
}

func (e *Editor) processKey(gtx layout.Context) {
	if e.rr.Changed() {
		e.events = append(e.events, ChangeEvent{})
	}
	for _, ke := range gtx.Events(&e.eventKey) {
		e.blinkStart = gtx.Now
		switch ke := ke.(type) {
		case key.FocusEvent:
			e.focused = ke.Focus
			e.focusHook(ke.Focus)
		case key.Event:
			if !e.focused {
				break
			}
			if e.submit && (ke.Name == key.NameReturn || ke.Name == key.NameEnter) {
				if ke.State == key.Release {
					break
				}
				if !ke.Modifiers.Contain(key.ModShift) {
					e.events = append(e.events, SubmitEvent{
						Text: e.Text(),
					})
					e.submitHook(e.Text())
					return
				}
			}
			if e.command(ke) {
				e.Caret.scroll = true
				e.scroller.Stop()
			}
		case key.EditEvent:
			e.Caret.scroll = true
			e.scroller.Stop()
			e.append(ke.Text)
		case clipboard.Event:
			e.Caret.scroll = true
			e.scroller.Stop()
			e.append(ke.Text)
		}
		if e.rr.Changed() {
			e.events = append(e.events, ChangeEvent{})
			e.changeHook(e.Text())
		}
	}
}

func (e *Editor) moveLines(distance int) {
	e.moveToLine(e.Caret.x+e.Caret.xoff, e.Caret.Line+distance)
}

func (e *Editor) command(k key.Event) bool {
	modSkip := key.ModCtrl
	if runtime.GOOS == "darwin" {
		modSkip = key.ModAlt
	}
	if k.State == key.Release {
		return false
	}
	switch k.Name {
	case key.NameReturn, key.NameEnter:
		e.append("\n")
	case key.NameDeleteBackward:
		if k.Modifiers == modSkip {
			e.deleteWord(-1)
		} else {
			e.Delete(-1)
		}
	case key.NameDeleteForward:
		if k.Modifiers == modSkip {
			e.deleteWord(1)
		} else {
			e.Delete(1)
		}
	case key.NameUpArrow:
		e.moveLines(-1)
	case key.NameDownArrow:
		e.moveLines(+1)
	case key.NameLeftArrow:
		if k.Modifiers == modSkip {
			e.moveWord(-1)
		} else {
			e.Move(-1)
		}
	case key.NameRightArrow:
		if k.Modifiers == modSkip {
			e.moveWord(1)
		} else {
			e.Move(1)
		}
	case key.NamePageUp:
		e.movePages(-1)
	case key.NamePageDown:
		e.movePages(+1)
	case key.NameHome:
		e.moveStart()
	case key.NameEnd:
		e.moveEnd()
	default:
		return false
	}
	return true
}

// Focus requests the input focus for the _editor.
func (e *Editor) Focus() {
	e.requestFocus = true
}

// Focused returns whether the editor is focused or not.
func (e *Editor) Focused() bool {
	return e.focused
}

// Layout lays out the editor.
func (e *Editor) Layout(gtx layout.Context, sh text.Shaper, font text.Font, size unit.Value) layout.Dimensions {
	textSize := fixed.I(gtx.Px(size))
	if e.font != font || e.textSize != textSize {
		e.invalidate()
		e.font = font
		e.textSize = textSize
	}
	maxWidth := gtx.Constraints.Max.X
	if e.singleLine {
		maxWidth = Inf
	}
	if maxWidth != e.maxWidth {
		e.maxWidth = maxWidth
		e.invalidate()
	}
	if sh != e.shaper {
		e.shaper = sh
		e.invalidate()
	}
	if e.mask != e.lastMask {
		e.lastMask = e.mask
		e.invalidate()
	}
	e.makeValid()
	e.processEvents(gtx)
	e.makeValid()
	if viewSize := gtx.Constraints.Constrain(e.dims.Size); viewSize != e.viewSize {
		e.viewSize = viewSize
		e.invalidate()
	}
	e.makeValid()
	return e.layout(gtx)
}

func (e *Editor) layout(gtx layout.Context) layout.Dimensions {
	// Adjust scrolling for new viewport and layout.
	e.scrollRel(0, 0)
	if e.Caret.scroll {
		e.Caret.scroll = false
		e.scrollToCaret()
	}
	off := image.Point{
		X: -e.scrollOff.X,
		Y: -e.scrollOff.Y,
	}
	cl := textPadding(e.lines)
	cl.Max = cl.Max.Add(e.viewSize)
	it := lineIterator{
		Lines:     e.lines,
		Clip:      cl,
		Alignment: e.alignment,
		Width:     e.viewSize.X,
		Offset:    off,
	}
	e.shapes = e.shapes[:0]
	for {
		lo, off, ok := it.Next()
		if !ok {
			break
		}
		path := e.shaper.Shape(e.font, e.textSize, lo)
		e.shapes = append(e.shapes, line{off, path})
	}
	key.InputOp{Tag: &e.eventKey}.Add(gtx.Ops)
	if e.requestFocus {
		key.FocusOp{Focus: true}.Add(gtx.Ops)
		key.SoftKeyboardOp{Show: true}.Add(gtx.Ops)
	}
	e.requestFocus = false
	pointerPadding := gtx.Px(unit.Dp(4))
	r := image.Rectangle{Max: e.viewSize}
	r.Min.X -= pointerPadding
	r.Min.Y -= pointerPadding
	r.Max.X += pointerPadding
	r.Max.X += pointerPadding
	pointer.Rect(r).Add(gtx.Ops)
	pointer.CursorNameOp{Name: pointer.CursorText}.Add(gtx.Ops)
	if !e.singleLine {
		e.scroller.Add(gtx.Ops)
	}
	e.clicker.Add(gtx.Ops)
	e.Caret.on = false
	if e.focused {
		now := gtx.Now
		dt := now.Sub(e.blinkStart)
		blinking := dt < maxBlinkDuration
		const timePerBlink = time.Second / blinksPerSecond
		nextBlink := now.Add(timePerBlink/2 - dt%(timePerBlink/2))
		if blinking {
			redraw := op.InvalidateOp{At: nextBlink}
			redraw.Add(gtx.Ops)
		}
		e.Caret.on = e.focused && (!blinking || dt%timePerBlink < timePerBlink/2)
	}
	
	return layout.Dimensions{Size: e.viewSize, Baseline: e.dims.Baseline}
}

func (e *Editor) PaintText(gtx layout.Context) {
	cl := textPadding(e.lines)
	cl.Max = cl.Max.Add(e.viewSize)
	for _, shape := range e.shapes {
		stack := op.Push(gtx.Ops)
		op.Offset(layout.FPt(shape.offset)).Add(gtx.Ops)
		shape.clip.Add(gtx.Ops)
		clip.Rect(cl.Sub(shape.offset)).Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Pop()
	}
}

func (e *Editor) PaintCaret(gtx layout.Context) {
	if !e.Caret.on {
		return
	}
	e.makeValid()
	carWidth := fixed.I(gtx.Px(unit.Dp(1)))
	carX := e.Caret.x
	carY := e.Caret.y
	
	defer op.Push(gtx.Ops).Pop()
	carX -= carWidth / 2
	carAsc, carDesc := -e.lines[e.Caret.Line].Bounds.Min.Y, e.lines[e.Caret.Line].Bounds.Max.Y
	carRect := image.Rectangle{
		Min: image.Point{X: carX.Ceil(), Y: carY - carAsc.Ceil()},
		Max: image.Point{X: carX.Ceil() + carWidth.Ceil(), Y: carY + carDesc.Ceil()},
	}
	carRect = carRect.Add(image.Point{
		X: -e.scrollOff.X,
		Y: -e.scrollOff.Y,
	})
	cl := textPadding(e.lines)
	// Account for caret width to each side.
	whalf := (carWidth / 2).Ceil()
	if cl.Max.X < whalf {
		cl.Max.X = whalf
	}
	if cl.Min.X > -whalf {
		cl.Min.X = -whalf
	}
	cl.Max = cl.Max.Add(e.viewSize)
	carRect = cl.Intersect(carRect)
	if !carRect.Empty() {
		st := op.Push(gtx.Ops)
		clip.Rect(carRect).Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		st.Pop()
	}
}

// Len is the length of the editor contents.
func (e *Editor) Len() int {
	return e.rr.len()
}

// Text returns the contents of the editor.
func (e *Editor) Text() string {
	return e.rr.String()
}

// SetText replaces the contents of the editor.
func (e *Editor) SetText(s string) *Editor {
	// this isn't necessary for normal inputs but should be done to password inputs, it isn't expensive anyway
	e.rr.Zero()
	e.rr = editBuffer{}
	e.Caret.xoff = 0
	e.prepend(s)
	return e
}

func (e *Editor) scrollBounds() image.Rectangle {
	var b image.Rectangle
	if e.singleLine {
		if len(e.lines) > 0 {
			b.Min.X = align(e.alignment, e.lines[0].Width, e.viewSize.X).Floor()
			if b.Min.X > 0 {
				b.Min.X = 0
			}
		}
		b.Max.X = e.dims.Size.X + b.Min.X - e.viewSize.X
	} else {
		b.Max.Y = e.dims.Size.Y - e.viewSize.Y
	}
	return b
}

func (e *Editor) scrollRel(dx, dy int) {
	e.scrollAbs(e.scrollOff.X+dx, e.scrollOff.Y+dy)
}

func (e *Editor) scrollAbs(x, y int) {
	e.scrollOff.X = x
	e.scrollOff.Y = y
	b := e.scrollBounds()
	if e.scrollOff.X > b.Max.X {
		e.scrollOff.X = b.Max.X
	}
	if e.scrollOff.X < b.Min.X {
		e.scrollOff.X = b.Min.X
	}
	if e.scrollOff.Y > b.Max.Y {
		e.scrollOff.Y = b.Max.Y
	}
	if e.scrollOff.Y < b.Min.Y {
		e.scrollOff.Y = b.Min.Y
	}
}

func (e *Editor) moveCoord(pos image.Point) {
	var (
		prevDesc fixed.Int26_6
		carLine  int
		y        int
	)
	for _, l := range e.lines {
		y += (prevDesc + l.Ascent).Ceil()
		prevDesc = l.Descent
		if y+prevDesc.Ceil() >= pos.Y+e.scrollOff.Y {
			break
		}
		carLine++
	}
	x := fixed.I(pos.X + e.scrollOff.X)
	e.moveToLine(x, carLine)
	e.Caret.xoff = 0
}

func (e *Editor) layoutText(s text.Shaper) ([]text.Line, layout.Dimensions) {
	e.rr.Reset()
	var r io.Reader = &e.rr
	if e.mask != 0 {
		e.maskReader.Reset(&e.rr, e.mask)
		r = &e.maskReader
	}
	var lines []text.Line
	if s != nil {
		lines, _ = s.Layout(e.font, e.textSize, e.maxWidth, r)
	} else {
		lines, _ = nullLayout(r)
	}
	dims := linesDimens(lines)
	for i := 0; i < len(lines)-1; i++ {
		// To avoid layout flickering while editing, assume a soft newline takes up all available space.
		if lay := lines[i].Layout; len(lay.Text) > 0 {
			r := lay.Text[len(lay.Text)-1]
			if r != '\n' {
				dims.Size.X = e.maxWidth
				break
			}
		}
	}
	return lines, dims
}

// CaretPos returns the line & column numbers of the caret.
func (e *Editor) CaretPos() (line, col int) {
	e.makeValid()
	return e.Caret.Line, e.Caret.Col
}

// CaretCoords returns the coordinates of the caret, relative to the
// editor itself.
func (e *Editor) CaretCoords() f32.Point {
	e.makeValid()
	return f32.Pt(float32(e.Caret.x)/64, float32(e.Caret.y))
}

func (e *Editor) layoutCaret() (line, col int, x fixed.Int26_6, y int) {
	var idx int
	var prevDesc fixed.Int26_6
loop:
	for {
		x = 0
		col = 0
		l := e.lines[line]
		y += (prevDesc + l.Ascent).Ceil()
		prevDesc = l.Descent
		for _, adv := range l.Layout.Advances {
			if idx == e.rr.caret {
				break loop
			}
			x += adv
			_, s := e.rr.runeAt(idx)
			idx += s
			col++
		}
		if line == len(e.lines)-1 || idx > e.rr.caret {
			break
		}
		line++
	}
	x += align(e.alignment, e.lines[line].Width, e.viewSize.X)
	return
}

func (e *Editor) invalidate() {
	e.valid = false
}

// Delete runes from the caret position. The sign of runes specifies the
// direction to delete: positive is forward, negative is backward.
func (e *Editor) Delete(runes int) {
	e.rr.deleteRunes(runes)
	e.Caret.xoff = 0
	e.invalidate()
}

// Insert inserts text at the caret, moving the caret forward.
func (e *Editor) Insert(s string) {
	e.append(s)
	e.Caret.scroll = true
	e.invalidate()
}

func (e *Editor) append(s string) {
	if e.singleLine {
		s = strings.ReplaceAll(s, "\n", "")
	}
	e.prepend(s)
	e.rr.caret += len(s)
}

func (e *Editor) prepend(s string) {
	e.rr.prepend(s)
	e.Caret.xoff = 0
	e.invalidate()
}

func (e *Editor) movePages(pages int) {
	e.makeValid()
	y := e.Caret.y + pages*e.viewSize.Y
	var (
		prevDesc fixed.Int26_6
		carLine2 int
	)
	y2 := e.lines[0].Ascent.Ceil()
	for i := 1; i < len(e.lines); i++ {
		if y2 >= y {
			break
		}
		l := e.lines[i]
		h := (prevDesc + l.Ascent).Ceil()
		prevDesc = l.Descent
		if y2+h-y >= y-y2 {
			break
		}
		y2 += h
		carLine2++
	}
	e.moveToLine(e.Caret.x+e.Caret.xoff, carLine2)
}

func (e *Editor) moveToLine(x fixed.Int26_6, line int) {
	e.makeValid()
	if line < 0 {
		line = 0
	}
	if line >= len(e.lines) {
		line = len(e.lines) - 1
	}
	
	prevDesc := e.lines[line].Descent
	for e.Caret.Line < line {
		e.moveEnd()
		l := e.lines[e.Caret.Line]
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
		e.Caret.y += (prevDesc + l.Ascent).Ceil()
		e.Caret.Col = 0
		prevDesc = l.Descent
		e.Caret.Line++
	}
	for e.Caret.Line > line {
		e.moveStart()
		l := e.lines[e.Caret.Line]
		_, s := e.rr.runeBefore(e.rr.caret)
		e.rr.caret -= s
		e.Caret.y -= (prevDesc + l.Ascent).Ceil()
		prevDesc = l.Descent
		e.Caret.Line--
		l = e.lines[e.Caret.Line]
		e.Caret.Col = len(l.Layout.Advances) - 1
	}
	
	e.moveStart()
	l := e.lines[line]
	e.Caret.x = align(e.alignment, l.Width, e.viewSize.X)
	// Only move past the end of the last line
	end := 0
	if line < len(e.lines)-1 {
		end = 1
	}
	// Move to rune closest to x.
	for i := 0; i < len(l.Layout.Advances)-end; i++ {
		adv := l.Layout.Advances[i]
		if e.Caret.x >= x {
			break
		}
		if e.Caret.x+adv-x >= x-e.Caret.x {
			break
		}
		e.Caret.x += adv
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
		e.Caret.Col++
	}
	e.Caret.xoff = x - e.Caret.x
}

// Move the caret: positive distance moves forward, negative distance moves
// backward.
func (e *Editor) Move(distance int) {
	e.makeValid()
	for ; distance < 0 && e.rr.caret > 0; distance++ {
		if e.Caret.Col == 0 {
			// Move to end of previous line.
			e.moveToLine(fixed.I(e.maxWidth), e.Caret.Line-1)
			continue
		}
		l := e.lines[e.Caret.Line].Layout
		_, s := e.rr.runeBefore(e.rr.caret)
		e.rr.caret -= s
		e.Caret.Col--
		e.Caret.x -= l.Advances[e.Caret.Col]
	}
	for ; distance > 0 && e.rr.caret < e.rr.len(); distance-- {
		l := e.lines[e.Caret.Line].Layout
		// Only move past the end of the last line
		end := 0
		if e.Caret.Line < len(e.lines)-1 {
			end = 1
		}
		if e.Caret.Col >= len(l.Advances)-end {
			// Move to start of next line.
			e.moveToLine(0, e.Caret.Line+1)
			continue
		}
		e.Caret.x += l.Advances[e.Caret.Col]
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
		e.Caret.Col++
	}
	e.Caret.xoff = 0
}

func (e *Editor) moveStart() {
	e.makeValid()
	lo := e.lines[e.Caret.Line].Layout
	for i := e.Caret.Col - 1; i >= 0; i-- {
		_, s := e.rr.runeBefore(e.rr.caret)
		e.rr.caret -= s
		e.Caret.x -= lo.Advances[i]
	}
	e.Caret.Col = 0
	e.Caret.xoff = -e.Caret.x
}

func (e *Editor) moveEnd() {
	e.makeValid()
	l := e.lines[e.Caret.Line]
	// Only move past the end of the last line
	end := 0
	if e.Caret.Line < len(e.lines)-1 {
		end = 1
	}
	lo := l.Layout
	for i := e.Caret.Col; i < len(lo.Advances)-end; i++ {
		adv := lo.Advances[i]
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
		e.Caret.x += adv
		e.Caret.Col++
	}
	a := align(e.alignment, l.Width, e.viewSize.X)
	e.Caret.xoff = l.Width + a - e.Caret.x
}

// moveWord moves the caret to the next word in the specified direction.
// Positive is forward, negative is backward.
// Absolute values greater than one will skip that many words.
func (e *Editor) moveWord(distance int) {
	e.makeValid()
	// split the distance information into constituent parts to be
	// used independently.
	words, direction := distance, 1
	if distance < 0 {
		words, direction = distance*-1, -1
	}
	// atEnd if caret is at either side of the buffer.
	atEnd := func() bool {
		return e.rr.caret == 0 || e.rr.caret == e.rr.len()
	}
	// next returns the appropriate rune given the direction.
	next := func() (r rune) {
		if direction < 0 {
			r, _ = e.rr.runeBefore(e.rr.caret)
		} else {
			r, _ = e.rr.runeAt(e.rr.caret)
		}
		return r
	}
	for ii := 0; ii < words; ii++ {
		for r := next(); unicode.IsSpace(r) && !atEnd(); r = next() {
			e.Move(direction)
		}
		e.Move(direction)
		for r := next(); !unicode.IsSpace(r) && !atEnd(); r = next() {
			e.Move(direction)
		}
	}
}

// deleteWord the next word(s) in the specified direction. Unlike moveWord, deleteWord treats whitespace as a word
// itself.
//
// Positive is forward, negative is backward.
//
// Absolute values greater than one will delete that many words.
func (e *Editor) deleteWord(distance int) {
	e.makeValid()
	// split the distance information into constituent parts to be
	// used independently.
	words, direction := distance, 1
	if distance < 0 {
		words, direction = distance*-1, -1
	}
	// atEnd if offset is at or beyond either side of the buffer.
	atEnd := func(offset int) bool {
		idx := e.rr.caret + offset*direction
		return idx <= 0 || idx >= e.rr.len()
	}
	// next returns the appropriate rune given the direction and offset.
	next := func(offset int) (r rune) {
		idx := e.rr.caret + offset*direction
		if idx < 0 {
			idx = 0
		} else if idx > e.rr.len() {
			idx = e.rr.len()
		}
		if direction < 0 {
			r, _ = e.rr.runeBefore(idx)
		} else {
			r, _ = e.rr.runeAt(idx)
		}
		return r
	}
	var runes = 1
	for ii := 0; ii < words; ii++ {
		if r := next(runes); unicode.IsSpace(r) {
			for r := next(runes); unicode.IsSpace(r) && !atEnd(runes); r = next(runes) {
				runes += 1
			}
		} else {
			for r := next(runes); !unicode.IsSpace(r) && !atEnd(runes); r = next(runes) {
				runes += 1
			}
		}
	}
	e.Delete(runes * direction)
}

func (e *Editor) scrollToCaret() {
	e.makeValid()
	l := e.lines[e.Caret.Line]
	if e.singleLine {
		var dist int
		if d := e.Caret.x.Floor() - e.scrollOff.X; d < 0 {
			dist = d
		} else if d := e.Caret.x.Ceil() - (e.scrollOff.X + e.viewSize.X); d > 0 {
			dist = d
		}
		e.scrollRel(dist, 0)
	} else {
		miny := e.Caret.y - l.Ascent.Ceil()
		maxy := e.Caret.y + l.Descent.Ceil()
		var dist int
		if d := miny - e.scrollOff.Y; d < 0 {
			dist = d
		} else if d := maxy - (e.scrollOff.Y + e.viewSize.Y); d > 0 {
			dist = d
		}
		e.scrollRel(0, dist)
	}
}

// NumLines returns the number of lines in the editor.
func (e *Editor) NumLines() int {
	e.makeValid()
	return len(e.lines)
}

func nullLayout(r io.Reader) ([]text.Line, error) {
	rr := bufio.NewReader(r)
	var rerr error
	var n int
	var buf bytes.Buffer
	for {
		r, s, err := rr.ReadRune()
		n += s
		buf.WriteRune(r)
		if err != nil {
			rerr = err
			break
		}
	}
	return []text.Line{
		{
			Layout: text.Layout{
				Text:     buf.String(),
				Advances: make([]fixed.Int26_6, n),
			},
		},
	}, rerr
}

func (s ChangeEvent) isEditorEvent() {}
func (s SubmitEvent) isEditorEvent() {}
