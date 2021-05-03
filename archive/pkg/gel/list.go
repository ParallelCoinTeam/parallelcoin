package gel

import (
	"image"
	"time"
	
	"github.com/p9c/gio/gesture"
	"github.com/p9c/gio/io/pointer"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/clip"
)

type scrollChild struct {
	size image.Point
	call op.CallOp
}

// List displays a subsection of a potentially infinitely large underlying list. List accepts user input to scroll the
// subsection.
type List struct {
	axis l.Axis
	// ScrollToEnd instructs the list to stay scrolled to the far end position once reached. A List with ScrollToEnd ==
	// true and Position.BeforeEnd == false draws its content with the last item at the bottom of the list area.
	scrollToEnd bool
	// Alignment is the cross axis alignment of list elements.
	alignment   l.Alignment
	scroll      gesture.Scroll
	scrollDelta int
	// position is updated during Layout. To save the list scroll position, just save Position after Layout finishes. To
	// scroll the list programmatically, update Position (e.g. restore it from a saved value) before calling Layout.
	// nextUp, nextDown Position
	position Position
	Len      int
	// maxSize is the total size of visible children.
	maxSize  int
	children []scrollChild
	dir      iterationDir
	
	// all below are additional fields to implement the scrollbar
	*Window
	// we store the constraints here instead of in the `cs` field
	ctx                 l.Context
	sideScroll          gesture.Scroll
	disableScroll       bool
	drag                gesture.Drag
	recentPageClick     time.Time
	color               string
	active              string
	background          string
	currentColor        string
	scrollWidth         int
	setScrollWidth      int
	length              int
	prevLength          int
	w                   ListElement
	pageUp, pageDown    *Clickable
	dims                DimensionList
	cross               int
	view, total, before int
	top, middle, bottom int
	lastWidth           int
	recalculateTime     time.Time
	recalculate         bool
	notFirst            bool
	leftSide            bool
}

// List returns a new scrollable List widget
func (w *Window) List() (li *List) {
	li = &List{
		Window:          w,
		pageUp:          w.WidgetPool.GetClickable(),
		pageDown:        w.WidgetPool.GetClickable(),
		color:           "DocText",
		background:      "Transparent",
		active:          "Primary",
		scrollWidth:     int(w.TextSize.Scale(0.75).V),
		setScrollWidth:  int(w.TextSize.Scale(0.75).V),
		recalculateTime: time.Now().Add(-time.Second),
		recalculate:     true,
	}
	li.currentColor = li.color
	return
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
	// OffsetLast is the signed distance in pixels from the bottom edge to the
	// bottom edge of the child at index First+Count.
	OffsetLast int
	// Count is the number of visible children.
	Count int
}

const (
	iterateNone iterationDir = iota
	iterateForward
	iterateBackward
)

