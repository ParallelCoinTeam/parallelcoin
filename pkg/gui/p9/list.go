package p9

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
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
	alignment l.Alignment

	ctx         l.Context
	scroll      gesture.Scroll
	scrollDelta int

	// Position is updated during Layout. To save the list scroll position, just save Position after Layout finishes. To
	// scroll the list programmatically, update Position (e.g. restore it from a saved value) before calling Layout.
	position         Position
	nextUp, nextDown Position
	len              int

	// maxSize is the total size of visible children.
	maxSize  int
	children []scrollChild
	dir      iterationDir

	length           int
	w                ListElement
	pageUp, pageDown *Clickable
}

// List returns a new scrollable List widget
func (th *Theme) List() (out *List) {
	out = &List{th: th, pageUp: th.Clickable(), pageDown: th.Clickable()}
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
	li.length = length
	return li
}

func (li *List) ListElement(w ListElement) *List {
	li.w = w
	return li
}

type DimensionList []l.Dimensions

func (d DimensionList) GetTotal(axis l.Axis) (total int) {
	for i := range d {
		if axis == l.Horizontal {
			total += d[i].Baseline*2 + d[i].Size.X
		} else {
			total += d[i].Baseline*2 + d[i].Size.Y
		}
	}
	return
}

func (d DimensionList) PositionToCoordinate(position Position, axis l.Axis) (coordinate int) {
	for i := 0; i <= position.First; i++ {
		if axis == l.Horizontal {
			coordinate += d[i].Baseline*2 + d[i].Size.X
		} else {
			coordinate += d[i].Baseline*2 + d[i].Size.Y
		}
	}
	return coordinate + position.Offset
}

func (d DimensionList) CoordinateToPosition(coordinate int, axis l.Axis) (position Position) {
	cursor := 0
	for i := range d {
		if axis == l.Horizontal {
			cursor += d[i].Baseline*2 + d[i].Size.X
		} else {
			cursor += d[i].Baseline*2 + d[i].Size.Y
		}
		if cursor > coordinate {
			if i == 0 {
				position.Offset = coordinate - cursor
				position.BeforeEnd = true
				break
			}
			// back up
			if axis == l.Horizontal {
				cursor -= d[i].Baseline*2 + d[i].Size.X
			} else {
				cursor -= d[i].Baseline*2 + d[i].Size.Y
			}
			position.First = i - 1
			position.Offset = coordinate - cursor
			position.BeforeEnd = true
			if i == len(d)-1 {
				if position.Offset == 0 {
					position.BeforeEnd = false
				}
			}
			break
		}
	}
	return
}

