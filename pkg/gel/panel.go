package gel

import (
	"fmt"
	"gioui.org/layout"
)

var (
	itemValue = item{
		i: 0,
	}
	body = &ScrollBarBody{
		pressed: false,
		Do: func(n interface{}) {
			itemValue.doSlide(n.(int))
		},
		OperateValue: 1,
	}
)

type Panel struct {
	Name          string
	TotalHeight   int
	VisibleHeight int
	TotalOffset   int
	//panelObject       []func()
	PanelContentLayout *layout.List
	PanelObjectHeight  int
	ScrollUnit         float32
	ScrollBar          *ScrollBar
}

func (p *Panel) Layout(gtx *layout.Context) {
	p.ScrollBar = &ScrollBar{
		body: body,
	}
	p.TotalOffset = p.TotalHeight - p.VisibleHeight

	p.ScrollUnit = float32(p.ScrollBar.BodyHeight) / float32(p.TotalHeight)
	p.ScrollBar.CursorHeight = int(p.ScrollUnit * float32(p.VisibleHeight))
	p.ScrollBar.Cursor = float32(p.PanelContentLayout.Position.Offset * int(p.ScrollUnit))
	fmt.Println("bodyHeight:", p.ScrollBar.BodyHeight)
	fmt.Println("cursorHeight:", p.ScrollBar.CursorHeight)

	fmt.Println("totalOffset:", p.TotalOffset)
	fmt.Println("scrollUnit:", p.ScrollUnit)

	//fmt.Println("cursor:", p.scrollBar.body.Cursor)
	fmt.Println("visibleHeight:", p.VisibleHeight)

	fmt.Println("total:", p.TotalHeight)
	//fmt.Println("offset:", p.panelContent.Position.Offset)
}
