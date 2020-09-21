package component

import (
	"github.com/stalker-loki/pod/pkg/gui/gelook"
)

type DuoUIcomponent struct {
	Name    string
	Version string
	Theme   *gelook.DuoUItheme
	M       interface{}
	V       func()
	C       func()
}
