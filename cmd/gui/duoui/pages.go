package duoui

import (
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
)

func DuoUIoverview(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if duo.m.DuoUIcontext.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 0, 0)
		page.Layout(duo.m.DuoUIcontext, func() {
			viewport.Layout(duo.m.DuoUIcontext,
				layout.Flexed(0.5, func() {
					cs := duo.m.DuoUIcontext.Constraints
					helpers.DuoUIdrawRectangle(duo.m.DuoUIcontext, cs.Width.Max, cs.Height.Max, duo.m.DuoUItheme.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					components.DuoUIbalanceWidget(duo.m, rc)
				}),
				layout.Flexed(0.5, func() {
					cs := duo.m.DuoUIcontext.Constraints
					helpers.DuoUIdrawRectangle(duo.m.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff424242"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					components.DuoUIlatestTxsWidget(duo.m, rc)
				}),
			)
		})
	}
}
func DuoUIsend(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.m.DuoUIcontext, components.DuoUIsend(duo.m))
	}
}
func DuoUIreceive(duo *DuoUI) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.m.DuoUIcontext, func() {
			duo.m.DuoUItheme.H5("receive :").Layout(duo.m.DuoUIcontext)
		})
	}
}

func DuoUIaddressbook(duo *DuoUI) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.m.DuoUIcontext, func() {
			duo.m.DuoUItheme.H5("addressbook :").Layout(duo.m.DuoUIcontext)
		})
	}
}
func DuoUIsettings(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 0, 0)
		page.Layout(duo.m.DuoUIcontext, components.DuoUIsettingsWidget(duo.m, cx, rc))
	}
}

func DuoUInetwork(duo *DuoUI) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.m.DuoUIcontext, func() {
			duo.m.DuoUItheme.H5("network :").Layout(duo.m.DuoUIcontext)
		})
	}
}

func DuoUIhistory(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.m.DuoUIcontext, components.DuoUItransactionsWidget(duo.m, cx, rc))
	}
}

func DuoUIconsole(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.m.DuoUIcontext, components.DuoUIconsoleWidget(duo.m, cx, rc))
	}
}

func DuoUItrace(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.m.DuoUIcontext, func() {
		page := duo.m.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.m.DuoUIcontext, func() {
			DuoUIloader(duo.m)
		})
	}
}
