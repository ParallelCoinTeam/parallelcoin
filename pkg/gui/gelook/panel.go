package gelook

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/gel"
)

type DuoUIPanel struct {
	PanelObject interface{}
	ScrollBar   *ScrollBar
	container   DuoUIContainer
}

func (t *DuoUITheme) DuoUIPanel() DuoUIPanel {
	return DuoUIPanel{
		container: t.DuoUIContainer(0, t.Colors["Light"]),
	}
}
func (p *DuoUIPanel) panelLayout(gtx *layout.Context, panel *gel.Panel, row func(i int, in interface{})) func() {
	return func() {
		visibleObjectsNumber := 0
		panel.PanelContentLayout.Layout(gtx, panel.PanelObjectsNumber, func(i int) {
			row(i, p.PanelObject)
			visibleObjectsNumber = visibleObjectsNumber + 1
			panel.VisibleObjectsNumber = visibleObjectsNumber
		})
	}
}

func (p *DuoUIPanel) Layout(gtx *layout.Context, panel *gel.Panel, row func(i int, in interface{})) {
	p.container.Layout(gtx, layout.NW, func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Flexed(1, p.panelLayout(gtx, panel, row)),
			layout.Rigid(func() {
				if panel.PanelObjectsNumber > panel.VisibleObjectsNumber {
					p.ScrollBarLayout(gtx, panel)
				}
			}),
		)
		//fmt.Println("scrollUnit:", panel.ScrollUnit)
		//fmt.Println("ScrollBar.Slider.Height:", panel.ScrollBar.Slider.Height)
		//fmt.Println("PanelObjectsNumber:", panel.PanelObjectsNumber)

		panel.Layout(gtx)
	})
}
