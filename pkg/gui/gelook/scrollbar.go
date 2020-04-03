package gelook

import (
	"fmt"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/pkg/gui/gel"
	"image"
)

var (
	widgetButtonUp   = new(gel.Button)
	widgetButtonDown = new(gel.Button)
)

type ScrollBar struct {
	ColorBg      string
	BorderRadius [4]float32
	OperateValue interface{}
	slider       *ScrollBarSlider
	up           IconButton
	down         IconButton
	container    DuoUIcontainer
}

type ScrollBarSlider struct {
	container DuoUIcontainer
	Icon      DuoUIicon
}

func (t *DuoUItheme) ScrollBar() *ScrollBar {
	slider := &ScrollBarSlider{
		container: t.DuoUIcontainer(0, t.Colors["Primary"]),
		Icon:      *t.Icons["Grab"],
	}
	slider.container.CornerRadius = 8
	scrollbar := &ScrollBar{
		ColorBg:      t.Colors["DarkGrayII"],
		BorderRadius: [4]float32{},
		slider:       slider,
		up:           t.IconButton(t.Icons["Up"]),
		down:         t.IconButton(t.Icons["Down"]),
		container:    t.DuoUIcontainer(0, t.Colors["Light"]),
	}
	scrollbar.container.PaddingLeft = 16
	return scrollbar
}

func (p *DuoUIpanel) ScrollBarLayout(gtx *layout.Context, panel *gel.Panel) {
	p.ScrollBar.container.Layout(gtx, layout.Center, func() {

		layout.Flex{
			Axis: layout.Vertical,
			//Alignment:layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func() {
				for widgetButtonUp.Clicked(gtx) {
					if panel.PanelContentLayout.Position.First > 0 {
						//p.panelContent.Position.First = p.panelContent.Position.First - int(p.ScrollBar.body.CursorHeight)
						panel.PanelContentLayout.Position.First = panel.PanelContentLayout.Position.First - 1
						panel.PanelContentLayout.Position.Offset = 0
					}
				}
				p.ScrollBar.up.Padding = unit.Dp(0)
				p.ScrollBar.up.Size = unit.Dp(16)
				p.ScrollBar.up.Color = HexARGB("ffcfcfcf")
				p.ScrollBar.up.Layout(gtx, widgetButtonUp)
			}),
			layout.Flexed(1, func() {
				p.bodyLayout(gtx, panel)
			}),
			layout.Rigid(func() {
				for widgetButtonDown.Clicked(gtx) {
					if panel.PanelContentLayout.Position.BeforeEnd {
						//p.panelContent.Position.First = p.panelContent.Position.First + int(p.ScrollBar.body.CursorHeight)
						panel.PanelContentLayout.Position.First = panel.PanelContentLayout.Position.First + 1
						panel.PanelContentLayout.Position.Offset = 0
					}
				}
				p.ScrollBar.down.Padding = unit.Dp(0)
				p.ScrollBar.down.Size = unit.Dp(16)
				p.ScrollBar.down.Color = HexARGB("ffcfcfcf")
				p.ScrollBar.down.Layout(gtx, widgetButtonDown)
			}),
		)
	})
}

func (p *DuoUIpanel) bodyLayout(gtx *layout.Context, panel *gel.Panel) {

	layout.Flex{
		Axis: layout.Vertical,
		//Alignment:layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func() {
			cs := gtx.Constraints

			pointer.Rect(
				image.Rectangle{Max: image.Point{X: cs.Width.Max, Y: cs.Height.Max}},
			).Add(gtx.Ops)
			pointer.InputOp{Key: panel.ScrollBar.Slider}.Add(gtx.Ops)
			//DuoUIdrawRectangle(gtx, 16, 30, p.ScrollBar.ColorBg, [4]float32{0, 0, 0, 0}, [4]float32{8, 8, 8, 8})
			//cs := gtx.Constraints
			layout.Center.Layout(gtx, func() {

				layout.Inset{
					Top: unit.Dp(float32(panel.PanelContentLayout.Position.First) * panel.ScrollUnit),
				}.Layout(gtx, func() {

					gtx.Constraints.Width.Min = panel.ScrollBar.Size
					//gtx.Constraints.Width.Max = panel.ScrollBar.Size
					gtx.Constraints.Height.Min = panel.ScrollBar.Slider.CursorHeight
					fmt.Println(".ScrollBar.Slider.CursorHeight", panel.ScrollBar.Slider.CursorHeight)
					fmt.Println(" ScrollBar.Size", panel.ScrollBar.Size)
					if panel.ScrollBar.Slider.CursorHeight < panel.ScrollBar.Size {
						panel.ScrollBar.Slider.CursorHeight = panel.ScrollBar.Size * 2
					}
					//DuoUIdrawRectangle(gtx, panel.ScrollBar.Size, panel.ScrollBar.Slider.CursorHeight, p.ScrollBar.slider.ColorBg, [4]float32{8, 8, 8, 8}, [4]float32{8, 8, 8, 8})
					//layout.Center.Layout(gtx, func() {
					p.ScrollBar.slider.container.Layout(gtx, layout.Center, func() {

						p.ScrollBar.slider.Icon.Color = HexARGB("ffcfcfcf")
						//p.ScrollBar.slider.Icon.op.Rect.Inset(0)
						p.ScrollBar.slider.Icon.Layout(gtx, unit.Px(float32(panel.ScrollBar.Size/2)))
					})
					//})
				})
				panel.ScrollBar.Slider.Layout(gtx)
			})
		}),
	)
}
