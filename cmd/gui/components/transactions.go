package components

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
)

var (
	transList = &layout.List{
		Axis: layout.Vertical,
	}
	allTxs              = new(widget.CheckBox)
	mintedTxs           = new(widget.CheckBox)
	immatureTxs         = new(widget.CheckBox)
	sentTxs             = new(widget.CheckBox)
	receivedTxs         = new(widget.CheckBox)
	transactionsCounter = new(widget.Counter)
)

func DuoUItransactionsWidget(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) func() {
	return func() {
		rc.Txs.TxsListNumber = 55
		rc.GetDuoUITransactionsExcertps(cx)

		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(duo.DuoUIcontext,
			layout.Rigid(func() {
				cs := duo.DuoUIcontext.Constraints
				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 48, duo.DuoUItheme.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

				in := layout.UniformInset(unit.Dp(8))
				in.Layout(duo.DuoUIcontext, func() {

					layout.Flex{
						Spacing: layout.SpaceBetween,
					}.Layout(duo.DuoUIcontext,
						layout.Rigid(func() {
							layout.Flex{}.Layout(duo.DuoUIcontext,
								layout.Rigid(func() {
									duo.DuoUItheme.DuoUIcheckBox("ALL").Layout(duo.DuoUIcontext, allTxs)
								}),
								layout.Rigid(func() {
									duo.DuoUItheme.DuoUIcheckBox("MINTED").Layout(duo.DuoUIcontext, mintedTxs)
								}),
								layout.Rigid(func() {
									duo.DuoUItheme.DuoUIcheckBox("IMATURE").Layout(duo.DuoUIcontext, immatureTxs)
								}),
								layout.Rigid(func() {
									duo.DuoUItheme.DuoUIcheckBox("SENT").Layout(duo.DuoUIcontext, sentTxs)
								}),
								layout.Rigid(func() {
									duo.DuoUItheme.DuoUIcheckBox("RECEIVED").Layout(duo.DuoUIcontext, receivedTxs)
								}),
							)
						}),
						layout.Rigid(func() {
							layout.Flex{}.Layout(duo.DuoUIcontext,
								layout.Rigid(func() {

									//parallel.DuoUIcounter(duo)

								}),
							)
						}),
					)
				})
			}),
			layout.Flexed(1, func() {

				in := layout.UniformInset(unit.Dp(16))
				in.Layout(duo.DuoUIcontext, func() {
					layout.Flex{
						Axis: layout.Vertical,
					}.Layout(duo.DuoUIcontext,
						// Balance status item
						layout.Rigid(func() {
							cs := duo.DuoUIcontext.Constraints
							//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ff424242", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

							//const n = 5
							//list.Layout(duo.DuoUIcontext, n, func(i int) {
							//	txt := fmt.Sprintf("List element #%d", i)
							//
							//	duo.DuoUItheme.H3(txt).Layout(duo.DuoUIcontext)
							//})
							//transList := &layout.List{
							//	Axis: layout.Vertical,
							//}

							//amount := duo.DuoUItheme.H5(fmt.Sprintf("%0.8f", rc.Txs.Txs))
							//amount.Color = helpers.RGB(0x003300)
							//amount.Color = helpers.Alpha(1.0, amount.Color)
							//amount.Alignment = text.End
							//amount.Font.Variant = "Mono"
							//amount.Font.Weight = text.Bold
							//amount.Layout(duo.DuoUIcontext)

							transList.Layout(duo.DuoUIcontext, len(rc.Txs.Txs), func(i int) {
								// Invert list
								//i = len(txs.Txs) - 1 - i
								//
								t := rc.Txs.Txs[i]
								a := 1.0
								//const duration = 5
								helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 1, helpers.HexARGB("ff535353"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

								layout.Flex{
									Spacing: layout.SpaceBetween,
								}.Layout(duo.DuoUIcontext,
									layout.Rigid(func() {
										layout.Flex{
											Axis: layout.Vertical,
										}.Layout(duo.DuoUIcontext,
											layout.Rigid(func() {
												num := duo.DuoUItheme.Body1(fmt.Sprint(i))
												num.Font.Typeface = "bariol"
												num.Color = helpers.Alpha(a, num.Color)
												num.Layout(duo.DuoUIcontext)
											}),
											layout.Rigid(func() {
												tim := duo.DuoUItheme.Body1(t.TxID)
												tim.Font.Typeface = "bariol"
												tim.Color = helpers.Alpha(a, tim.Color)
												tim.Layout(duo.DuoUIcontext)
											}),
											layout.Rigid(func() {
												amount := duo.DuoUItheme.H5(fmt.Sprintf("%0.8f", t.Amount))
												amount.Font.Typeface = "bariol"
												amount.Color = helpers.RGB(0x003300)
												amount.Color = helpers.Alpha(a, amount.Color)
												amount.Alignment = text.End
												amount.Font.Variant = "Mono"
												amount.Font.Weight = text.Bold
												amount.Layout(duo.DuoUIcontext)
											}),
											layout.Rigid(func() {
												sat := duo.DuoUItheme.Body1(t.Category)
												sat.Font.Typeface = "bariol"
												sat.Color = helpers.Alpha(a, sat.Color)
												sat.Layout(duo.DuoUIcontext)
											}),
											layout.Rigid(func() {
												l := duo.DuoUItheme.Body2(t.Time)
												l.Font.Typeface = "bariol"
												l.Color = duo.DuoUItheme.Color.Hint
												l.Color = helpers.Alpha(a, l.Color)
												l.Layout(duo.DuoUIcontext)
											}),
										)
									}),
									layout.Rigid(func() {
										sat := duo.DuoUItheme.Body1(fmt.Sprintf("%0.8f", t.Amount))
										sat.Font.Typeface = "bariol"
										sat.Color = helpers.Alpha(a, sat.Color)
										sat.Layout(duo.DuoUIcontext)
									}),
								)
							})
						}))
				})
			}),
		)
	}
}