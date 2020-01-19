package componentsWidgets

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/io/pointer"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"image"
)

var (
	transList = &layout.List{
		Axis: layout.Vertical,
	}
	allTxs      = new(widget.CheckBox)
	mintedTxs   = new(widget.CheckBox)
	immatureTxs = new(widget.CheckBox)
	sentTxs     = new(widget.CheckBox)
	receivedTxs = new(widget.CheckBox)

	increase = &Button{
		Name:         "increase",
		OperateValue: 1,
	}
	decrease = &Button{
		Name:         "decrease",
		OperateValue: 1,
	}
	reset = &Button{
		Name:         "reset",
		OperateValue: 0,
	}
	itemValue = item{
		i: 5,
	}
)

type Button struct {
	pressed      bool
	Name         string
	Do           func(interface{})
	ColorBg      string
	BorderRadius [4]float32
	OperateValue interface{}
}
type item struct {
	i int
}

func (it *item) doIncrease(n int) {
	it.i = it.i + int(n)
}

func (it *item) doDecrease(n int) {
	it.i = it.i - int(n)
}
func (it *item) doReset() {
	it.i = 0
}

func DuoUItransactionsWidget(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	rc.Txs.TxsListNumber = itemValue.i
	rc.GetDuoUITransactionsExcertps(duo, cx)

	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 48, helpers.HexARGB("ff3030cf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			in := layout.UniformInset(unit.Dp(8))
			in.Layout(duo.DuoUIcontext, func() {

				layout.Flex{
					Spacing: layout.SpaceBetween,
				}.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {
						layout.Flex{}.Layout(duo.DuoUIcontext,
							layout.Rigid(func() {
								duo.DuoUItheme.DuoUIcheckBox("All").Layout(duo.DuoUIcontext, allTxs)
							}),
							layout.Rigid(func() {
								duo.DuoUItheme.DuoUIcheckBox("Minted").Layout(duo.DuoUIcontext, mintedTxs)
							}),
							layout.Rigid(func() {
								duo.DuoUItheme.DuoUIcheckBox("Immature").Layout(duo.DuoUIcontext, immatureTxs)
							}),
							layout.Rigid(func() {
								duo.DuoUItheme.DuoUIcheckBox("Sent").Layout(duo.DuoUIcontext, sentTxs)
							}),
							layout.Rigid(func() {
								duo.DuoUItheme.DuoUIcheckBox("Received").Layout(duo.DuoUIcontext, receivedTxs)
							}),
						)
					}),
					layout.Rigid(func() {
						layout.Flex{}.Layout(duo.DuoUIcontext,
							layout.Rigid(func() {

								counter(duo)

							}),
						)
					}),
				)
			})
		}),
		layout.Flexed(1, func() {

			in := layout.UniformInset(unit.Dp(16))
			in.Layout(duo.DuoUIcontext, func() {
				duo.DuoUIcomponents.Status.Layout.Layout(duo.DuoUIcontext,
					// Balance status item
					layout.Rigid(func() {
						cs := duo.DuoUIcontext.Constraints
						//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff424242"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

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
											num.Color = helpers.Alpha(a, num.Color)
											num.Layout(duo.DuoUIcontext)
										}),
										layout.Rigid(func() {
											tim := duo.DuoUItheme.Body1(t.TxID)
											tim.Color = helpers.Alpha(a, tim.Color)
											tim.Layout(duo.DuoUIcontext)
										}),
										layout.Rigid(func() {
											amount := duo.DuoUItheme.H5(fmt.Sprintf("%0.8f", t.Amount))
											amount.Color = helpers.RGB(0x003300)
											amount.Color = helpers.Alpha(a, amount.Color)
											amount.Alignment = text.End
											amount.Font.Variant = "Mono"
											amount.Font.Weight = text.Bold
											amount.Layout(duo.DuoUIcontext)
										}),
										layout.Rigid(func() {
											sat := duo.DuoUItheme.Body1(t.Category)
											sat.Color = helpers.Alpha(a, sat.Color)
											sat.Layout(duo.DuoUIcontext)
										}),
										layout.Rigid(func() {
											l := duo.DuoUItheme.Body2(t.Time)
											l.Color = duo.DuoUItheme.Color.Hint
											l.Color = helpers.Alpha(a, l.Color)
											l.Layout(duo.DuoUIcontext)
										}),
									)
								}),
								layout.Rigid(func() {
									sat := duo.DuoUItheme.Body1(fmt.Sprintf("%0.8f", t.Amount))
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

//////////////////////////
/////////////////////////

func counter(duo *models.DuoUI) {

	layout.Stack{}.Layout(duo.DuoUIcontext,
		layout.Stacked(func() {
			layout.Flex{}.Layout(duo.DuoUIcontext,
				layout.Flexed(0.4, func() {
					decrease.Do = func(n interface{}) {
						itemValue.doDecrease(n.(int))
					}
					decrease.Layout(duo)
				}),
				layout.Flexed(0.2, func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(duo.DuoUIcontext,
						layout.Rigid(func() {
							//cs := duo.DuoUIcontext.Constraints
							//helpers.DrawRectangle(duo.DuoUIcontext, cs.Width.Max, 120, helpers.HexARGB("ff3030cf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							in := layout.UniformInset(unit.Dp(0))
							in.Layout(duo.DuoUIcontext, func() {
								layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
									duo.DuoUItheme.Body2(fmt.Sprint(itemValue.i)).Layout(duo.DuoUIcontext)
								})
							})
						}),
						layout.Flexed(1, func() {
							reset.Do = func(interface{}) {
								itemValue.doReset()
							}
							reset.Layout(duo)
						}),
					)
				}),
				layout.Flexed(0.4, func() {
					increase.Do = func(n interface{}) {
						itemValue.doIncrease(n.(int))
					}
					increase.Layout(duo)
				}),
			)
		}),
	)
}

func (b *Button) Layout(duo *models.DuoUI) {
	for _, e := range duo.DuoUIcontext.Events(b) { // HLevent
		if e, ok := e.(pointer.Event); ok { // HLevent
			switch e.Type { // HLevent
			case pointer.Press: // HLevent
				b.pressed = true // HLevent
				b.Do(b.OperateValue)
			case pointer.Release: // HLevent
				b.pressed = false // HLevent
			}
		}
	}

	cs := duo.DuoUIcontext.Constraints
	colorBg := helpers.HexARGB("ff30cfcf")
	colorBorder := helpers.HexARGB("ffcf3030")
	border := unit.Dp(1)

	if b.pressed {
		colorBg = helpers.HexARGB("ffcf30cf")
		colorBorder = helpers.HexARGB("ff303030")
		border = unit.Dp(3)
	}
	pointer.Rect( // HLevent
		image.Rectangle{Max: image.Point{X: cs.Width.Max, Y: cs.Height.Max}}, // HLevent
	).Add(duo.DuoUIcontext.Ops)                       // HLevent
	pointer.InputOp{Key: b}.Add(duo.DuoUIcontext.Ops) // HLevent
	//helpers.DrawRectangle(gtx, cs.Width.Max, cs.Height.Max, colorBorder, b.BorderRadius, unit.Dp(0))
	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 32, colorBorder, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

	in := layout.UniformInset(border)
	in.Layout(duo.DuoUIcontext, func() {
		//helpers.DrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, colorBg, b.BorderRadius, unit.Dp(0))
		helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 30, colorBg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		//cs := gtx.Constraints
		label := duo.DuoUItheme.Caption(b.Name)
		label.Alignment = text.Middle
		label.Layout(duo.DuoUIcontext)
	})
}
