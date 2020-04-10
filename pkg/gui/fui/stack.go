package fui

import "gioui.org/layout"

type stack struct {
	*layout.Stack
	children []layout.StackChild
}

// Stack starts a chain of widgets to compose into a stack
func Stack() (out *stack) {
	out = &stack{
		Stack: &layout.Stack{},
	}
	return
}

// functions to chain widgets to stack (first is lowest last highest)

// Stacked appends a widget to the stack, the stack's dimensions will be
// computed from the largest widget in the stack
func (s *stack) Stacked(w layout.Widget) (out *stack) {
	s.children = append(s.children, layout.Stacked(w))
	return s
}

// Expanded lays out a widget with the same max constraints as the stack
func (s *stack) Expanded(w layout.Widget) (out *stack) {
	s.children = append(s.children, layout.Expanded(w))
	return s
}

// Layout runs the ops queue configured in the stack
func (s *stack) Layout(c *layout.Context) {
	s.Stack.Layout(c, s.children...)
}

// Prep the ops queue configured in the stack
func (s *stack) Prep(c *layout.Context) func() {
	return func() { s.Layout(c) }
}
