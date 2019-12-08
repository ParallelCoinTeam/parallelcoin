// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/skip2/go-qrcode"

	xlayout "eliasnaur.com/giox/layout"

	"gioui.org/font/gofont"
)

type fill color.RGBA

var (
	theme *material.Theme
)

func init() {
	gofont.Register()
	theme = material.NewTheme()
	theme.Color.Primary = rgb(0x1b5e20)
	//theme.TextSize = unit.Sp(20)
}

type App struct {
	wallet    *Wallet
	trans     []Transaction
	transList *layout.List
	qrBtn     widget.Button
	showQR    bool
	addrQR    paint.ImageOp
}

func main() {
	flag.Parse()
	go func() {
		runUI()
	}()
	app.Main()
}

func NewApp() *App {
	wallet, err := NewWallet(*pubAddr, *host)
	if err != nil {
		log.Fatal(err)
	}
	qr, err := qrcode.New(strings.ToUpper(*pubAddr), qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}
	qr.BackgroundColor = rgb(0xe8f5e9)
	return &App{
		transList: &layout.List{
			Axis: layout.Vertical,
		},
		addrQR: paint.NewImageOp(qr.Image(256)),
		wallet: wallet,
	}
}

func runUI() {
	a := NewApp()
	w := app.NewWindow()
	gtx := &layout.Context{
		Queue: w.Queue(),
	}
	for {
		select {
		case e := <-a.wallet.events:
			switch e := e.(type) {
			case TransactionEvent:
				a.trans = append(a.trans, e.Trans)
				w.Invalidate()
			}
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				if err := e.Err; err != nil {
					log.Fatal(err)
				}
				return
			case system.FrameEvent:
				gtx.Reset(e.Config, e.Size)
				a.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func (a *App) Layout(gtx *layout.Context) {
	xlayout.Format(gtx, "stack(southeast, r(max(_)), r(inset(16dp, _)))",
		func() {
			a.layoutMain(gtx)
		},
		func() {
			for a.qrBtn.Clicked(gtx) {
				a.showQR = !a.showQR
			}
			theme.IconButton(qrIcn).Layout(gtx, &a.qrBtn)
		},
	)
	if a.showQR {
		xlayout.Format(gtx, "center(inset(16dp, _))",
			func() {
				Corners(unit.Dp(10)).Layout(gtx, func() {
					xlayout.Format(gtx, "stack(center, r(_), e(min(inset(16dp, south(_)))))",
						func() {
							sz := gtx.Constraints.Width.Constrain(gtx.Px(unit.Dp(500)))
							a.addrQR.Add(gtx.Ops)
							paint.PaintOp{
								Rect: f32.Rectangle{
									Max: f32.Point{
										X: float32(sz), Y: float32(sz),
									},
								},
							}.Add(gtx.Ops)
							gtx.Dimensions.Size = image.Point{X: sz, Y: sz}
						},
						func() {
							theme.Body1(*pubAddr).Layout(gtx)
						},
					)
				})
			},
		)
	}
}

func (a *App) layoutMain(gtx *layout.Context) {
	const f = `
		vflex(
			r(vcap(200dp, vmax(
				stack(center, 
					e(max(_))
					r(hflex(baseline, r(_), r(_)))
				)
			))),
			r(_)
		)
	`
	xlayout.Format(gtx, f,
		func() {
			fill(rgb(0x1b5e20)).Layout(gtx)
		},
		func() {
			bal := fmt.Sprintf("%d", a.wallet.Balance())
			amt := theme.H2(bal)
			amt.Color = rgb(0xffffff)
			if len(a.trans) > 0 {
				last := a.trans[len(a.trans)-1]
				const duration = .2
				if dt := time.Now().Sub(last.added).Seconds(); dt < duration {
					dt /= duration
					dt = dt * dt * dt
					const scale = 0.1
					amt.Font.Size.V *= float32(1 + scale - dt*scale)
					op.InvalidateOp{}.Add(gtx.Ops)
				}
			}
			amt.Layout(gtx)
		},
		func() {
			sat := theme.H6(" sat")
			sat.Color = rgb(0xffffff)
			sat.Layout(gtx)
		},
		func() {
			a.layoutTrans(gtx)
		},
	)
}

func (a *App) layoutTrans(gtx *layout.Context) {
	now := time.Now()
	a.transList.Layout(gtx, len(a.trans), func(i int) {
		// Invert list
		i = len(a.trans) - 1 - i
		t := a.trans[i]
		a := 1.0
		const duration = 0.5
		if dt := now.Sub(t.added).Seconds(); dt < duration {
			op.InvalidateOp{}.Add(gtx.Ops)
			a = dt / duration
			a *= a
		}
		const f = `
		inset(16dp,
			vflex(
				r(inset(0dp0dp4dp0dp, hflex(baseline,
					r(_),
					r(hmax(east(hflex(baseline, r(_), r(_)))))
				))),
				r(_)
			)
		)`
		xlayout.Format(gtx, f,
			func() {
				tim := theme.Body1(t.FormatTime())
				tim.Color = alpha(a, tim.Color)
				tim.Layout(gtx)
			},
			func() {
				amount := theme.H5(t.FormatAmount())
				amount.Color = rgb(0x003300)
				amount.Color = alpha(a, amount.Color)
				amount.Alignment = text.End
				amount.Font.Variant = "Mono"
				amount.Font.Weight = text.Bold
				amount.Layout(gtx)
			},
			func() {
				sat := theme.Body1(" sat")
				sat.Color = alpha(a, sat.Color)
				sat.Layout(gtx)
			},
			func() {
				l := theme.Body2(t.Hash)
				l.Color = theme.Color.Hint
				l.Color = alpha(a, l.Color)
				l.Layout(gtx)
			},
		)
	})
}

func alpha(a float64, col color.RGBA) color.RGBA {
	col.A = byte(float64(col.A) * a)
	col.R = byte(float64(col.R) * a)
	col.G = byte(float64(col.G) * a)
	col.B = byte(float64(col.B) * a)
	return col
}

type Corners unit.Value

func (c Corners) Layout(gtx *layout.Context, w layout.Widget) {
	var macro op.MacroOp
	macro.Record(gtx.Ops)
	w()
	macro.Stop()
	sz := gtx.Dimensions.Size
	rr := float32(gtx.Px(unit.Value(c)))
	var stack op.StackOp
	stack.Push(gtx.Ops)
	rrect(gtx.Ops, float32(sz.X), float32(sz.Y), rr, rr, rr, rr)
	macro.Add(gtx.Ops)
	stack.Pop()
}

func rgb(c uint32) color.RGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func (f fill) Layout(gtx *layout.Context) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: color.RGBA(f)}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

// https://pomax.github.io/bezierinfo/#circles_cubic.
func rrect(ops *op.Ops, width, height, se, sw, nw, ne float32) {
	w, h := float32(width), float32(height)
	const c = 0.55228475 // 4*(sqrt(2)-1)/3
	var b paint.Path
	b.Begin(ops)
	b.Move(f32.Point{X: w, Y: h - se})
	b.Cube(f32.Point{X: 0, Y: se * c}, f32.Point{X: -se + se*c, Y: se}, f32.Point{X: -se, Y: se}) // SE
	b.Line(f32.Point{X: sw - w + se, Y: 0})
	b.Cube(f32.Point{X: -sw * c, Y: 0}, f32.Point{X: -sw, Y: -sw + sw*c}, f32.Point{X: -sw, Y: -sw}) // SW
	b.Line(f32.Point{X: 0, Y: nw - h + sw})
	b.Cube(f32.Point{X: 0, Y: -nw * c}, f32.Point{X: nw - nw*c, Y: -nw}, f32.Point{X: nw, Y: -nw}) // NW
	b.Line(f32.Point{X: w - ne - nw, Y: 0})
	b.Cube(f32.Point{X: ne * c, Y: 0}, f32.Point{X: ne, Y: ne - ne*c}, f32.Point{X: ne, Y: ne}) // NE
	b.End().Add(ops)
}
