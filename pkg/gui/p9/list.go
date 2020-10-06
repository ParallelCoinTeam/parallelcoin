package p9

import lo "gioui.org/layout"

type _list struct {
	*lo.List
}

// List returns a new scrollable _list widget
func (th *Theme) List() (out *_list) {
	out = &_list{
		List: &lo.List{},
	}
	return
}

// Vertical sets the axis to vertical (default implicit is horizontal)
func (l *_list) Vertical() (out *_list) {
	l.List.Axis = lo.Vertical
	return l
}

// ScrollToEnd sets the _list to add new items to the end and push older ones
// up/left and initial render has scroll to the end (or bottom) of the _list
func (l *_list) ScrollToEnd() (out *_list) {
	l.List.ScrollToEnd = true
	return l
}

// Fn runs the layout in the configured context. The ListElement function
// returns the widget at the given index
func (l *_list) Fn(c *lo.Context, length int, w lo.ListElement) {
	l.List.Layout(*c, length, w)
}
