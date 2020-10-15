package counter

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
	"image"
)

type (
	D = layout.Dimensions
	C = layout.Context
	W = layout.Widget
)
type CounterStyle struct {
	increase   *iconBtn
	decrease   *iconBtn
	reset      *iconBtn
	input      material.EditorStyle
	Font       text.Font
	TextSize   unit.Value
	TxColor    string
	BgColor    string
	BtnBgColor string
	shaper     text.Shaper
	c          *Counter
}

func CounterSt(t *theme.Theme, cc *Counter) CounterStyle {
	return CounterStyle{
		// ToDo Replace theme's buttons with counter exclusive buttons, set icons for increase/decrease
		increase: iconButton(t.T, t.Icons["Down"], cc.CounterIncrease, t.Colors["Light"]),
		decrease: iconButton(t.T, t.Icons["Up"], cc.CounterDecrease, t.Colors["Light"]),
		input:    material.Editor(t.T, cc.CounterInput, "0"),
		//pageFunction: pageFunction,
		Font: text.Font{
			Typeface: t.Fonts["Primary"],
		},
		TxColor:    t.Colors["Light"],
		BgColor:    t.Colors["Dark"],
		BtnBgColor: t.Colors["Primary"],
		TextSize:   unit.Dp(float32(18)),
		shaper:     t.T.Shaper,
		c:          cc,
	}
}

func (c CounterStyle) Layout(th *theme.Theme, value string) func(gtx C) D {

	return func(gtx C) D {
		return lyt.Format(gtx, "hflex(middle,r(_),r(_))",
			box.BoxEditor(th, func(gtx C) D {
				e := material.Editor(th.T, c.c.CounterInput, "0.00000000")
				return e.Layout(gtx)
			}),
			func(gtx C) D {
				return lyt.Format(gtx, "vflex(middle,r(_),r(_))",
					func(gtx C) D {
						for c.c.CounterDecrease.Clicked() {
							c.c.Decrease()
							c.c.PageFunction()
						}
						return c.decrease.Layout(gtx)
					},
					func(gtx C) D {
						for c.c.CounterIncrease.Clicked() {
							c.c.Increase()
							c.c.PageFunction()
						}
						return c.increase.Layout(gtx)
					})
			})
	}
	//return func(gtx layout.Context) layout.Dimensions {
	//	c.c.CounterInput.SetText(value)
	//	bgColor := c.BgColor
	//	return layout.Stack{Alignment: layout.Center}.Layout(g,
	//		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
	//			return helper.Fill(g, helper.HexARGB(bgColor))
	//		}),
	//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	//			return layout.Center.Layout(g, func(gtx layout.Context) layout.Dimensions {
	//				return layout.Flex{
	//					Spacing:   layout.SpaceAround,
	//					Alignment: layout.Middle,
	//				}.Layout(gtx,
	//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//						for c.c.CounterDecrease.Clicked() {
	//							c.c.Decrease()
	//							c.c.PageFunction()
	//						}
	//						return c.decrease.Layout(gtx)
	//					}),
	//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//							return layout.Inset{
	//								Top:    unit.Dp(0),
	//								Right:  unit.Dp(16),
	//								Bottom: unit.Dp(0),
	//								Left:   unit.Dp(16),
	//							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//								return layout.Flex{
	//									Axis:      layout.Vertical,
	//									Spacing:   layout.SpaceAround,
	//									Alignment: layout.Middle,
	//								}.Layout(gtx,
	//									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//										c.input.Font.Typeface = c.Font.Typeface
	//										c.input.Color = helper.HexARGB(c.TxColor)
	//										for _, e := range c.c.CounterInput.Events() {
	//											switch e.(type) {
	//											case widget.ChangeEvent:
	//												//if i, err := strconv.Atoi(c.c.CounterInput.Text()); err == nil {
	//												//	c.c.Value = i
	//												//}
	//											}
	//										}
	//										return c.input.Layout(gtx)
	//										paint.ColorOp{Color: helper.HexARGB(c.TxColor)}.Add(gtx.Ops)
	//										return material.Body1(th.T, value).Layout(gtx)
	//									}))
	//							})
	//						})
	//					}),
	//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//						for c.c.CounterIncrease.Clicked() {
	//							c.c.Increase()
	//							c.c.PageFunction()
	//						}
	//						return c.increase.Layout(gtx)
	//					}))
	//			})
	//		}),
	//	)
	//
	//}
}

type iconBtn struct {
	theme  *material.Theme
	button *widget.Clickable
	icon   *widget.Icon
	word   string
	bg     string
}

func iconButton(t *material.Theme, i *widget.Icon, c *widget.Clickable, bg string) *iconBtn {
	return &iconBtn{
		theme:  t,
		button: c,
		icon:   i,
		bg:     bg,
	}
}

func (b iconBtn) Layout(gtx layout.Context) layout.Dimensions {
	btn := material.ButtonLayout(b.theme, b.button)
	btn.CornerRadius = unit.Dp(0)
	btn.Background = helper.HexARGB(b.bg)
	return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			layIcon := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					var d layout.Dimensions
					if b.icon != nil {
						size := gtx.Px(unit.Dp(12))
						b.icon.Layout(gtx, unit.Px(float32(size)))
						d = layout.Dimensions{
							Size: image.Point{X: size, Y: size},
						}
					}
					return d
				})
			})

			return iconAndLabel.Layout(gtx, layIcon)
		})
	})
}
