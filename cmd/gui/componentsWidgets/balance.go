package componentsWidgets

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
)

var (
	itemsList = &layout.List{
		Axis: layout.Vertical,
	}
	singleItem = &layout.Flex{
		Axis:layout.Horizontal,
		Spacing:layout.SpaceBetween,
	}
	Icon, _ = components.NewDuoUIicon(icons.EditorMonetizationOn)
)
func listItem(duo *models.DuoUI, name, value string){
	layout.Flex{
		Axis:layout.Horizontal,
		Spacing:layout.SpaceBetween,
	}.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(0), Right:unit.Dp(0)}.Layout(duo.DuoUIcontext, func() {
						if Icon != nil {
							Icon.Color = helpers.HexARGB("ff303030")
							Icon.Layout(duo.DuoUIcontext, unit.Px(float32(32)))
						}
						duo.DuoUIcontext.Dimensions = layout.Dimensions{
							Size: image.Point{X: 32, Y: 32},
						}
					})
				}),
				layout.Rigid(func() {
					txt := duo.DuoUItheme.H6(name)
					txt.Color = duo.DuoUIconfiguration.SecondaryTextColor
					txt.Layout(duo.DuoUIcontext)
				}),

			)

		}),
		layout.Rigid(func() {
			value := duo.DuoUItheme.H5(value)
			value.Color = color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}
			value.Alignment = text.End
			value.Layout(duo.DuoUIcontext)
		}),
	)
}
func DuoUIbalanceWidget(duo *models.DuoUI, rc *rcd.RcVar) {
	in := layout.UniformInset(unit.Dp(16))
	in.Layout(duo.DuoUIcontext, func() {
cs := duo.DuoUIcontext.Constraints
		navButtons := []func(){
			func() {
				listItem(duo, "Balance :", rc.Balance + " " + duo.DuoUIconfiguration.Abbrevation)
			},
			func(){
				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max , 1, helpers.HexARGB("ffbdbdbd"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
			},
			func() {
				listItem(duo, "Unconfirmed :", rc.Unconfirmed)
			},
			func(){
				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max , 1, helpers.HexARGB("ffbdbdbd"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
			},
			func() {
				listItem(duo, "Transactions :", fmt.Sprint(rc.Transactions.TxsNumber))
			},
		}
		itemsList.Layout(duo.DuoUIcontext, len(navButtons), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, navButtons[i])
		})
	})
	//
	//in := layout.Inset{
	//	Top:    unit.Dp(15),
	//	Right:  unit.Dp(30),
	//	Bottom: unit.Dp(15),
	//	Left:   unit.Dp(30),
	//}
	//in.Layout(duo.DuoUIcontext, func() {
	//	layout.Flex{
	//		Axis: layout.Vertical,
	//	}.Layout(duo.DuoUIcontext,
	//		layout.Rigid(func() {
	//			balanceTxt := duo.DuoUItheme.H6("Balance :")
	//			balanceTxt.Color = duo.DuoUIconfiguration.SecondaryTextColor
	//			balanceTxt.Layout(duo.DuoUIcontext)
	//		}),
	//		layout.Rigid(func() {
	//			balanceVal := duo.DuoUItheme.H4(rc.Balance + " " + duo.DuoUIconfiguration.Abbrevation)
	//			balanceVal.Color = duo.DuoUIconfiguration.PrimaryTextColor
	//			balanceVal.Alignment = text.End
	//			balanceVal.Layout(duo.DuoUIcontext)
	//		}),
	//		layout.Rigid(func() {
	//			balanceUnconfirmed := duo.DuoUItheme.H6("Unconfirmed :" + rc.Unconfirmed)
	//			balanceUnconfirmed.Color = duo.DuoUIconfiguration.SecondaryTextColor
	//			balanceUnconfirmed.Alignment = text.End
	//			balanceUnconfirmed.Layout(duo.DuoUIcontext)
	//		}),
	//		layout.Rigid(func() {
	//			txsNumber := duo.DuoUItheme.H6("Transactions :" + fmt.Sprint(rc.Transactions.TxsNumber))
	//			txsNumber.Color = duo.DuoUIconfiguration.SecondaryTextColor
	//			txsNumber.Alignment = text.End
	//			txsNumber.Layout(duo.DuoUIcontext)
	//		}),
	//
	//	)
	//
	//})
}
