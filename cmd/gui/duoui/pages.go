package duoui

import (
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
)

func DuoUIoverview(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if duo.Model.DuoUIcontext.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 0, 0)
		page.Layout(duo.Model.DuoUIcontext, func() {
			viewport.Layout(duo.Model.DuoUIcontext,
				layout.Flexed(0.5, func() {
					cs := duo.Model.DuoUIcontext.Constraints
					helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, cs.Width.Max, cs.Height.Max, duo.Model.DuoUItheme.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					components.DuoUIbalanceWidget(duo.Model, rc)
				}),
				layout.Flexed(0.5, func() {
					cs := duo.Model.DuoUIcontext.Constraints
					helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff424242"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					components.DuoUIlatestTxsWidget(duo.Model, rc)
				}),
			)
		})
	}
}
func DuoUIsend(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.Model.DuoUIcontext, components.DuoUIsend(duo.Model))
	}
}
func DuoUIreceive(duo *DuoUI) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.Model.DuoUIcontext, func() {
			duo.Model.DuoUItheme.H5("receive :").Layout(duo.Model.DuoUIcontext)
		})
	}
}

func DuoUIaddressbook(duo *DuoUI) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.Model.DuoUIcontext, func() {
			duo.Model.DuoUItheme.H5("addressbook :").Layout(duo.Model.DuoUIcontext)
		})
	}
}
func DuoUIsettings(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 0, 0)
		page.Layout(duo.Model.DuoUIcontext, components.DuoUIsettingsWidget(duo.Model, cx, rc))
	}
}

func DuoUInetwork(duo *DuoUI) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.Model.DuoUIcontext, func() {
			duo.Model.DuoUItheme.H5("network :").Layout(duo.Model.DuoUIcontext)
		})
	}
}

func DuoUIhistory(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.Model.DuoUIcontext, components.DuoUItransactionsWidget(duo.Model, cx, rc))
	}
}

func DuoUIconsole(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.Model.DuoUIcontext, components.DuoUIconsoleWidget(duo.Model, cx, rc))
	}
}

func DuoUItrace(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		page := duo.Model.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 10, 10)
		page.Layout(duo.Model.DuoUIcontext, func() {
			DuoUIloader(duo.Model)
		})
	}
}
