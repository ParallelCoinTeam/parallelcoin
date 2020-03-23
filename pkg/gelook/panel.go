package gelook

import (
	"gioui.org/layout"
)

type item struct {
	i int
}

func (it *item) doSlide(n int) {
	it.i = it.i + n
}

type Panel struct {
	Size                 int
	visibleObjectsNumber int
	//totalOffset       int
	PanelContentLayout *layout.List
	PanelObject        interface{}
	PanelObjectsNumber int
	ScrollBar          *ScrollBar
	ScrollUnit         int
}

func (t *DuoUItheme) DuoUIpanel(object interface{}) *Panel {
	return &Panel{
		PanelContentLayout: &layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
		PanelObject: object,
		Size:        16,
		ScrollBar:   t.ScrollBar(),
	}
}

func (p *Panel) panelLayout(gtx *layout.Context, row func(i int, in interface{})) func() {
	return func() {

		//visibleObjectsNumber := 0
		p.PanelContentLayout.Layout(gtx, p.PanelObjectsNumber, func(i int) {
			//p.panelObject[i]
			row(i, p.PanelObject)

			//visibleObjectsNumber = visibleObjectsNumber + 1
			//p.visibleObjectsNumber = visibleObjectsNumber
			//p.totalHeight = p.totalHeight + content[i]()
		})
		//p.visibleHeight = gtx.Constraints.Height.Max
	}
}

func (p *Panel) Layout(gtx *layout.Context, row func(i int, in interface{})) {
	//p.panelObjectsNumber = len(p.panelObject)
	layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		layout.Flexed(1, p.panelLayout(gtx, row)),
		layout.Rigid(func() {
			//if p.totalOffset > 0 {
			p.SliderLayout(gtx)
			//}
		}),
	)
	//p.scrollUnit = p.scrollBar.body.Height / p.panelObjectsNumber
	cursorHeight := p.visibleObjectsNumber * p.ScrollUnit
	if cursorHeight > 30 {
		p.ScrollBar.body.CursorHeight = cursorHeight
	}

	//fmt.Println("bodyHeight:", p.scrollBar.body.Height)
	//fmt.Println("visibleObjectsNumber:", p.visibleObjectsNumber)
	//fmt.Println("scrollBarbodyPosition:", p.scrollBar.body.Position)
	//fmt.Println("scrollUnit:", p.scrollUnit)
	//fmt.Println("cursor:", p.PanelContentLayout.Position.Offset)
	//fmt.Println("First:", p.PanelContentLayout.Position.First)
	//fmt.Println("BeforeEnd:", p.PanelContentLayout.Position.BeforeEnd)
}
