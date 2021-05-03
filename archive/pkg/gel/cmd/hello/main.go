package main

import (
	l "github.com/p9c/gio/layout"
	"github.com/p9c/qu"

	"github.com/p9c/gel"
)

type State struct {
	*gel.Window
}

func NewState(quit qu.C) *State {
	return &State{
		Window: gel.NewWindowP9(quit),
	}
}

func main() {
	quit := qu.T()
	state := NewState(quit)
	var e error
	I.Ln("logging")
	rootWidget := state.rootWidget()
	if e = state.Window.
		Size(48, 32).
		Title("hello world").
		Open().
		Run(rootWidget, quit.Q, quit); E.Chk(e) {
	}
}

func (s *State) rootWidget() l.Widget {
	return s.Direction().Center().Embed(s.H2("hello world!").Fn).Fn
}
