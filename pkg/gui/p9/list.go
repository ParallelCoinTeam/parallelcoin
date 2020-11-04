package p9

import (
	"image"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

// List displays a subsection of a potentially infinitely large underlying list. List accepts user input to scroll the
// subsection.
type List struct {
	th   *Theme
	axis l.Axis
	// ScrollToEnd instructs the list to stay scrolled to the far end position once reached. A List with ScrollToEnd ==
	// true and Position.BeforeEnd == false draws its content with the last item at the bottom of the list area.
	scrollToEnd bool
	// Alignment is the cross axis alignment of list elements.
	alignment     l.Alignment
	disableScroll bool
	ctx           l.Context
	scroll        gesture.Scroll
	sideScroll    gesture.Scroll
	scrollDelta   int
	// Position is updated during Layout. To save the list scroll position, just save Position after Layout finishes. To
	// scroll the list programmatically, update Position (e.g. restore it from a saved value) before calling Layout.
	position Position
	// nextUp, nextDown Position
	len             int
	drag            gesture.Drag
	color           string
	active          string
	background      string
	currentColor    string
	scrollWidth     int
	setScrollWidth  int
	scrollBarPad    int
	setScrollBarPad int
	// maxSize is the total size of visible children.
	maxSize             int
	children            []scrollChild
	dir                 iterationDir
	length              int
	prevLength          int
	w                   ListElement
	pageUp, pageDown    *Clickable
	dims                DimensionList
	view, total, before int
	top, middle, bottom int
	lastWidth           int
	recalculateTime     time.Time
	recalculate         bool
	notFirst            bool
}

// List returns a new scrollable List widget
func (th *Theme) List() (li *List) {
	li = &List{
		th:              th,
		pageUp:          th.Clickable(),
		pageDown:        th.Clickable(),
		color:           "Primary",
		background:      "Transparent",
		active:          "DocBg",
		scrollWidth:     int(th.TextSize.Scale(1).V),
		setScrollWidth:  int(th.TextSize.Scale(1).V),
		scrollBarPad:    int(th.TextSize.Scale(0.25).V),
		setScrollBarPad: int(th.TextSize.Scale(0.25).V),
		recalculateTime: time.Now().Add(-time.Second),
		recalculate:     true,
	}
	li.currentColor = li.color
	return
}

// Vertical sets the axis to vertical (default implicit is horizontal)
func (li *List) Vertical() (out *List) {
	li.axis = l.Vertical
	return li
}

func (li *List) Start() *List {
	li.alignment = l.Start
	return li
}

func (li *List) End() *List {
	li.alignment = l.End
	return li
}

func (li *List) Middle() *List {
	li.alignment = l.Middle
	return li
}

func (li *List) Baseline() *List {
	li.alignment = l.Baseline
	return li
}

// ScrollToEnd sets the List to add new items to the end and push older ones up/left and initial render has scroll
// to the end (or bottom) of the List
func (li *List) ScrollToEnd() (out *List) {
	li.scrollToEnd = true
	return li
}

func (li *List) Length(length int) *List {
	li.prevLength = li.length
	li.length = length
	return li
}

func (li *List) DisableScroll(disable bool) *List {
	li.disableScroll = disable
	if disable {
		li.scrollWidth = 0
		li.scrollBarPad = 0
	} else {
		li.scrollWidth = li.setScrollWidth
		li.scrollBarPad = li.setScrollBarPad
	}
	return li
}

func (li *List) ListElement(w ListElement) *List {
	li.w = w
	return li
}

func (li *List) ScrollWidth(width int) *List {
	li.scrollWidth = width
	li.setScrollWidth = width
	return li
}

func (li *List) Color(color string) *List {
	li.color = color
	li.currentColor = li.color
	return li
}

func (li *List) Background(color string) *List {
	li.background = color
	return li
}

func (li *List) Active(color string) *List {
	li.active = color
	return li
}

func (li *List) Slice(gtx l.Context, widgets ...l.Widget) l.Widget {
	return li.Length(len(widgets)).Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
		return widgets[index](gtx)
	}).Fn
}

