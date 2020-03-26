package gel

import (
	"gioui.org/layout"
)

type Panel struct {
	Size                 int
	VisibleObjectsNumber int
	//totalOffset       int
	PanelContentLayout *layout.List
	PanelObject        interface{}
	PanelObjectsNumber int
	ScrollBar          *ScrollBar
	ScrollUnit         float32
}

func NewPanel() *Panel {
	itemValue := item{
		i: 0,
	}
	return &Panel{
		PanelContentLayout: &layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
		ScrollBar: &ScrollBar{
			Size: 16,
			Slider: &Slider{
				Do: func(n interface{}) {
					itemValue.doSlide(n.(int))
				},
				OperateValue: 1,
				pressed:      false,
			},
			Up:   new(Button),
			Down: new(Button),
		},
	}
}

func (p *Panel) Layout(gtx *layout.Context) {
	if p.PanelObjectsNumber > 0 {
		p.ScrollUnit = float32(p.ScrollBar.Slider.Height) / float32(p.PanelObjectsNumber)
	}
	cursorHeight := int(float32(p.VisibleObjectsNumber) * p.ScrollUnit)
	if cursorHeight > p.Size {
		p.ScrollBar.Slider.CursorHeight = cursorHeight
	}
	if p.ScrollBar.Slider.pressed {
		cs := gtx.Constraints
		if p.ScrollBar.Slider.Position >= 0 && p.ScrollBar.Slider.Position <= cs.Height.Max-p.ScrollBar.Slider.CursorHeight {
			p.ScrollBar.Slider.Cursor = p.ScrollBar.Slider.Position
			p.PanelContentLayout.Position.First = int(float32(p.ScrollBar.Slider.Position) / p.ScrollUnit)
			p.PanelContentLayout.Position.Offset = 0
		}

	}
}
