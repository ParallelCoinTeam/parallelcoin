package component

import (
	"github.com/p9c/pod/pkg/gui/theme"
)

type DuoUIcomponent struct {
	Name    string
	Version string
	Theme   *theme.DuoUItheme
	M       interface{}
	V       func()
	C       func()
}