// init prepares the list for iterating through its children with next.
func (li *List) init(gtx l.Context, length int) {
	if li.more() {
		panic("unfinished child")
	}
	li.ctx = gtx
	li.maxSize = 0
	li.children = li.children[:0]
	li.Len = length
	li.update()
	if li.canScrollToEnd() || li.position.First > length {
		li.position.Offset = 0
		li.position.First = length
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

// canScrollToEnd returns true if there is room to scroll further towards the end
func (li *List) canScrollToEnd() bool {
	return li.scrollToEnd && !li.position.BeforeEnd
}

// Dragging reports whether the List is being dragged.
func (li *List) Dragging() bool {
	return li.scroll.State() == gesture.StateDragging ||
		li.sideScroll.State() == gesture.StateDragging
}

// update the scrolling
func (li *List) update() {
	d := li.scroll.Scroll(li.ctx.Metric, li.ctx, li.ctx.Now, gesture.Axis(li.axis))
	d += li.sideScroll.Scroll(li.ctx.Metric, li.ctx, li.ctx.Now, gesture.Axis(li.axis))
	li.scrollDelta = d
	li.position.Offset += d
}

// next advances to the next child.
func (li *List) next() {
	li.dir = li.nextDir()
	// The user scroll offset is applied after scrolling to list end.
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
	_, vSize := axisMainConstraint(li.axis, li.ctx.Constraints)
	last := li.position.First + len(li.children)
	// Clamp offset.
	if li.maxSize-li.position.Offset < vSize && last == li.Len {
		li.position.Offset = li.maxSize - vSize
	}
	if li.position.Offset < 0 && li.position.First == 0 {
		li.position.Offset = 0
	}
	switch {
	case len(li.children) == li.Len:
		return iterateNone
	case li.maxSize-li.position.Offset < vSize:
		return iterateForward
	case li.position.Offset < 0:
		return iterateBackward
	}
	return iterateNone
}

// End the current child by specifying its dimensions.
func (li *List) end(dims l.Dimensions, call op.CallOp) {
	child := scrollChild{dims.Size, call}
	mainSize := axisConvert(li.axis, child.size).X
	li.maxSize += mainSize
	switch li.dir {
	case iterateForward:
		li.children = append(li.children, child)
	case iterateBackward:
		li.children = append(li.children, scrollChild{})
		copy(li.children[1:], li.children)
		li.children[0] = child
		li.position.First--
		li.position.Offset += mainSize
	default:
		panic("call Next before End")
	}
	li.dir = iterateNone
}

// layout the List and return its dimensions.
func (li *List) layout(macro op.MacroOp) l.Dimensions {
	if li.more() {
		panic("unfinished child")
	}
	mainMin, mainMax := axisMainConstraint(li.axis, li.ctx.Constraints)
	children := li.children
	// Skip invisible children
	for len(children) > 0 {
		sz := children[0].size
		mainSize := axisConvert(li.axis, sz).X
		if li.position.Offset < mainSize {
			// First child is partially visible.
			break
		}
		li.position.First++
		li.position.Offset -= mainSize
		children = children[1:]
	}
	size := -li.position.Offset
	var maxCross int
	for i, child := range children {
		sz := axisConvert(li.axis, child.size)
		if c := sz.Y; c > maxCross {
			maxCross = c
		}
		size += sz.X
		if size >= mainMax {
			children = children[:i+1]
			break
		}
	}
	li.position.Count = len(children)
	li.position.OffsetLast = mainMax - size
	ops := li.ctx.Ops
	pos := -li.position.Offset
	// ScrollToEnd lists are end aligned.
	if space := li.position.OffsetLast; li.scrollToEnd && space > 0 {
		pos += space
	}
	for _, child := range children {
		sz := axisConvert(li.axis, child.size)
		var cross int
		switch li.alignment {
		case l.End:
			cross = maxCross - sz.Y
		case l.Middle:
			cross = (maxCross - sz.Y) / 2
		}
		childSize := sz.X
		max := childSize + pos
		if max > mainMax {
			max = mainMax
		}
		min := pos
		if min < 0 {
			min = 0
		}
		r := image.Rectangle{
			Min: axisConvert(li.axis, image.Pt(min, -Inf)),
			Max: axisConvert(li.axis, image.Pt(max, Inf)),
		}
		stack := op.Save(ops)
		clip.Rect(r).Add(ops)
		pt := axisConvert(li.axis, image.Pt(pos, cross))
		op.Offset(Fpt(pt)).Add(ops)
		child.call.Add(ops)
		stack.Load()
		pos += childSize
	}
	atStart := li.position.First == 0 && li.position.Offset <= 0
	atEnd := li.position.First+len(children) == li.Len && mainMax >= pos
	if atStart && li.scrollDelta < 0 || atEnd && li.scrollDelta > 0 {
		li.scroll.Stop()
		li.sideScroll.Stop()
	}
	li.position.BeforeEnd = !atEnd
	if pos < mainMin {
		pos = mainMin
	}
	if pos > mainMax {
		pos = mainMax
	}
	dims := axisConvert(li.axis, image.Pt(pos, maxCross))
	call := macro.Stop()
	defer op.Save(ops).Load()
	bounds := image.Rectangle{Max: dims}
	pointer.Rect(bounds).Add(ops)
	// li.sideScroll.Add(ops, bounds)
	// li.scroll.Add(ops, bounds)
	
	var min, max int
	if o := li.position.Offset; o > 0 {
		// Use the size of the invisible part as scroll boundary.
		min = -o
	} else if li.position.First > 0 {
		min = -Inf
	}
	if o := li.position.OffsetLast; o < 0 {
		max = -o
	} else if li.position.First+li.position.Count < li.Len {
		max = Inf
	}
	scrollRange := image.Rectangle{
		Min: axisConvert(li.axis, image.Pt(min, 0)),
		Max: axisConvert(li.axis, image.Pt(max, 0)),
	}
	li.scroll.Add(ops, scrollRange)
	li.sideScroll.Add(ops, scrollRange)
	
	call.Add(ops)
	return l.Dimensions{Size: dims}
}

// Everything below is extensions on the original from gioui.org/layout

// Position returns the current position of the scroller
func (li *List) Position() Position {
	return li.position
}

// SetPosition sets the position of the scroller
func (li *List) SetPosition(position Position) {
	li.position = position
}

// JumpToStart moves the position to the start
func (li *List) JumpToStart() {
	li.position = Position{}
}

// JumpToEnd moves the position to the end
func (li *List) JumpToEnd() {
	li.position = Position{
		BeforeEnd: false,
		First:     len(li.dims),
		Offset:    axisMain(li.axis, li.dims[len(li.dims)-1].Size),
	}
}

// Vertical sets the axis to vertical (default implicit is horizontal)
func (li *List) Vertical() (out *List) {
	li.axis = l.Vertical
	return li
}

// Start sets the alignment to start
func (li *List) Start() *List {
	li.alignment = l.Start
	return li
}

// End sets the alignment to end
func (li *List) End() *List {
	li.alignment = l.End
	return li
}

// Middle sets the alignment to middle
func (li *List) Middle() *List {
	li.alignment = l.Middle
	return li
}

// Baseline sets the alignment to baseline
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

// LeftSide sets the scroller to be on the opposite side from usual
func (li *List) LeftSide(b bool) (out *List) {
	li.leftSide = b
	return li
}

// Length sets the new length for the list
func (li *List) Length(length int) *List {
	li.prevLength = li.length
	li.length = length
	return li
}

// DisableScroll turns off the scrollbar
func (li *List) DisableScroll(disable bool) *List {
	li.disableScroll = disable
	if disable {
		li.scrollWidth = 0
	} else {
		li.scrollWidth = li.setScrollWidth
	}
	return li
}

// ListElement defines the function that returns list elements
func (li *List) ListElement(w ListElement) *List {
	li.w = w
	return li
}

// ScrollWidth sets the width of the scrollbar
func (li *List) ScrollWidth(width int) *List {
	li.scrollWidth = width
	li.setScrollWidth = width
	return li
}

// Color sets the primary color of the scrollbar grabber
func (li *List) Color(color string) *List {
	li.color = color
	li.currentColor = li.color
	return li
}

// Background sets the background color of the scrollbar
func (li *List) Background(color string) *List {
	li.background = color
	return li
}

// Active sets the color of the scrollbar grabber when it is being operated
func (li *List) Active(color string) *List {
	li.active = color
	return li
}

func (li *List) Slice(gtx l.Context, widgets ...l.Widget) l.Widget {
	return li.Length(len(widgets)).Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
		return widgets[index](gtx)
	},
	).Fn
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
		li.scrollBarSize = li.scrollWidth // + li.scrollBarPad
		gtx1 := CopyContextDimensionsWithMaxAxis(gtx, li.axis)
		// generate the dimensions for all the list elements
		li.dims = GetDimensionList(gtx1, li.length, li.w)
		li.recalculateTime = time.Time{}
		li.recalculate = false
	}
	_, li.view = axisMainConstraint(li.axis, gtx.Constraints)
	_, li.cross = axisCrossConstraint(li.axis, gtx.Constraints)
	li.total, li.before = li.dims.GetSizes(li.position, li.axis)
	if li.total == 0 {
		// if there is no children just return a big empty box
		return EmptyFromSize(gtx.Constraints.Max)(gtx)
	}
	if li.total < li.view {
		// if the contents fit the view, don't show the scrollbar
		li.top, li.middle, li.bottom = 0, 0, 0
		li.scrollWidth = 0
	} else {
		li.scrollWidth = li.setScrollWidth
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
		containerFlex := li.Theme.VFlex()
		if !li.leftSide {
			containerFlex.Rigid(li.embedWidget(li.scrollWidth /* + int(li.TextSize.True)/4)*/))
			containerFlex.Rigid(EmptySpace(int(li.TextSize.V)/4, int(li.TextSize.V)/4))
		}
		containerFlex.Rigid(
			li.VFlex().
				Rigid(
					func(gtx l.Context) l.Dimensions {
						pointer.Rect(image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X,
							Y: gtx.Constraints.Max.Y,
						},
						},
						).Add(gtx.Ops)
						li.drag.Add(gtx.Ops)
						return li.Theme.Flex().
							Rigid(li.pageUpDown(li.dims, li.view, li.total,
								// li.scrollBarPad+
								li.scrollWidth, li.top, false,
							),
							).
							Rigid(li.grabber(li.dims, li.scrollWidth, li.middle,
								li.view, gtx.Constraints.Max.X,
							),
							).
							Rigid(li.pageUpDown(li.dims, li.view, li.total,
								// li.scrollBarPad+
								li.scrollWidth, li.bottom, true,
							),
							).
							Fn(gtx)
					},
				).
				Fn,
		)
		if li.leftSide {
			containerFlex.Rigid(EmptySpace(int(li.TextSize.V)/4, int(li.TextSize.V)/4))
			containerFlex.Rigid(li.embedWidget(li.scrollWidth)) // li.scrollWidth)) // + li.scrollBarPad))
		}
		container = containerFlex.Fn
	} else {
		containerFlex := li.Theme.Flex()
		if !li.leftSide {
			containerFlex.Rigid(li.embedWidget(li.scrollWidth + int(li.TextSize.V)/2)) // + li.scrollBarPad))
			containerFlex.Rigid(EmptySpace(int(li.TextSize.V)/2, int(li.TextSize.V)/2))
		}
		containerFlex.Rigid(
			li.Fill(li.background, l.Center, li.TextSize.V/4, 0, li.Flex().
				Rigid(
					func(gtx l.Context) l.Dimensions {
						pointer.Rect(image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X,
							Y: gtx.Constraints.Max.Y,
						},
						},
						).Add(gtx.Ops)
						li.drag.Add(gtx.Ops)
						return li.Theme.Flex().Vertical().
							Rigid(li.pageUpDown(li.dims, li.view, li.total,
								li.scrollWidth, li.top, false,
							),
							).
							Rigid(li.grabber(li.dims,
								li.scrollWidth, li.middle,
								li.view, gtx.Constraints.Max.X,
							),
							).
							Rigid(li.pageUpDown(li.dims, li.view, li.total,
								li.scrollWidth, li.bottom, true,
							),
							).
							Fn(gtx)
					},
				).
				Fn,
			).Fn,
		)
		if li.leftSide {
			containerFlex.Rigid(EmptySpace(int(li.TextSize.V)/2, int(li.TextSize.V)/2))
			containerFlex.Rigid(li.embedWidget(li.scrollWidth + int(li.TextSize.V)/2))
		}
		container = li.Fill(li.background, l.Center, li.TextSize.V/4, 0, containerFlex.Fn).Fn
	}
	return container(gtx)
}

