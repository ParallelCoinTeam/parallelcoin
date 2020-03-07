package theme

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/controller"
)

type item struct {
	i int
}

func (it *item) doSlide(n int) {
	it.i = it.i + n
}

type DuoUIpanel struct {
	Name               string
	totalHeight        int
	visibleHeight      int
	totalOffset        int
	panelContent       *func()
	panelContentLayout *layout.List
	panelObject        []func()
	panelObjectHeight  int
	scrollBar          *ScrollBar
	scrollUnit         float32
}

func (t *DuoUItheme) DuoUIpanel(content *func()) *DuoUIpanel {
	return &DuoUIpanel{
		Name:         "OneDuoUIpanel",
		panelContent: content,
		panelContentLayout: &layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
		scrollBar: t.ScrollBar(),
	}
}

func (p *DuoUIpanel) Layout(gtx *layout.Context, panel controller.Panel) {
	layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		layout.Flexed(1, panel.Panel(gtx, *p.panelContent)),
		layout.Rigid(func() {
			if p.totalOffset > 0 {
				p.SliderLayout(gtx)
			}
		}),
	)
}
