package p9

import l "gioui.org/layout"

type _stack struct {
	*l.Stack
	children []l.StackChild
}

// Stack starts a chain of widgets to compose into a stack
func (th *Theme) Stack() (out *_stack) {
	out = &_stack{
		Stack: &l.Stack{},
	}
	return
}

func (s *_stack) Alignment(alignment l.Direction) *_stack {
	s.Stack.Alignment = alignment
	return s
}

// functions to chain widgets to stack (first is lowest last highest)

// Stacked appends a widget to the stack, the stack's dimensions will be
// computed from the largest widget in the stack
func (s *_stack) Stacked(w l.Widget) (out *_stack) {
	s.children = append(s.children, l.Stacked(w))
	return s
}

// Expanded lays out a widget with the same max constraints as the stack
func (s *_stack) Expanded(w l.Widget) (out *_stack) {
	s.children = append(s.children, l.Expanded(w))
	return s
}

// Fn runs the ops queue configured in the stack
func (s *_stack) Fn(c l.Context) l.Dimensions {
	return s.Stack.Layout(c, s.children...)
}
