package gel

import l "github.com/p9c/gio/layout"

type Stack struct {
	*l.Stack
	children []l.StackChild
}

// Stack starts a chain of widgets to compose into a stack
func (w *Window) Stack() (out *Stack) {
	out = &Stack{
		Stack: &l.Stack{},
	}
	return
}

func (s *Stack) Alignment(alignment l.Direction) *Stack {
	s.Stack.Alignment = alignment
	return s
}

// functions to chain widgets to stack (first is lowest last highest)

// Stacked appends a widget to the stack, the stack's dimensions will be
// computed from the largest widget in the stack
func (s *Stack) Stacked(w l.Widget) (out *Stack) {
	s.children = append(s.children, l.Stacked(w))
	return s
}

// Expanded lays out a widget with the same max constraints as the stack
func (s *Stack) Expanded(w l.Widget) (out *Stack) {
	s.children = append(s.children, l.Expanded(w))
	return s
}

// Fn runs the ops queue configured in the stack
func (s *Stack) Fn(c l.Context) l.Dimensions {
	return s.Stack.Layout(c, s.children...)
}
