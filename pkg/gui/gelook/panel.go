package gelook

import (
	"gioui.org/layout"
	"github.com/stalker-loki/pod/pkg/gui/gel"
)

type DuoUIpanel struct {
	PanelObject interface{}
	ScrollBar   *ScrollBar
	container   DuoUIcontainer
}

func (t *DuoUItheme) DuoUIpanel() DuoUIpanel {
	return DuoUIpanel{
		container: t.DuoUIcontainer(0, t.Colors["Light"]),
	}
}
func (p *DuoUIpanel) panelLayout(gtx *layout.Context, panel *gel.Panel, row func(i int, in interface{})) func() {
	return func() {
		visibleObjectsNumber := 0
		panel.PanelContentLayout.Layout(gtx, panel.PanelObjectsNumber, func(i int) {
			row(i, p.PanelObject)
			visibleObjectsNumber = visibleObjectsNumber + 1
			panel.VisibleObjectsNumber = visibleObjectsNumber
		})
	}
}

func (p *DuoUIpanel) Layout(gtx *layout.Context, panel *gel.Panel, row func(i int, in interface{})) {
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
