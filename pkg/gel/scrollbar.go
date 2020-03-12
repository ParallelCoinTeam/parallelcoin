package gel

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	log "github.com/p9c/logi"
)

type item struct {
	i int
}

func (it *item) doSlide(n int) {
	it.i = it.i + n
}

type ScrollBar struct {
	//Height       float32
	body *ScrollBarBody
	//up   *ScrollBarButton
	//down *ScrollBarButton
	BodyHeight   int
	CursorHeight int
	Cursor       float32
	Position     float32
}
type ScrollBarBody struct {
	pressed      bool
	Do           func(interface{})
	OperateValue interface{}
}

func (s *ScrollBar) Layout(gtx *layout.Context) {
	s.BodyHeight = gtx.Constraints.Height.Max

	//// Flush clicks from before the previous frame.
	//b.clicks -= b.prevClicks
	//b.prevClicks = 0
	s.processEvents(gtx)
	//b.click.Add(gtx.Ops)
	//for len(b.history) > 0 {
	//	c := b.history[0]
	//	if gtx.Now().Sub(c.Time) < 1*time.Second {
	//		break
	//	}
	//	copy(b.history, b.history[1:])
	//	b.history = b.history[:len(b.history)-1]
	//}
}

func (s *ScrollBar) processEvents(gtx *layout.Context) {
	for _, e := range gtx.Events(s.body) {
		if e, ok := e.(pointer.Event); ok {
			s.Position = e.Position.Y - float32(s.CursorHeight/2)
			switch e.Type {
			case pointer.Press:
				s.body.pressed = true
				s.body.Do(s.body.OperateValue)
				//list.Position.First = int(s.Position)
				log.L.Debug("RADI PRESS")
			case pointer.Release:
				s.body.pressed = false
			}
		}
	}
}
