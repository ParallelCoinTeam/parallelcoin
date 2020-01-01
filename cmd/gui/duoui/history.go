package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/widgets"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIhistory(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(1, func() {


			DuoUIframe(duo,cx,rc,"ff558866", [4]float32{20, 50, 40, 100},[4]float32{0, 0, 0, 0} ,func(){
			widgets.DuoUItransactionsWidget(duo, cx, rc)

			})
		}),
	)
}
