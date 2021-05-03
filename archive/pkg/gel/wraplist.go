package gel

import (
	l "github.com/p9c/gio/layout"
	"golang.org/x/exp/shiny/text"
)

// WrapList is a generalised layout for creating lists from widgets lined up in one axis and wrapped into lines across
// a given axis. It can be used for an icon view, for a text console cell grid, for laying out selectable text.
type WrapList struct {
	*Window
	list    *List
	widgets []l.Widget
	axis    l.Axis
	dir     text.Direction
}

// WrapList creates a new WrapList
func (w *Window) WrapList() *WrapList {
	return &WrapList{
		Window: w,
		list:   w.WidgetPool.GetList(),
	}
}

// Axis sets the axis that will be scrollable
func (w *WrapList) Axis(axis l.Axis) *WrapList {
	w.axis = axis
	return w
}

// Direction sets the direction across the axis, for vertical, text.Forwards means left to right, text.Backwards means
// right to left, and for horizontal, text.Forwards means top to bottom, text.Backwards means bottom to top
func (w *WrapList) Direction(dir text.Direction) *WrapList {
	w.dir = dir
	return w
}

// Widgets loads a set of widgets into the WrapList
func (w *WrapList) Widgets(widgets []l.Widget) *WrapList {
	w.widgets = widgets
	return w
}

// Fn renders the WrapList in the current context
// todo: this needs to be cached and have a hook in the WrapList.Widgets method to trigger generation
func (w *WrapList) Fn(gtx l.Context) l.Dimensions {
	// first get the dimensions of the widget list
	if len(w.widgets) > 0 {
		le := func(gtx l.Context, index int) l.Dimensions {
			dims := w.widgets[index](gtx)
			return dims
		}
		gtx1 := CopyContextDimensionsWithMaxAxis(gtx, w.axis)
		// generate the dimensions for all the list elements
		allDims := GetDimensionList(gtx1, len(w.widgets), le)
		var out = [][]l.Widget{{}}
		cursor := 0
		runningTotal := 0
		_, width := axisCrossConstraint(w.axis, gtx.Constraints)
		for i := range allDims {
			runningTotal += axisCross(w.axis, allDims[i].Size)
			if runningTotal > width {
				cursor++
				runningTotal = 0
				out = append(out, []l.Widget{})
			}
			out[cursor] = append(out[cursor], w.widgets[i])
		}
		le2 := func(gtx l.Context, index int) l.Dimensions {
			o := w.Flex().AlignStart()
			if w.axis == l.Horizontal {
				o = w.VFlex().AlignStart()
			}
			// for _, y := range out {
			y := out[index]
			var outRow []l.Widget
			for j, x := range y {
				if w.dir == text.Forwards {
					outRow = append(outRow, x)
				} else {
					outRow = append(outRow, y[len(y)-1-j])
				}
			}
			for i := range outRow {
				o.Rigid(outRow[i])
			}
			// }
			return o.Fn(gtx)
		}
		return w.list.Vertical().Length(len(out)).ListElement(le2).Fn(gtx)
	} else {
		return l.Dimensions{}
	}
}
