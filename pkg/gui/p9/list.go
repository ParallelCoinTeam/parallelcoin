package p9

import l "gioui.org/layout"

type List struct {
	*l.List
	length int
	w      l.ListElement
}

// List returns a new scrollable List widget
func (th *Theme) List() (out *List) {
	out = &List{
		List: &l.List{},
	}
	return
}

// Vertical sets the axis to vertical (default implicit is horizontal)
func (li *List) Vertical() (out *List) {
	li.List.Axis = l.Vertical
	return li
}

// ScrollToEnd sets the List to add new items to the end and push older ones up/left and initial render has scroll to
// the end (or bottom) of the List
func (li *List) ScrollToEnd() (out *List) {
	li.List.ScrollToEnd = true
	return li
}

func (li *List) Length(length int) *List {
	li.length = length
	return li
}

func (li *List) ListElement(w l.ListElement) *List {
	li.w = w
	return li
}

// Fn runs the layout in the configured context. The ListElement function returns the widget at the given index
func (li *List) Fn(gtx l.Context) l.Dimensions {
	return li.List.Layout(gtx, li.length, li.w)
}