// Fn runs the layout in the configured context. The ListElement function returns the widget at the given index
func (li *List) Fn(gtx l.Context) l.Dimensions {
	if li.length == 0 {
		// if there is no children just return a big empty box
		return EmptyFromSize(gtx.Constraints.Max)(gtx)
	}
	if li.disableScroll {
		return li.embedWidget(0)(gtx)
	}
	if li.length != li.prevLength {
		li.recalculate = true
		li.recalculateTime = time.Now().Add(time.Millisecond * 100)
	} else if li.lastWidth != gtx.Constraints.Max.X && li.notFirst {
		li.recalculateTime = time.Now().Add(time.Millisecond * 100)
		li.recalculate = true
	}
	if !li.notFirst {
		li.recalculateTime = time.Now().Add(-time.Millisecond * 100)
		li.notFirst = true
	}
	li.lastWidth = gtx.Constraints.Max.X
	if li.recalculateTime.Sub(time.Now()) < 0 && li.recalculate {
		// return li.embedWidget(li.scrollWidth)(gtx)
		// } else {
		// if li.recalculate && !li.changing {
		Debug("recalculating")
		// get the size of the scrollbar
		li.th.scrollBarSize = li.scrollWidth + li.scrollBarPad
		// render the widgets onto a second context to get their dimensions
		gtx1 := CopyContextDimensions(gtx, gtx.Constraints.Max, li.axis)
		// generate the dimensions for all the list elements
		li.dims = GetDimensionList(gtx1, li.length, li.w)
		// li.recalculate = false
		li.recalculateTime = time.Time{}
		li.recalculate = false
	}
	_, li.view = axisMainConstraint(li.axis, gtx.Constraints)
	// Debugs(li.dims)
	li.total, li.before = li.dims.GetSizes(li.position, li.axis)
	if li.total == 0 {
		// if there is no children just return a big empty box
		return EmptyFromSize(gtx.Constraints.Max)(gtx)
	}
	if li.total < li.view {
		// Debug("not showing scrollbar", li.total, li.view)
		// if the contents fit the view, don't show the scrollbar
		li.top, li.middle, li.bottom = 0, 0, 0
		li.scrollWidth = 0
		li.scrollBarPad = 0
	} else {
		li.scrollWidth = li.setScrollWidth
		li.scrollBarPad = li.setScrollBarPad
		li.top = li.before * (li.view - li.scrollWidth) / li.total
		li.middle = li.view * (li.view - li.scrollWidth) / li.total
		li.bottom = (li.total - li.before - li.view) * (li.view - li.scrollWidth) / li.total
		if li.view < li.scrollWidth {
			li.middle = li.view
			li.top, li.bottom = 0, 0
		} else {
			li.middle += li.scrollWidth
		}
	}
	// now lay it all out and draw the list and scrollbar
	var container l.Widget
	if li.axis == l.Horizontal {
		container = li.th.VFlex().
			Rigid(li.embedWidget(li.scrollWidth + li.scrollBarPad*2)).
			Rigid(
				li.th.VFlex().
					Rigid(
						li.th.Fill(li.background, EmptySpace(0, li.scrollBarPad*2)).Fn,
					).
					Rigid(
						li.th.Fill(li.background,
							li.th.Flex().
								Rigid(li.pageUpDown(li.dims, li.view, li.total, li.scrollWidth, li.top, false)).
								Rigid(li.grabber(li.dims, li.scrollWidth, li.middle)).
								Rigid(li.pageUpDown(li.dims, li.view, li.total, li.scrollWidth, li.bottom, true)).
								Fn,
						).Fn,
					).
					Fn,
			).Fn
	} else {
		container = li.th.Flex().
			Rigid(li.embedWidget(li.scrollWidth + li.scrollBarPad*2)).
			Rigid(
				li.th.Fill(li.background,
					li.th.Flex().
						Rigid(
							EmptySpace(li.scrollBarPad*2, 0),
						).
						Rigid(
							li.th.Flex().Vertical().
								Rigid(li.pageUpDown(li.dims, li.view, li.total, li.scrollWidth, li.top, false)).
								Rigid(li.grabber(li.dims, li.scrollWidth, li.middle)).
								Rigid(li.pageUpDown(li.dims, li.view, li.total, li.scrollWidth, li.bottom, true)).
								Fn,
						).Fn,
				).Fn,
			).Fn
	}
	return container(gtx)
}

func (li *List) embedWidget(scrollWidth int) func(l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		if li.axis == l.Horizontal {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y - scrollWidth // - li.scrollBarPad
			gtx.Constraints.Max.Y = gtx.Constraints.Min.Y
		} else {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X - scrollWidth // - li.scrollBarPad
			gtx.Constraints.Max.X = gtx.Constraints.Min.X
		}
		return li.Layout(gtx, li.length, li.w)
	}
}