// EmbedWidget places the scrollable content
func (li *List) embedWidget(scrollWidth int) func(l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		if li.axis == l.Horizontal {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y - scrollWidth
			gtx.Constraints.Max.Y = gtx.Constraints.Min.Y
		} else {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X - scrollWidth
			gtx.Constraints.Max.X = gtx.Constraints.Min.X
		}
		return li.Layout(gtx, li.length, li.w)
	}
}

// pageUpDown creates the clickable areas either side of the grabber that trigger a page up/page down action
func (li *List) pageUpDown(dims DimensionList, view, total, x, y int, down bool) func(l.Context) l.Dimensions {
	button := li.pageUp
	if down {
		button = li.pageDown
	}
	return func(gtx l.Context) l.Dimensions {
		bounds := image.Rectangle{Max: gtx.Constraints.Max}
		pointer.Rect(bounds).Add(gtx.Ops)
		li.sideScroll.Add(gtx.Ops, bounds)
		return li.ButtonLayout(button.SetClick(func() {
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
		},
		).
			SetPress(func() { li.recentPageClick = time.Now() }),
		).Embed(
			li.Flex().
				Rigid(EmptySpace(x/4, y)).
				Rigid(
					li.Fill("scrim", l.Center, li.TextSize.V/4, 0, EmptySpace(x/2, y)).Fn,
				).
				Rigid(EmptySpace(x/4, y)).
				Fn,
		).Background("Transparent").CornerRadius(0).Fn(gtx)
	}
}

