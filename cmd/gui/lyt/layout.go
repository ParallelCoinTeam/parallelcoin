package lyt

import (
	"gioui.org/layout"
)

func (l *DuoUIlayouts) DuoUIlayouts() {
	l.View = &layout.Flex{Axis: layout.Vertical}
	l.Main = &layout.Flex{Axis: layout.Horizontal}
	l.Menu = &layout.Flex{Axis: layout.Vertical}
	l.Status = &layout.Flex{Axis: layout.Vertical}
	return
}
