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
	Debug("GetNewReceivingAddress")
	var addr util.Address
	var err error
	if addr, err = wg.WalletClient.GetNewAddress("default"); !Check(err) {
		Debug("getting new receiving address", addr.EncodeAddress(),
			"previous:", wg.State.currentReceivingAddress.String.Load())
		// save to addressbook
		var ae AddressEntry
		ae.Address = addr.EncodeAddress()
		var amt float64
		if amt, err = strconv.ParseFloat(
			wg.inputs["receiveAmount"].GetText(),
			64,
		); !Check(err) {
			if ae.Amount, err = util.NewAmount(amt); Check(err) {
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
		Debugs(wg.State.receiveAddresses)
		wg.State.SetReceivingAddress(addr)
		filename := filepath.Join(wg.cx.DataDir, "state.json")
		if err := wg.State.Save(filename, wg.cx.Config.WalletPass); Check(err) {
		}
		wg.invalidate <- struct{}{}
	}
}

func (wg *WalletGUI) GetNewReceivingQRCode(qrText string) {
	wg.currentReceiveRegenerate.Store(false)
	var qrc image.Image
	Debug("generating QR code")
	var err error
	if qrc, err = qrcode.Encode(qrText, 0, qrcode.ECLevelL, 4); !Check(err) {
		iop := paint.NewImageOp(qrc)
		wg.currentReceiveQRCode = &iop
		wg.currentReceiveQR = wg.ButtonLayout(
			wg.currentReceiveCopyClickable.SetClick(
				func() {
					Debug("clicked qr code copy clicker")
					if err := clipboard.WriteAll(qrText); Check(err) {
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
