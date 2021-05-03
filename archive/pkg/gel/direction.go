package gel

import l "github.com/p9c/gio/layout"

type Direction struct {
	l.Direction
	w l.Widget
}

// Direction creates a directional layout that sets its contents to align according to the configured direction (8
// cardinal directions and centered)
func (w *Window) Direction() (out *Direction) {
	out = &Direction{}
	return
}

// direction setters

// NW sets the relevant direction for the Direction layout
func (d *Direction) NW() (out *Direction) {
	d.Direction = l.NW
	return d
}

// N sets the relevant direction for the Direction layout
func (d *Direction) N() (out *Direction) {
	d.Direction = l.N
	return d
}

// NE sets the relevant direction for the Direction layout
func (d *Direction) NE() (out *Direction) {
	d.Direction = l.NE
	return d
}

// E sets the relevant direction for the Direction layout
func (d *Direction) E() (out *Direction) {
	d.Direction = l.E
	return d
}

// SE sets the relevant direction for the Direction layout
func (d *Direction) SE() (out *Direction) {
	d.Direction = l.SE
	return d
}

// S sets the relevant direction for the Direction layout
func (d *Direction) S() (out *Direction) {
	d.Direction = l.S
	return d
}

// SW sets the relevant direction for the Direction layout
func (d *Direction) SW() (out *Direction) {
	d.Direction = l.SW
	return d
}

// W sets the relevant direction for the Direction layout
func (d *Direction) W() (out *Direction) {
	d.Direction = l.W
	return d
}

// Center sets the relevant direction for the Direction layout
func (d *Direction) Center() (out *Direction) {
	d.Direction = l.Center
	return d
}

func (d *Direction) Embed(w l.Widget) *Direction {
	d.w = w
	return d
}

// Fn the given widget given the context and direction
func (d *Direction) Fn(c l.Context) l.Dimensions {
	return d.Direction.Layout(c, d.w)
}
