package p9

import l "gioui.org/layout"

type _direction struct {
	l.Direction
}

// Direction creates a directional layout that sets its contents to align according to the configured direction (8
// cardinal directions and centered)
func (th *Theme) Direction() (out *_direction) {
	out = &_direction{}
	return
}

// direction setters

// NW sets the relevant direction for the Direction layout
func (d *_direction) NW() (out *_direction) {
	d.Direction = l.NW
	return d
}

// N sets the relevant direction for the Direction layout
func (d *_direction) N() (out *_direction) {
	d.Direction = l.N
	return d
}

// NE sets the relevant direction for the Direction layout
func (d *_direction) NE() (out *_direction) {
	d.Direction = l.NE
	return d
}

// E sets the relevant direction for the Direction layout
func (d *_direction) E() (out *_direction) {
	d.Direction = l.E
	return d
}

// SE sets the relevant direction for the Direction layout
func (d *_direction) SE() (out *_direction) {
	d.Direction = l.SE
	return d
}

// S sets the relevant direction for the Direction layout
func (d *_direction) S() (out *_direction) {
	d.Direction = l.S
	return d
}

// SW sets the relevant direction for the Direction layout
func (d *_direction) SW() (out *_direction) {
	d.Direction = l.SW
	return d
}

// W sets the relevant direction for the Direction layout
func (d *_direction) W() (out *_direction) {
	d.Direction = l.W
	return d
}

// Center sets the relevant direction for the Direction layout
func (d *_direction) Center() (out *_direction) {
	d.Direction = l.Center
	return d
}

// Layout the given widget given the context and direction
func (d *_direction) Fn(c *l.Context, w l.Widget) {
	d.Direction.Layout(*c, w)
}
