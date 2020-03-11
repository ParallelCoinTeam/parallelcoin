package component

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
)

// func DuoUIqrCode(pubAddr string) {
//	//qr, err := qrcode.New(strings.ToUpper(pubAddr), qrcode.Medium)
//	//if err != nil {
//	//	log.L.Fatal(err)
//	//}
//	//qr.BackgroundColor = rgb(0xe8f5e9)
//	//addrQR := paint.NewImageOp(qr.Image(256))
//	return
// }

// func NewQrCode(pubAddr string) *model.DuoUIqrCode {
//	//qr, err := qrcode.New(strings.ToUpper(pubAddr), qrcode.Medium)
//	//if err != nil {
//	//	log.L.Fatal(err)
//	//}
//	//log.L.Info(pubAddr)
//	//qr.BackgroundColor = theme.HexARGB("ff3030cf")
//	//return &model.DuoUIqrCode{
//	//	AddrQR:  paint.NewImageOp(qr.Image(256)),
//	//	PubAddr: pubAddr,
//	//}
// }

func DuoUIqrCode(rc *rcd.RcVar, gtx *layout.Context) func() {
	return func() {
		sz := gtx.Constraints.Width.Constrain(gtx.Px(unit.Dp(500)))
		rc.QrCode.AddrQR.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Max: f32.Point{
					X: float32(sz), Y: float32(sz),
				},
			},
		}.Add(gtx.Ops)
		gtx.Dimensions.Size = image.Point{X: sz, Y: sz}
	}
}

func QrDialog(rc *rcd.RcVar, gtx *layout.Context) func() {
	return func() {
		// clipboard.Set(t.Address)
		rc.Dialog.Show = true
		rc.Dialog = &model.DuoUIdialog{
			Show: true,
			// Ok:   rc.DuoSend(passPharse, address, amount),
			Close: func() {

			},
			CustomField: DuoUIqrCode(rc, gtx),
			Cancel:      func() { rc.Dialog.Show = false },
			Title:       "Are you sure?",
			Text:        "Confirm ParallelCoin send",
		}
	}
}