func (li *List) pageUpDown(dims DimensionList, view, total, x, y int, down bool) func(l.Context) l.Dimensions {
	button := li.pageUp
	if down {
		button = li.pageDown
	}
	return func(gtx l.Context) l.Dimensions {
		pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
		li.sideScroll.Add(gtx.Ops)
		return li.th.ButtonLayout(button.SetClick(func() {
			current := dims.PositionToCoordinate(li.position, li.axis)
			var newPos int
			if down {
				if current+view > total {
					newPos = total - view
				} else {
					newPos = current + view
				}
			} else {
				newPos = current - view
				if newPos < 0 {
					newPos = 0
				}
			}
			li.position = dims.CoordinateToPosition(newPos, li.axis)
		})).Embed(
			li.th.Fill(li.background,
				EmptySpace(x, y),
			).Fn,
		).Background(li.background).CornerRadius(0).Fn(gtx)
	}
}

func (li *List) grabber(dims DimensionList, x, y int) func(l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		ax := gesture.Vertical
		if li.axis == l.Horizontal {
			ax = gesture.Horizontal
		}
		var de *pointer.Event
		for _, ev := range li.drag.Events(gtx.Metric, gtx, ax) {
			if ev.Type == pointer.Press ||
				ev.Type == pointer.Release ||
				ev.Type == pointer.Drag {
				de = &ev
			}
		}
		if de != nil {
			// respond to the event
			if de.Type == pointer.Press || de.Type == pointer.Drag {
				li.currentColor = li.active
			}
			if de.Type == pointer.Release {
				li.currentColor = li.color
			}
			if de.Type == pointer.Drag {
				current := dims.PositionToCoordinate(li.position, li.axis)
				var d int
				if li.axis == l.Horizontal {
					d = int(de.Position.X) + current
				} else {
					d = int(de.Position.Y) + current
				}
				total := dims.GetTotal(gtx, li.axis)
				if d > total {
					d = total - 1
				}
				li.position = dims.CoordinateToPosition(d, li.axis)
			}
			// if de.Type == pointer.Scroll {
		}
		defer op.Push(gtx.Ops).Pop()
		pointer.Rect(image.Rectangle{Max: image.Point{X: x, Y: y}}).Add(gtx.Ops)
		li.drag.Add(gtx.Ops)
		pointer.Rect(image.Rectangle{Max: image.Point{X: x, Y: y}}).Add(gtx.Ops)
		li.sideScroll.Add(gtx.Ops)
		return li.th.Fill(li.currentColor,
			EmptySpace(x, y),
		).Fn(gtx)
	}
}

type scrollChild struct {
	size image.Point
	call op.CallOp
}

// ListElement is a function that computes the dimensions of a list element.
type ListElement func(gtx l.Context, index int) l.Dimensions

type iterationDir uint8

// Position is a List scroll offset represented as an offset from the top edge of a child element.
type Position struct {
	// BeforeEnd tracks whether the List position is before the very end. We use "before end" instead of "at end" so
	// that the zero value of a Position struct is useful.
	//
	// When laying out a list, if ScrollToEnd is true and BeforeEnd is false, then First and Offset are ignored, and the
	// list is drawn with the last item at the bottom. If ScrollToEnd is false then BeforeEnd is ignored.
	BeforeEnd bool
	// First is the index of the first visible child.
	First int
	// Offset is the distance in pixels from the top edge to the child at index First.
	Offset int
}

const (
	iterateNone iterationDir = iota
	iterateForward
	iterateBackward
)

// init prepares the list for iterating through its children with next.
func (li *List) init(gtx l.Context, len int) {
	if li.more() {
		panic("unfinished child")
	}
	li.ctx = gtx
	li.maxSize = 0
	li.children = li.children[:0]
	li.len = len
	li.update()
	if li.canScrollToEnd() || li.position.First > len {
		li.position.Offset = 0
		li.position.First = len
	}
}

// Layout the List.
func (li *List) Layout(gtx l.Context, len int, w ListElement) l.Dimensions {
	li.init(gtx, len)
	crossMin, crossMax := axisCrossConstraint(li.axis, gtx.Constraints)
	gtx.Constraints = axisConstraints(li.axis, 0, Inf, crossMin, crossMax)
	macro := op.Record(gtx.Ops)
	for li.next(); li.more(); li.next() {
		child := op.Record(gtx.Ops)
		dims := w(gtx, li.index())
		call := child.Stop()
		li.end(dims, call)
	}
	return li.layout(macro)
}

func (li *List) canScrollToEnd() bool {
	return li.scrollToEnd && !li.position.BeforeEnd
}

// Dragging reports whether the List is being dragged.
func (li *List) Dragging() bool {
	return li.scroll.State() == gesture.StateDragging
}

func (li *List) update() {
	d := li.scroll.Scroll(li.ctx.Metric, li.ctx, li.ctx.Now, gesture.Axis(li.axis))
	d += li.sideScroll.Scroll(li.ctx.Metric, li.ctx, li.ctx.Now, gesture.Axis(li.axis))
	li.scrollDelta = d
	li.position.Offset += d
}