// Fn runs the layout in the configured context. The ListElement function returns the widget at the given index
func (li *List) Fn(gtx l.Context) l.Dimensions {
	// Debug("position", li.position, )
	// get the size of the scrollbar
	scrollWidth := int(li.th.TextSize.V)
	// render the widgets onto a second context to get their dimensions
	var ops op.Ops
	gtx1 := l.NewContext(&ops, system.FrameEvent{})
	// set constraints for same width infinite length
	if li.axis == l.Horizontal {
		gtx1.Constraints.Max.Y = gtx.Constraints.Max.Y - scrollWidth
		gtx1.Constraints.Max.X = Inf
	} else {
		gtx1.Constraints.Max.Y = Inf
		gtx1.Constraints.Max.X = gtx.Constraints.Max.X - scrollWidth
	}
	// generate the dimensions for all the list elements
	var dims DimensionList
	var total, inView int
	var before int
	if li.axis == l.Horizontal {
		inView = gtx.Constraints.Max.X
	} else {
		inView = gtx.Constraints.Max.Y
	}
	for i := 0; i < li.length; i++ {
		d := li.w(gtx1, i)
		dims = append(dims, d)
		// Debug(d)
		// depending on the axis gather the totals that create the total length
		if li.axis == l.Horizontal {
			el := d.Size.X + d.Baseline*2
			total += el
			if i < li.position.First {
				before += el
			}
		} else {
			el := d.Size.Y + d.Baseline*2
			total += el
			if i < li.position.First {
				before += el
			}
		}
	}
	before += li.position.Offset
	// compute the new positions for page up and page down

	proportion := float32(inView) / float32(total)
	if proportion > 1 {
		proportion = 1
	}
	var last, first float32
	first = float32(before) / float32(total)
	last = 1 - first - proportion
	// Debug("in view", inView, "proportion", proportion, "before", before, "total", total, "first", first, "last", last)
	// now lay it all out and draw the list and scrollbar
	var container l.Widget
	if li.axis == l.Horizontal {
		container = li.th.Flex().Flexed(1,
			li.th.Flex().Vertical().
				Rigid(
					func(l.Context) l.Dimensions {
						gtx.Constraints.Min.Y = gtx.Constraints.Max.Y - scrollWidth
						gtx.Constraints.Max.Y = gtx.Constraints.Min.Y
						return li.Layout(gtx, li.length, li.w)
					},
				).
				Rigid(
					li.th.Flex().Vertical().
						Flexed(first,
							li.th.Fill("DocBg").Embed(
								func(gtx l.Context) l.Dimensions {
									x := gtx.Constraints.Max.X
									return l.Dimensions{
										Size: image.Point{
											Y: scrollWidth,
											X: x,
										},
									}
								},
							).Fn,
						).
						Flexed(proportion,
							li.th.Fill("Primary").Embed(
								func(gtx l.Context) l.Dimensions {
									return l.Dimensions{
										Size: image.Point{
											Y: scrollWidth,
											X: gtx.Constraints.Max.X,
										},
									}
								},
							).Fn,
						).
						Flexed(last,
							li.th.Fill("DocBg").Embed(
								func(gtx l.Context) l.Dimensions {
									x := gtx.Constraints.Max.X
									return l.Dimensions{
										Size: image.Point{
											Y: scrollWidth,
											X: x,
										},
									}
								},
							).Fn,
						).
						Fn,
				).Fn,
		).Fn
	} else {
		container = li.th.Flex().Vertical().Flexed(1,
			li.th.Flex().
				Rigid(
					func(l.Context) l.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X - scrollWidth
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
						return li.Layout(gtx, li.length, li.w)
					},
				).
				Rigid(
					li.th.Flex().Vertical().
						Flexed(first,
							li.th.ButtonLayout(li.pageUp.SetClick(func() {
								current := dims.PositionToCoordinate(li.position, li.axis)
								upPos := current - inView
								if upPos < 0 {
									upPos = 0
								}
								// jump new position next page up
								li.position = dims.CoordinateToPosition(upPos, li.axis)
							})).Embed(
								li.th.Fill("DocBg").Embed(
									func(gtx l.Context) l.Dimensions {
										y := gtx.Constraints.Max.Y
										return l.Dimensions{
											Size: image.Point{
												X: scrollWidth,
												Y: y,
											},
										}
									},
								).Fn,
							).Fn,
						).
						Flexed(proportion,
							li.th.Fill("Primary").Embed(
								func(gtx l.Context) l.Dimensions {
									return l.Dimensions{
										Size: image.Point{
											X: scrollWidth,
											Y: gtx.Constraints.Max.Y,
										},
									}
								},
							).Fn,
						).
						Flexed(last,
							li.th.ButtonLayout(li.pageDown.SetClick(func() {
								current := dims.PositionToCoordinate(li.position, li.axis)
								var downPos int
								if current+inView > total {
									downPos = total - inView
								} else {
									downPos = current + inView
								}
								li.position = dims.CoordinateToPosition(downPos, li.axis)
							})).Embed(
								li.th.Fill("DocBg").Embed(
									func(gtx l.Context) l.Dimensions {
										y := gtx.Constraints.Max.Y
										return l.Dimensions{
											Size: image.Point{
												X: scrollWidth,
												Y: y,
											},
										}
									},
								).Fn,
							).Fn,
						).
						Fn,
				).Fn,
		).Fn
	}
	return container(gtx)
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
