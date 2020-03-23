package gel

import "gioui.org/layout"

type item struct {
	i int
}

func (it *item) doSlide(n int) {
	it.i = it.i + n
}

type Panel struct {
	Size                 int
	VisibleObjectsNumber int
	//totalOffset       int
	PanelContentLayout *layout.List
	PanelObject        interface{}
	PanelObjectsNumber int
	ScrollBar          *ScrollBar
	ScrollUnit         int
}

func (p *Panel) Layout(gtx *layout.Context) {
	if p.ScrollBar.Body.pressed {
		cs := gtx.Constraints
		if p.ScrollBar.Body.Position >= 0 && p.ScrollBar.Body.Position <= cs.Height.Max-p.ScrollBar.Body.CursorHeight {
			p.ScrollBar.Body.Cursor = p.ScrollBar.Body.Position
			p.PanelContentLayout.Position.First = p.ScrollBar.Body.Position / p.ScrollUnit
			p.PanelContentLayout.Position.Offset = 0
			//p.panelContent.Position.First = int(p.ScrollBar.body.Cursor)
		}
		//colorBg = "ffcf30cf"
		//colorBorder = "ff303030"
		//border = unit.Dp(0)
	}

}
