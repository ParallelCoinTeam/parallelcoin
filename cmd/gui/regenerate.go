package gui

import (
	"image"
	"path/filepath"
	"strconv"
	"time"
	
	"gioui.org/op/paint"
	"github.com/atotto/clipboard"
	
	"github.com/p9c/pod/pkg/coding/qrcode"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) GetNewReceivingAddress() {
	dbg.Ln("GetNewReceivingAddress")
	var addr util.Address
	var e error
	if addr, e = wg.WalletClient.GetNewAddress("default"); !dbg.Chk(e) {
		dbg.Ln("getting new receiving address", addr.EncodeAddress(),
			"previous:", wg.State.currentReceivingAddress.String.Load())
		// save to addressbook
		var ae AddressEntry
		ae.Address = addr.EncodeAddress()
		var amt float64
		if amt, e = strconv.ParseFloat(
			wg.inputs["receiveAmount"].GetText(),
			64,
		); !dbg.Chk(e) {
			if ae.Amount, e = util.NewAmount(amt); dbg.Chk(e) {
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
		dbg.S(wg.State.receiveAddresses)
		wg.State.SetReceivingAddress(addr)
		filename := filepath.Join(wg.cx.DataDir, "state.json")
		if e := wg.State.Save(filename, wg.cx.Config.WalletPass); dbg.Chk(e) {
		}
		wg.invalidate <- struct{}{}
	}
}

func (wg *WalletGUI) GetNewReceivingQRCode(qrText string) {
	wg.currentReceiveRegenerate.Store(false)
	var qrc image.Image
	dbg.Ln("generating QR code")
	var e error
	if qrc, e = qrcode.Encode(qrText, 0, qrcode.ECLevelL, 4); !dbg.Chk(e) {
		iop := paint.NewImageOp(qrc)
		wg.currentReceiveQRCode = &iop
		wg.currentReceiveQR = wg.ButtonLayout(
			wg.currentReceiveCopyClickable.SetClick(
				func() {
					dbg.Ln("clicked qr code copy clicker")
					if e := clipboard.WriteAll(qrText); dbg.Chk(e) {
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
