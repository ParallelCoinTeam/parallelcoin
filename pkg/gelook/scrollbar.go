package gelook

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/pkg/gel"
	"image"
)

type ScrollBar struct {
	ColorBg      string
	BorderRadius [4]float32
	OperateValue interface{}
	//Height       float32
	body *ScrollBarBody
	up   *ScrollBarButton
	down *ScrollBarButton
}

type ScrollBarBody struct {
	ColorBg      string
	Position     int
	Cursor       int
	Height       int
	CursorHeight int
	Icon         DuoUIicon
}

type ScrollBarButton struct {
	button      DuoUIbutton
	Height      int
	insetTop    float32
	insetRight  float32
	insetBottom float32
	insetLeft   float32
	iconSize    int
	iconPadding float32
}

func (t *DuoUItheme) ScrollBar() *ScrollBar {
	//itemValue := item{
	//	i: 0,
	//}
	up := &ScrollBarButton{
		button: t.DuoUIbutton(t.Fonts["Primary"], "", "", "ff669944", "", "", "Up", "ff558822", 0, 42, 0, 0, 0, 0),
		//button:      t.DuoUIbutton(t.Icons["Up"]),
		insetTop:    0,
		insetRight:  0,
		insetBottom: 0,
		insetLeft:   0,
		//iconSize:    p.size,
		iconPadding: 0,
	}
	down := &ScrollBarButton{
		button: t.DuoUIbutton(t.Fonts["Primary"], "", "", "ff669944", "", "", "Down", "ff558822", 0, 42, 0, 0, 0, 0),
		//button:      t.DuoUIbutton(t.Icons["Down"]),
		Height:   16,
		iconSize: 16,
	}
	body := &ScrollBarBody{
		ColorBg:  "",
		Position: 0,
		Cursor:   0,
		Icon:     *t.Icons["Grab"],
	}
	return &ScrollBar{
		ColorBg:      "ff447733",
		BorderRadius: [4]float32{},
		OperateValue: 1,
		//ListPosition: 0,
		//Height: 16,
		body: body,
		up:   up,
		down: down,
	}
}
func (s *ScrollBarButton) scrollBarButton() *DuoUIbutton {
	button := s.button
	//button.Inset.Top = unit.Dp(0)
	//button.Inset.Bottom = unit.Dp(0)
	//button.Inset.Right = unit.Dp(0)
	//button.Inset.Left = unit.Dp(0)
	//button.Size = unit.Dp(32)
	//button.Padding = unit.Dp(0)
	return &button
}

func (p *Panel) SliderLayout(gtx *layout.Context, panel *gel.Panel) {
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func() {
			for panel.ScrollBar.Up.Clicked(gtx) {
				if panel.PanelContentLayout.Position.First > 0 {
					//p.panelContent.Position.First = p.panelContent.Position.First - int(p.ScrollBar.body.CursorHeight)
					panel.PanelContentLayout.Position.First = panel.PanelContentLayout.Position.First - 1
					panel.PanelContentLayout.Position.Offset = 0
				}
			}
			p.ScrollBar.up.button.IconLayout(gtx, panel.ScrollBar.Up)
		}),
		layout.Flexed(1, func() {
			p.bodyLayout(gtx, panel)
		}),
		layout.Rigid(func() {
			for panel.ScrollBar.Down.Clicked(gtx) {
				if panel.PanelContentLayout.Position.BeforeEnd {
					//p.panelContent.Position.First = p.panelContent.Position.First + int(p.ScrollBar.body.CursorHeight)
					panel.PanelContentLayout.Position.First = panel.PanelContentLayout.Position.First + 1
					panel.PanelContentLayout.Position.Offset = 0
				}
			}
			p.ScrollBar.down.button.BgColor = HexARGB("ff004455")
			p.ScrollBar.down.button.IconLayout(gtx, panel.ScrollBar.Down)
		}),
	)
	panel.Layout(gtx)
}

func (p *Panel) bodyLayout(gtx *layout.Context, panel *gel.Panel) {
	//for _, e := range gtx.Events(p.ScrollBar.body) {
	//	if e, ok := e.(pointer.Event); ok {
	//		if e.Position.Y > 0 {
	//			p.ScrollBar.body.Position = int(e.Position.Y) - (p.ScrollBar.body.CursorHeight / 2)
	//		}
	//		switch e.Type {
	//		case pointer.Press:
	//			p.ScrollBar.body.pressed = true
	//			p.ScrollBar.body.Do(p.ScrollBar.body.OperateValue)
	//		case pointer.Release:
	//			p.ScrollBar.body.pressed = false
	//		}
	//	}
	//}
	cs := gtx.Constraints
	p.ScrollBar.body.Height = cs.Height.Max
	sliderBg := "ff005588"
	colorBg := "ff4477cf"
	colorBorder := "ffcf3030"
	border := unit.Dp(0)
	//if panel.ScrollBar.Body.Pressed {
	//	if p.ScrollBar.body.Position >= 0 && p.ScrollBar.body.Position <= cs.Height.Max-p.ScrollBar.body.CursorHeight {
	//		p.ScrollBar.body.Cursor = p.ScrollBar.body.Position
	//		p.PanelContentLayout.Position.First = p.ScrollBar.body.Position / p.ScrollUnit
	//		p.PanelContentLayout.Position.Offset = 0
	//		//p.panelContent.Position.First = int(p.ScrollBar.body.Cursor)
	//	}
	//	colorBg = "ffcf30cf"
	//	colorBorder = "ff303030"
	//	border = unit.Dp(0)
	//}
	pointer.Rect(
		image.Rectangle{Max: image.Point{X: cs.Width.Max, Y: cs.Height.Max}},
	).Add(gtx.Ops)
	pointer.InputOp{Key: p.ScrollBar.body}.Add(gtx.Ops)
	DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, colorBorder, [4]float32{0, 0, 0, 0}, [4]float32{8, 8, 8, 8})
	layout.UniformInset(border).Layout(gtx, func() {
		cs := gtx.Constraints
		DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, colorBg, [4]float32{0, 0, 0, 0}, [4]float32{8, 8, 8, 8})
		//cs := gtx.Constraints
		layout.Flex{
			Axis: layout.Vertical,
			//Alignment:layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.Center.Layout(gtx, func() {
					layout.Inset{
						Top: unit.Dp(float32(panel.PanelContentLayout.Position.First * panel.ScrollUnit)),
					}.Layout(gtx, func() {
						//gtx.Dimensions.Size.Y= p.ScrollBar.body.CursorHeight
						gtx.Constraints.Height.Min = 50
						DuoUIdrawRectangle(gtx, 50, 50, sliderBg, [4]float32{8, 8, 8, 8}, [4]float32{8, 8, 8, 8})
						layout.Center.Layout(gtx, func() {
							p.ScrollBar.body.Icon.Color = HexARGB("ff554499")
							p.ScrollBar.body.Icon.Layout(gtx, unit.Px(float32(p.Size)))
						})
					})
				})
			}),
		)
	})
}
