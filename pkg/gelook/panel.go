package gelook

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gel"
)

type Panel struct {
	Size int
	//totalOffset       int
	//PanelContentLayout *layout.List
	PanelObject interface{}
	//PanelObjectsNumber int
	ScrollBar *ScrollBar
}

//func (t *DuoUItheme) DuoUIpanel(object interface{}) *Panel {
//	return &Panel{
//		PanelContentLayout: &layout.List{
//			Axis:        layout.Vertical,
//			ScrollToEnd: false,
//		},
//		PanelObject: object,
//		Size:        16,
//		ScrollBar:   t.ScrollBar(),
//	}
//}

func (p *Panel) panelLayout(gtx *layout.Context, panel *gel.Panel, row func(i int, in interface{})) func() {
	return func() {
		visibleObjectsNumber := 0
		panel.PanelContentLayout.Layout(gtx, panel.PanelObjectsNumber, func(i int) {
			row(i, p.PanelObject)
			panel.VisibleObjectsNumber = visibleObjectsNumber
		})
	}
}

func (p *Panel) Layout(gtx *layout.Context, panel *gel.Panel, row func(i int, in interface{})) {
	//p.PanelObjectsNumber = len(p.PanelObject)
	layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		layout.Flexed(1, p.panelLayout(gtx, panel, row)),
		layout.Rigid(func() {
			//if p.totalOffset > 0 {
			p.SliderLayout(gtx, panel)
			//}
		}),
	)
	panel.ScrollUnit = p.ScrollBar.body.Height / panel.PanelObjectsNumber
	cursorHeight := panel.VisibleObjectsNumber * panel.ScrollUnit
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
	panel.Layout(gtx)
}
