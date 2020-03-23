package gelook

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gel"
)

type Panel struct {
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
			visibleObjectsNumber = visibleObjectsNumber + 1
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
			if panel.PanelObjectsNumber > panel.VisibleObjectsNumber {
				p.SliderLayout(gtx, panel)
			}
		}),
	)
	if panel.PanelObjectsNumber > 0 {
		panel.ScrollUnit = panel.ScrollBar.Body.Height / panel.PanelObjectsNumber
	}

	cursorHeight := panel.VisibleObjectsNumber * panel.ScrollUnit
	if cursorHeight > 30 {
		panel.ScrollBar.Body.CursorHeight = cursorHeight
	}

	//fmt.Println("cursorHeight:", cursorHeight)
	//fmt.Println("scrollUnit:", panel.ScrollUnit)
	//fmt.Println("visibleObjectsNumber:", panel.VisibleObjectsNumber)
	//fmt.Println("scrollBarbodCursorHeight:", panel.ScrollBar.Body.CursorHeight)
	//fmt.Println("scrollBarbodyPosition:", panel.ScrollBar.Body.Position)
	//fmt.Println("PanelObjectsNumber:", panel.PanelObjectsNumber)
	//fmt.Println("ScrollBar.body.Height:", p.ScrollBar.body.Height)
	//fmt.Println("cursor:", panel.PanelContentLayout.Position.Offset)
	//fmt.Println("First:", panel.PanelContentLayout.Position.First)
	//fmt.Println("BeforeEnd:", panel.PanelContentLayout.Position.BeforeEnd)
	panel.Layout(gtx)
}
