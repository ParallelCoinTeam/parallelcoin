package gel

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
)

type ScrollBar struct {
	Size         int
	OperateValue interface{}
	//Height       float32
	Body *ScrollBarBody
	Up   *Button
	Down *Button
}

type ScrollBarBody struct {
	pressed      bool
	Do           func(interface{})
	ColorBg      string
	Position     int
	Cursor       int
	OperateValue interface{}
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

func (s *ScrollBar) Layout(gtx *layout.Context) {
	// s.BodyHeight = gtx.Constraints.Height.Max

	// // Flush clicks from before the previous frame.
	// b.clicks -= b.prevClicks
	// b.prevClicks = 0
	s.processEvents(gtx)
	// b.click.Add(gtx.Ops)
	// for len(b.history) > 0 {
	//	c := b.history[0]
	//	if gtx.Now().Sub(c.Time) < 1*time.Second {
	//		break
	//	}
	//	copy(b.history, b.history[1:])
	//	b.history = b.history[:len(b.history)-1]
	// }
}

func (s *ScrollBar) processEvents(gtx *layout.Context) {
	for _, e := range gtx.Events(s.Body) {
		if e, ok := e.(pointer.Event); ok {
			//s.Body.Position = e.Position.Y - float32(s.CursorHeight/2)
			switch e.Type {
			case pointer.Press:
				s.Body.pressed = true
				s.Body.Do(s.Body.OperateValue)
				// list.Position.First = int(s.Position)
				L.Debug("RADI PRESS")
			case pointer.Release:
				s.Body.pressed = false
			}
		}
	}
}
