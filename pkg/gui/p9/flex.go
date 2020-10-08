package p9

import l "gioui.org/layout"

type _flex struct {
	flex     l.Flex
	ctx      *l.Context
	children []l.FlexChild
}

// Flex creates a new flex layout
func (th *Theme) Flex() (out *_flex) {
	out = &_flex{}
	return
}

// Alignment setters

// AlignStart sets alignment for layout from Start
func (f *_flex) AlignStart() (out *_flex) {
	f.flex.Alignment = l.Start
	return f
}

// AlignEnd sets alignment for layout from End
func (f *_flex) AlignEnd() (out *_flex) {
	f.flex.Alignment = l.End
	return f
}

// AlignMiddle sets alignment for layout from Middle
func (f *_flex) AlignMiddle() (out *_flex) {
	f.flex.Alignment = l.Middle
	return f
}

// AlignBaseline sets alignment for layout from Baseline
func (f *_flex) AlignBaseline() (out *_flex) {
	f.flex.Alignment = l.Baseline
	return f
}

// Axis setters

// Vertical sets axis to vertical, otherwise it is horizontal
func (f *_flex) Vertical() (out *_flex) {
	f.flex.Axis = l.Vertical
	return f
}

// Spacing setters

// SpaceStart sets the corresponding flex spacing parameter
func (f *_flex) SpaceStart() (out *_flex) {
	f.flex.Spacing = l.SpaceStart
	return f
}

// SpaceEnd sets the corresponding flex spacing parameter
func (f *_flex) SpaceEnd() (out *_flex) {
	f.flex.Spacing = l.SpaceEnd
	return f
}

// SpaceSides sets the corresponding flex spacing parameter
func (f *_flex) SpaceSides() (out *_flex) {
	f.flex.Spacing = l.SpaceSides
	return f
}

// SpaceAround sets the corresponding flex spacing parameter
func (f *_flex) SpaceAround() (out *_flex) {
	f.flex.Spacing = l.SpaceAround
	return f
}

// SpaceBetween sets the corresponding flex spacing parameter
func (f *_flex) SpaceBetween() (out *_flex) {
	f.flex.Spacing = l.SpaceBetween
	return f
}

// SpaceEvenly sets the corresponding flex spacing parameter
func (f *_flex) SpaceEvenly() (out *_flex) {
	f.flex.Spacing = l.SpaceEvenly
	return f
}

// Rigid inserts a rigid widget into the flex
func (f *_flex) Rigid(w l.Widget) (out *_flex) {
	f.children = append(f.children, l.Rigid(w))
	return f
}

// Flexed inserts a flexed widget into the flex
func (f *_flex) Flexed(wgt float32, w l.Widget) (out *_flex) {
	f.children = append(f.children, l.Flexed(wgt, w))
	return f
}

// Fn runs the ops in the context using the FlexChildren inside it
func (f *_flex) Fn(c l.Context) l.Dimensions {
	return f.flex.Layout(c, f.children...)
}
