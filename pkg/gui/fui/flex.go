package fui

import "gioui.org/layout"

type flex struct {
	flex     layout.Flex
	ctx      *layout.Context
	children []layout.FlexChild
}

// Flex creates a new flex layout
func Flex() (out *flex) {
	out = &flex{}
	return
}

// Alignment setters

// AlignStart sets alignment for layout from Start
func (f *flex) AlignStart() (out *flex) {
	f.flex.Alignment = layout.Start
	return f
}

// AlignEnd sets alignment for layout from End
func (f *flex) AlignEnd() (out *flex) {
	f.flex.Alignment = layout.End
	return f
}

// AlignMiddle sets alignment for layout from Middle
func (f *flex) AlignMiddle() (out *flex) {
	f.flex.Alignment = layout.Middle
	return f
}

// AlignBaseline sets alignment for layout from Baseline
func (f *flex) AlignBaseline() (out *flex) {
	f.flex.Alignment = layout.Baseline
	return f
}

// Axis setters

// Vertical sets axis to vertical, otherwise it is horizontal
func (f *flex) Vertical() (out *flex) {
	f.flex.Axis = layout.Vertical
	return f
}

// Spacing setters

// SpaceStart sets the corresponding flex spacing parameter
func (f *flex) SpaceStart() (out *flex) {
	f.flex.Spacing = layout.SpaceStart
	return f
}

// SpaceEnd sets the corresponding flex spacing parameter
func (f *flex) SpaceEnd() (out *flex) {
	f.flex.Spacing = layout.SpaceEnd
	return f
}

// SpaceSides sets the corresponding flex spacing parameter
func (f *flex) SpaceSides() (out *flex) {
	f.flex.Spacing = layout.SpaceSides
	return f
}

// SpaceAround sets the corresponding flex spacing parameter
func (f *flex) SpaceAround() (out *flex) {
	f.flex.Spacing = layout.SpaceAround
	return f
}

// SpaceBetween sets the corresponding flex spacing parameter
func (f *flex) SpaceBetween() (out *flex) {
	f.flex.Spacing = layout.SpaceBetween
	return f
}

// SpaceEvenly sets the corresponding flex spacing parameter
func (f *flex) SpaceEvenly() (out *flex) {
	f.flex.Spacing = layout.SpaceEvenly
	return f
}

// Rigid inserts a rigid widget into the flex
func (f *flex) Rigid(w layout.Widget) (out *flex) {
	f.children = append(f.children, layout.Rigid(w))
	return f
}

// Flexed inserts a flexed widget into the flex
func (f *flex) Flexed(wgt float32, w layout.Widget) (out *flex) {
	f.children = append(f.children, layout.Flexed(wgt, w))
	return f
}

// Layout runs the ops in the context using the FlexChildren inside it
func (f *flex) Layout(c *layout.Context) {
	f.flex.Layout(c, f.children...)
}

// Prep a layout but return it instead of running it (return a layout.Widget)
func (f *flex) Prep(c *layout.Context) func() {
	return func() { f.flex.Layout(c, f.children...) }
}