// next advances to the next child.
func (li *List) next() {
	li.dir = li.nextDir()
	// The user scroll offset is applied after scrolling to
	// list end.
	if li.canScrollToEnd() && !li.more() && li.scrollDelta < 0 {
		li.position.BeforeEnd = true
		li.position.Offset += li.scrollDelta
		li.dir = li.nextDir()
	}
}

// index is current child's position in the underlying list.
func (li *List) index() int {
	switch li.dir {
	case iterateBackward:
		return li.position.First - 1
	case iterateForward:
		return li.position.First + len(li.children)
	default:
		panic("Index called before Next")
	}
}

// more reports whether more children are needed.
func (li *List) more() bool {
	return li.dir != iterateNone
}

func (li *List) nextDir() iterationDir {
	_, vsize := axisMainConstraint(li.axis, li.ctx.Constraints)
	last := li.position.First + len(li.children)
	// Clamp offset.
	if li.maxSize-li.position.Offset < vsize && last == li.len {
		li.position.Offset = li.maxSize - vsize
	}
	if li.position.Offset < 0 && li.position.First == 0 {
		li.position.Offset = 0
	}
	switch {
	case len(li.children) == li.len:
		return iterateNone
	case li.maxSize-li.position.Offset < vsize:
		return iterateForward
	case li.position.Offset < 0:
		return iterateBackward
	}
	return iterateNone
}

// End the current child by specifying its dimensions.
func (li *List) end(dims l.Dimensions, call op.CallOp) {
	child := scrollChild{dims.Size, call}
	mainSize := axisMain(li.axis, child.size)
	li.maxSize += mainSize
	switch li.dir {
	case iterateForward:
		li.children = append(li.children, child)
	case iterateBackward:
		li.children = append([]scrollChild{child}, li.children...)
		li.position.First--
		li.position.Offset += mainSize
	default:
		panic("call Next before End")
	}
	li.dir = iterateNone
}

// Layout the List and return its dimensions.
func (li *List) layout(macro op.MacroOp) l.Dimensions {
	if li.more() {
		panic("unfinished child")
	}
	mainMin, mainMax := axisMainConstraint(li.axis, li.ctx.Constraints)
	children := li.children
	// Skip invisible children
	for len(children) > 0 {
		sz := children[0].size
		mainSize := axisMain(li.axis, sz)
		if li.position.Offset <= mainSize {
			break
		}
		li.position.First++
		li.position.Offset -= mainSize
		children = children[1:]
	}
	size := -li.position.Offset
	var maxCross int
	for i, child := range children {
		sz := child.size
		if c := axisCross(li.axis, sz); c > maxCross {
			maxCross = c
		}
		size += axisMain(li.axis, sz)
		if size >= mainMax {
			children = children[:i+1]
			break
		}
	}
	ops := li.ctx.Ops
	pos := -li.position.Offset
	// ScrollToEnd lists are end aligned.
	if space := mainMax - size; li.scrollToEnd && space > 0 {
		pos += space
	}
	for _, child := range children {
		sz := child.size
		var cross int
		switch li.alignment {
		case l.End:
			cross = maxCross - axisCross(li.axis, sz)
		case l.Middle:
			cross = (maxCross - axisCross(li.axis, sz)) / 2
		}
		childSize := axisMain(li.axis, sz)
		max := childSize + pos
		if max > mainMax {
			max = mainMax
		}
		min := pos
		if min < 0 {
			min = 0
		}
		r := image.Rectangle{
			Min: axisPoint(li.axis, min, -Inf),
			Max: axisPoint(li.axis, max, Inf),
		}
		stack := op.Push(ops)
		clip.Rect(r).Add(ops)
		op.Offset(toPointF(axisPoint(li.axis, pos, cross))).Add(ops)
		child.call.Add(ops)
		stack.Pop()
		pos += childSize
	}
	atStart := li.position.First == 0 && li.position.Offset <= 0
	atEnd := li.position.First+len(children) == li.len && mainMax >= pos
	if atStart && li.scrollDelta < 0 || atEnd && li.scrollDelta > 0 {
		li.scroll.Stop()
	}
	li.position.BeforeEnd = !atEnd
	if pos < mainMin {
		pos = mainMin
	}
	if pos > mainMax {
		pos = mainMax
	}
	dims := axisPoint(li.axis, pos, maxCross)
	call := macro.Stop()
	defer op.Push(li.ctx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: dims}).Add(ops)
	li.scroll.Add(ops)
	call.Add(ops)
	return l.Dimensions{Size: dims}
}
