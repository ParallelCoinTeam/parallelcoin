package gelook

import (
	"gioui.org/layout"
	"github.com/p9c/gel"
)

var (
	c *gel.ScrollBar
)

type item struct {
	i int
}

func (it *item) doSlide(n int) {
	it.i = it.i + n
}

type DuoUIpanel struct {
	Name        string
	panelObject []func()
	scrollBar   *ScrollBar
}

func (t *DuoUItheme) DuoUIpanel(content func()) *DuoUIpanel {
	return &DuoUIpanel{
		Name: "OneDuoUIpanel",
		panelObject: []func(){
			content,
		},
		scrollBar: t.ScrollBar(c),
	}
}

func (p *DuoUIpanel) Layout(gtx *layout.Context, panel *gel.Panel) {
	layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		layout.Flexed(1, func() {
			panel.PanelContentLayout.Layout(gtx, len(p.panelObject), func(i int) {
				p.panelObject[i]()
				panel.TotalHeight = gtx.Dimensions.Size.Y
			})
			panel.VisibleHeight = gtx.Constraints.Height.Max
		}),
		layout.Rigid(func() {
			// if panel.TotalOffset > 0 {
			// p.scrollBar = t.ScrollBar(32)
			p.scrollBar.Layout(gtx,
				panel.PanelContentLayout.Position.Offset,
				panel.ScrollUnit,
			)
			// }
		}),
	)
	panel.Layout(gtx)
}
