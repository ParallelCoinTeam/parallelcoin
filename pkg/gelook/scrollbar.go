package gelook

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/gel"
	"github.com/p9c/logi"
)

var (
	widgetButtonUp   = new(gel.Button)
	widgetButtonDown = new(gel.Button)
)

type ScrollBar struct {
	ColorBg      string
	BorderRadius [4]float32
	OperateValue interface{}
	// Height       float32
	// scrollUnit     float32
	// positionOffset int
	// cursor         float32
	size       int
	body       *ScrollBarBody
	up         *ScrollBarButton
	down       *ScrollBarButton
	controller *gel.ScrollBar
}

type ScrollBarBody struct {
	ColorBg string
	Height  int
	Icon    DuoUIicon
}

type ScrollBarButton struct {
	button      *IconButton
	Height      int
	iconColor   string
	iconBgColor string
	insetTop    float32
	insetRight  float32
	insetBottom float32
	insetLeft   float32
	size        int
	iconPadding float32
}

func (t *DuoUItheme) ScrollBar(c *gel.ScrollBar) *ScrollBar {
	buttonUp := t.IconButton(t.Icons["iconUp"])
	buttonDown := t.IconButton(t.Icons["iconOK"])
	up := &ScrollBarButton{
		button: &buttonUp,
		// Height: w,
		// size:   w,
	}
	down := &ScrollBarButton{
		button: &buttonDown,
		// Height: w,
		// size:   w,
	}
	body := &ScrollBarBody{
		ColorBg: "ff445588",
		Icon:    *t.Icons["iconGrab"],
	}
	return &ScrollBar{
		size:         t.scrollBarSize,
		ColorBg:      "ff885566",
		BorderRadius: [4]float32{},
		OperateValue: 1,
		// ListPosition: 0,
		// Height: 16,
		// scrollUnit:     scrollUnit,
		// positionOffset: positionOffset,
		// cursor:         cursor,
		controller: c,
		body:       body,
		up:         up,
		down:       down,
	}
}

func (s *ScrollBar) Layout(gtx *layout.Context, positionOffset int, scrollUnit float32) {
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func() {
			for widgetButtonUp.Clicked(gtx) {
				positionOffset = positionOffset - int(s.controller.CursorHeight)
			}
			s.up.scrollBarButton(s.size).Layout(gtx, widgetButtonUp)
		}),
		layout.Flexed(1, s.bodyLayout(gtx,
			positionOffset,
			scrollUnit)),
		layout.Rigid(func() {
			for widgetButtonDown.Clicked(gtx) {
				positionOffset = positionOffset + int(s.controller.CursorHeight)
			}
			s.down.scrollBarButton(s.size).Layout(gtx, widgetButtonDown)
		}),
	)
}

func (s *ScrollBarButton) scrollBarButton(size int) *IconButton {
	button := *s.button
	// button..Inset.Top = unit.Dp(0)
	// button.Inset.Bottom = unit.Dp(0)
	// button.Inset.Right = unit.Dp(0)
	// button.Inset.Left = unit.Dp(0)
	button.Background = HexARGB("ffff0000")
	button.Icon.Color = HexARGB("ff882266")
	button.Icon.size = unit.Dp(float32(size))
	button.Size = unit.Dp(float32(size))
	button.Padding = unit.Dp(0)
	return &button
}
func (s *ScrollBar) bodyLayout(gtx *layout.Context, positionOffset int, scrollUnit float32) func() {
	return func() {

		cs := gtx.Constraints
		sliderBg := "ff558899"
		colorBg := "ff30cfcf"
		colorBorder := "ffcf3030"
		border := unit.Dp(0)
		// if s.body.pressed {
		if s.controller.Position >= 0 && s.controller.Position <= float32(cs.Height.Max-s.controller.CursorHeight) {
			s.controller.Cursor = s.controller.Position
			positionOffset = int(float32(s.controller.Cursor) / scrollUnit)
		}
		colorBg = "ffcf30cf"
		colorBorder = "ff303030"
		border = unit.Dp(0)
		// }
		pointer.Rect(
			image.Rectangle{Max: image.Point{X: cs.Width.Max, Y: cs.Height.Max}},
		).Add(gtx.Ops)
		pointer.InputOp{Key: s.body}.Add(gtx.Ops)
		DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, colorBorder, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.UniformInset(border).Layout(gtx, func() {
			cs := gtx.Constraints
			DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, colorBg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			// cs := gtx.Constraints
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					layout.Center.Layout(gtx, func() {
						layout.Inset{
							Top: unit.Dp(s.controller.Cursor),
						}.Layout(gtx, func() {
							// cs := gtx.Constraints
							if s.controller.CursorHeight > s.size {
								s.body.Height = s.controller.CursorHeight
							} else {
								s.body.Height = s.size
							}

							DuoUIdrawRectangle(gtx, s.size, s.controller.CursorHeight, sliderBg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							// DuoUIdrawRectangle(gtx, 30, 111, sliderBg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

							layout.Center.Layout(gtx, func() {
								s.body.Icon.Color = HexARGB("ff554499")
								s.body.Icon.Layout(gtx, unit.Px(float32(32)))
							})
						})
					})
				}),
			)
		})
		s.controller.Layout(gtx)
		logi.L.Info("RADI Constraints")
		logi.L.Info(s.controller.BodyHeight)
	}
}
