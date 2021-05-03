package gel

import (
	"sort"
	
	l "github.com/p9c/gio/layout"
)

// WidgetSize is a widget with a specification of the minimum size to select it for viewing.
// Note that the widgets you put in here should be wrapped in func(l.Context) l.Dimensions otherwise
// any parameters retrieved from the controlling state variable will be from initialization and not
// at execution of the widget in the render process
type WidgetSize struct {
	Size   float32
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
	*Theme
	Widgets
	size int32
}

func (th *Theme) Responsive(size int32, widgets Widgets) *Responsive {
	return &Responsive{Theme: th, size: size, Widgets: widgets}
}

func (r *Responsive) Fn(gtx l.Context) l.Dimensions {
	out := func(l.Context) l.Dimensions {
		return l.Dimensions{}
	}
	sort.Sort(r.Widgets)
	for i := range r.Widgets {
		if float32(r.size)/r.TextSize.V >= r.Widgets[i].Size {
			out = r.Widgets[i].Widget
			// D.Ln("selected widget for responsive with scale", r.size, "width", r.Widgets[i].Size)
			break
		}
	}
	return out(gtx)
}
