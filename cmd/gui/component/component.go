package component

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/theme"
)

type DuoUIcomponent struct {
	Name    string
	Version string
	Context *layout.Context
	Theme   *theme.DuoUItheme
	M       interface{}
	V       func()
	C       func()
}
