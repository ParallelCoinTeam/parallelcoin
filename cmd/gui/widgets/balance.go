package widgets

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
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

)

func DuoUIbalanceWidget(duo *models.DuoUI, rc *rcd.RcVar) {

	navButtons := []func(){
		func() {
			singleItem.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
									balanceTxt := duo.DuoUItheme.H6("Balance :")
									balanceTxt.Color = duo.DuoUIconfiguration.SecondaryTextColor
									balanceTxt.Layout(duo.DuoUIcontext)
				}),
				layout.Rigid(func() {
					balance := duo.DuoUItheme.Body2(rc.Balance + " " + duo.DuoUIconfiguration.Abbrevation)
					balance.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					balance.Alignment = text.End
					balance.Layout(duo.DuoUIcontext)
				}),
			)
		},
		func() {

			singleItem.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					balanceUnconfirmedTxt := duo.DuoUItheme.H6("Unconfirmed:")
					balanceUnconfirmedTxt.Color = duo.DuoUIconfiguration.SecondaryTextColor
					balanceUnconfirmedTxt.Alignment = text.End
					balanceUnconfirmedTxt.Layout(duo.DuoUIcontext)
				}),
				layout.Rigid(func() {
					balanceUnconfirmed := duo.DuoUItheme.H6(rc.Unconfirmed)
					balanceUnconfirmed.Color = duo.DuoUIconfiguration.SecondaryTextColor
					balanceUnconfirmed.Alignment = text.End
					balanceUnconfirmed.Layout(duo.DuoUIcontext)
				}),
			)



		},
		func() {

			singleItem.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					txsNumberTxt := duo.DuoUItheme.H6("Transactions :")
					txsNumberTxt.Color = duo.DuoUIconfiguration.SecondaryTextColor
					txsNumberTxt.Alignment = text.End
					txsNumberTxt.Layout(duo.DuoUIcontext)
				}),
				layout.Rigid(func() {
					txsNumber := duo.DuoUItheme.H6(fmt.Sprint(rc.Transactions.TxsNumber))
					txsNumber.Color = duo.DuoUIconfiguration.SecondaryTextColor
					txsNumber.Alignment = text.End
					txsNumber.Layout(duo.DuoUIcontext)
				}),
			)

		},
	}
	itemsList.Layout(duo.DuoUIcontext, len(navButtons), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, navButtons[i])
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
