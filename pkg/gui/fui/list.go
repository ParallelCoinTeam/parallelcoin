package fui

import "gioui.org/layout"

type list struct {
	*layout.List
}

// List returns a new scrollable list widget
func List() (out *list) {
	out = &list{
		List: &layout.List{},
	}
	return
}

// Vertical sets the axis to vertical (default implicit is horizontal)
func (l *list) Vertical() (out *list) {
	l.List.Axis = layout.Vertical
	return l
}

// ScrollToEnd sets the list to add new items to the end and push older ones
// up/left and initial render has scroll to the end (or bottom) of the list
func (l *list) ScrollToEnd() (out *list) {
	l.List.ScrollToEnd = true
	return l
}

// Layout runs the layout in the configured context. The ListElement function
// returns the widget at the given index
func (l *list) Layout(c *layout.Context, length int, w layout.ListElement) {
	l.List.Layout(c, length, w)
}

// Prep the layout in the configured context. The ListElement function
// returns the widget at the given index.
func (l *list) Prep(c *layout.Context, length int, w layout.ListElement) func() {
	return func() { l.List.Layout(c, length, w) }
}
