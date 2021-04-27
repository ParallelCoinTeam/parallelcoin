package gui

import (
	"image"
	"path/filepath"
	"strconv"
	"time"

	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/btcaddr"

	"github.com/atotto/clipboard"

	"github.com/p9c/gio/op/paint"

	"github.com/p9c/matrjoska/pkg/qrcode"
)

func (wg *WalletGUI) GetNewReceivingAddress() {
	D.Ln("GetNewReceivingAddress")
	var addr btcaddr.Address
	var e error
	if addr, e = wg.WalletClient.GetNewAddress("default"); !E.Chk(e) {
		D.Ln(
			"getting new receiving address", addr.EncodeAddress(),
			"previous:", wg.State.currentReceivingAddress.String.Load(),
		)
		// save to addressbook
		var ae AddressEntry
		ae.Address = addr.EncodeAddress()
		var amount float64
		if amount, e = strconv.ParseFloat(
			wg.inputs["receiveAmount"].GetText(),
			64,
		); !E.Chk(e) {
			if ae.Amount, e = amt.NewAmount(amount); E.Chk(e) {
			}
		}
		msg := wg.inputs["receiveMessage"].GetText()
		if len(msg) > 64 {
			msg = msg[:64]
		}
		ae.Message = msg
		ae.Created = time.Now()
		if wg.State.IsReceivingAddress() {
			wg.State.receiveAddresses = append(wg.State.receiveAddresses, ae)
		} else {
			wg.State.receiveAddresses = []AddressEntry{ae}
			wg.State.isAddress.Store(true)
		}
		D.S(wg.State.receiveAddresses)
		wg.State.SetReceivingAddress(addr)
		filename := filepath.Join(wg.cx.Config.DataDir.V(), "state.json")
		if e = wg.State.Save(filename, wg.cx.Config.WalletPass.Bytes(), false); E.Chk(e) {
		}
		wg.invalidate <- struct{}{}
	}
}

func (wg *WalletGUI) GetNewReceivingQRCode(qrText string) {
	wg.currentReceiveRegenerate.Store(false)
	var qrc image.Image
	D.Ln("generating QR code")
	var e error
	if qrc, e = qrcode.Encode(qrText, 0, qrcode.ECLevelL, 4); !E.Chk(e) {
		iop := paint.NewImageOp(qrc)
		wg.currentReceiveQRCode = &iop
		wg.currentReceiveQR = wg.ButtonLayout(
			wg.currentReceiveCopyClickable.SetClick(
				func() {
					D.Ln("clicked qr code copy clicker")
					if e := clipboard.WriteAll(qrText); E.Chk(e) {
					}
				},
			),
		).
			Background("white").
			Embed(
				wg.Inset(
					0.125,
					wg.Image().Src(*wg.currentReceiveQRCode).Scale(1).Fn,
				).Fn,
			).Fn
	}
}
