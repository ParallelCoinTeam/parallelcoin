package gel

import l "github.com/p9c/gio/layout"

type Flex struct {
	flex     l.Flex
	ctx      *l.Context
	children []l.FlexChild
}

// Flex creates a new flex layout
func (th *Theme) Flex() (out *Flex) {
	return new(Flex)
}

// VFlex creates a new vertical flex layout
func (th *Theme) VFlex() (out *Flex) {
	return new(Flex).Vertical()
}

// alignment setters

// AlignStart sets alignment for layout from Start
func (f *Flex) AlignStart() (out *Flex) {
	f.flex.Alignment = l.Start
	return f
}

// AlignEnd sets alignment for layout from End
func (f *Flex) AlignEnd() (out *Flex) {
	f.flex.Alignment = l.End
	return f
}

// AlignMiddle sets alignment for layout from Middle
func (f *Flex) AlignMiddle() (out *Flex) {
	f.flex.Alignment = l.Middle
	return f
}

// AlignBaseline sets alignment for layout from Baseline
func (f *Flex) AlignBaseline() (out *Flex) {
	f.flex.Alignment = l.Baseline
	return f
}

// Axis setters

// Vertical sets axis to vertical, otherwise it is horizontal
func (f *Flex) Vertical() (out *Flex) {
	f.flex.Axis = l.Vertical
	return f
}

// Spacing setters

// SpaceStart sets the corresponding flex spacing parameter
func (f *Flex) SpaceStart() (out *Flex) {
	f.flex.Spacing = l.SpaceStart
	return f
}

// SpaceEnd sets the corresponding flex spacing parameter
func (f *Flex) SpaceEnd() (out *Flex) {
	f.flex.Spacing = l.SpaceEnd
	return f
}

// SpaceSides sets the corresponding flex spacing parameter
func (f *Flex) SpaceSides() (out *Flex) {
	f.flex.Spacing = l.SpaceSides
	return f
}

// SpaceAround sets the corresponding flex spacing parameter
func (f *Flex) SpaceAround() (out *Flex) {
	f.flex.Spacing = l.SpaceAround
	return f
}

// SpaceBetween sets the corresponding flex spacing parameter
func (f *Flex) SpaceBetween() (out *Flex) {
	f.flex.Spacing = l.SpaceBetween
	return f
}

// SpaceEvenly sets the corresponding flex spacing parameter
func (f *Flex) SpaceEvenly() (out *Flex) {
	f.flex.Spacing = l.SpaceEvenly
	return f
}

// Rigid inserts a rigid widget into the flex
func (f *Flex) Rigid(w l.Widget) (out *Flex) {
	f.children = append(f.children, l.Rigid(w))
	return f
}

// Flexed inserts a flexed widget into the flex
func (f *Flex) Flexed(wgt float32, w l.Widget) (out *Flex) {
	f.children = append(f.children, l.Flexed(wgt, w))
	return f
}

// Fn runs the ops in the context using the FlexChildren inside it
func (f *Flex) Fn(c l.Context) l.Dimensions {
	return f.flex.Layout(c, f.children...)
}
