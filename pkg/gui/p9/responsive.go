package p9

import (
	"sort"

	l "gioui.org/layout"
)

type ContextWidget func(l.Context) l.Widget

type WidgetSize struct {
	Size   int
	Widget l.Widget
}

type Widgets []WidgetSize

func (w Widgets) Len() int {
	return len(w)
}

func (w Widgets) Less(i, j int) bool {
	// we want largest first so this uses greater than
	return w[i].Size > w[j].Size
}

func (w Widgets) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

type Responsive struct {
	Widgets
	size int
}

func (th *Theme) Responsive(size int, widgets Widgets) *Responsive {
	return &Responsive{size: size, Widgets: widgets}
}

func (r *Responsive) Embed(widgets Widgets) *Responsive {
	r.Widgets = widgets
	return r
}

func (r *Responsive) Fn(gtx l.Context) l.Dimensions {
	out := func(l.Context) l.Dimensions {
		return l.Dimensions{}
	}
	sort.Sort(r.Widgets)
	for i := range r.Widgets {
		if r.size >= r.Widgets[i].Size {
			out = r.Widgets[i].Widget
			break
		}
	}
	return out(gtx)
}
