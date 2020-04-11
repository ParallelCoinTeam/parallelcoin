package fui

import "gioui.org/layout"

type direction struct {
	layout.Direction
}

// Direction creates a directional layout that sets its contents to align
// according to the configured direction (8 cardinal directions and centered)
func Direction() (out *direction) {
	out = &direction{}
	return
}

// direction setters

// NW sets the relevant direction for the Direction layout
func (d *direction) NW() (out *direction) {
	d.Direction = layout.NW
	return d
}

// N sets the relevant direction for the Direction layout
func (d *direction) N() (out *direction) {
	d.Direction = layout.N
	return d
}

// NE sets the relevant direction for the Direction layout
func (d *direction) NE() (out *direction) {
	d.Direction = layout.NE
	return d
}

// E sets the relevant direction for the Direction layout
func (d *direction) E() (out *direction) {
	d.Direction = layout.E
	return d
}

// SE sets the relevant direction for the Direction layout
func (d *direction) SE() (out *direction) {
	d.Direction = layout.SE
	return d
}

// S sets the relevant direction for the Direction layout
func (d *direction) S() (out *direction) {
	d.Direction = layout.S
	return d
}

// SW sets the relevant direction for the Direction layout
func (d *direction) SW() (out *direction) {
	d.Direction = layout.SW
	return d
}

// W sets the relevant direction for the Direction layout
func (d *direction) W() (out *direction) {
	d.Direction = layout.W
	return d
}

// Center sets the relevant direction for the Direction layout
func (d *direction) Center() (out *direction) {
	d.Direction = layout.Center
	return d
}

// Layout the given widget given the context and direction
func (d *direction) Layout(c *layout.Context, w layout.Widget) {
	d.Direction.Layout(c, w)
}

// Prep the given widget given the context and direction
func (d *direction) Prep(c *layout.Context, w layout.Widget) func() {
	return func() { d.Direction.Layout(c, w) }
}
