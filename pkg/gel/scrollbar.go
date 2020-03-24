package gel

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
)

type item struct {
	i int
}

func (it *item) doSlide(n int) {
	it.i = it.i + n
}

type ScrollBar struct {
	Size   int
	Slider *Slider
	Up     *Button
	Down   *Button
}

type Slider struct {
	Do           func(interface{})
	OperateValue interface{}
	pressed      bool
	Position     int
	Cursor       int
	Height       int
	CursorHeight int
	//Icon         DuoUIicon
}

type ScrollBarButton struct {
	//button      DuoUIbutton
	Height      int
	insetTop    float32
	insetRight  float32
	insetBottom float32
	insetLeft   float32
	iconSize    int
	iconPadding float32
}

func (s *Slider) Layout(gtx *layout.Context) {
	//fmt.Println("He::", gtx.Constraints.Height.Max)
	//fmt.Println("wi::", gtx.Constraints.Width.Max)

	for _, e := range gtx.Events(s) {
		if e, ok := e.(pointer.Event); ok {
			//s.Body.Position = e.Position.Y - float32(s.CursorHeight/2)
			if e.Position.Y > 0 {
				s.Position = int(e.Position.Y) - (s.CursorHeight / 2)
			}
			switch e.Type {
			case pointer.Press:
				s.pressed = true
				s.Do(s.OperateValue)
				// list.Position.First = int(s.Position)
				L.Debug("RADI PRESS")
			case pointer.Release:
				s.pressed = false
			}
		}
	}
	cs := gtx.Constraints
	s.Height = cs.Height.Max
}