// grabber renders the grabber
func (li *List) grabber(dims DimensionList, x, y, viewAxis, viewCross int) func(l.Context) l.Dimensions {
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
			if de.Type == pointer.Press { // || de.Type == pointer.Drag {
			}
			if de.Type == pointer.Release {
			}
			if de.Type == pointer.Drag {
				// D.Ln("drag position", de.Position)
				if time.Now().Sub(li.recentPageClick) > time.Second/2 {
					total := dims.GetTotal(li.axis)
					var d int
					if li.axis == l.Horizontal {
						deltaX := int(de.Position.X)
						if deltaX > 8 || deltaX < -8 {
							d = deltaX * (total / viewAxis)
							li.SetPosition(dims.CoordinateToPosition(d, li.axis))
						}
					} else {
						deltaY := int(de.Position.Y)
						if deltaY > 8 || deltaY < -8 {
							d = deltaY * (total / viewAxis)
							li.SetPosition(dims.CoordinateToPosition(d, li.axis))
						}
					}
				}
				li.Window.Invalidate()
			}
		}
		defer op.Save(gtx.Ops).Load()
		bounds := image.Rectangle{Max: image.Point{X: x, Y: y}}
		pointer.Rect(bounds).Add(gtx.Ops)
		li.sideScroll.Add(gtx.Ops, bounds)
		return li.Flex().
			Rigid(
				li.Fill(li.currentColor, l.Center, 0, 0, EmptySpace(x, y)).
					Fn,
			).
			Fn(gtx)
	}
}
