package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op/clip"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/unit"
	"image"
	"image/color"
)


func DuoUIframe(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar, bgColor string,  padding [4]float32, borderRadius [4]float32, frame func()) {

	//brcol := helpers.HexARGB("ff303030")

	cs := duo.DuoUIcontext.Constraints
	hmin := cs.Width.Max
	vmin := cs.Height.Max
	layout.Stack{Alignment: layout.Center}.Layout(duo.DuoUIcontext,
		layout.Expanded(func() {
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(cs.Width.Min),
					Y: float32(cs.Height.Min),
				}},
			}.Op(duo.DuoUIcontext.Ops).Add(duo.DuoUIcontext.Ops)
		}),
		layout.Stacked(func() {
			cs.Width.Min = hmin
			cs.Height.Min = vmin
			layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
				layout.Inset{Top: unit.Dp(padding[0]), Bottom: unit.Dp(padding[1]), Left: unit.Dp(padding[2]), Right: unit.Dp(padding[3])}.Layout(duo.DuoUIcontext, func() {

					layout.Stack{Alignment: layout.Center}.Layout(duo.DuoUIcontext,
						layout.Expanded(func() {
							//rr := float32(duo.DuoUIcontext.Px(unit.Dp(4)))
							clip.Rect{
								Rect: f32.Rectangle{Max: f32.Point{
									X: float32(cs.Width.Min),
									Y: float32(cs.Height.Min),
								}},
								NE: borderRadius[0], NW: borderRadius[1], SE: borderRadius[2], SW: borderRadius[3],
							}.Op(duo.DuoUIcontext.Ops).Add(duo.DuoUIcontext.Ops)
							fill(duo.DuoUIcontext, helpers.HexARGB(bgColor))
						}),
						layout.Stacked(func() {
							cs.Width.Min = hmin
							cs.Height.Min = vmin
							layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
								//layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15), Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(duo.DuoUIcontext, func() {

								// Content <<<
								frame()

								// Content >>>

								//})
							})
						}),
					)
				})
			})
		}),
	)

}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
