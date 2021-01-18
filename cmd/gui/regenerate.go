package gui

import (
	"fmt"
	"gioui.org/op/paint"
	"github.com/atotto/clipboard"
	"github.com/p9c/pod/pkg/coding/qrcode"
	"github.com/p9c/pod/pkg/util"
	"image"
	"path/filepath"
	"strconv"
	"time"
)

func (wg *WalletGUI) GetNewReceivingAddress() {
	var addr util.Address
	var err error
	if addr, err = wg.WalletClient.GetNewAddress("default"); !Check(err) {
		// Debug("getting new address new receiving address", addr.EncodeAddress(),
		// 	"as prior was empty", wg.State.currentReceivingAddress.String.Load())
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
		ae.Message = wg.inputs["receiveMessage"].GetText()
		ae.Created = time.Now()
		wg.State.receiveAddresses = append(wg.State.receiveAddresses, ae)
		Debugs(wg.State.receiveAddresses)
		// TODO: update the receive addressbook widget
		
		wg.State.SetReceivingAddress(addr)
		wg.State.isAddress.Store(true)
		filename := filepath.Join(wg.cx.DataDir, "state.json")
		if err := wg.State.Save(filename, wg.cx.Config.WalletPass); Check(err) {
		}
		wg.invalidate <- struct{}{}
	}
	
}

func (wg *WalletGUI) GetNewReceivingQRCode() {
	wg.currentReceiveRegenerate.Store(false)
	var qrc image.Image
	Debug("generating QR code")
	var err error
	qrText := fmt.Sprintf(
		"parallelcoin:%s?amount=%s&message=%s",
		wg.State.currentReceivingAddress.Load().EncodeAddress(),
		wg.inputs["receiveAmount"].GetText(),
		wg.inputs["receiveMessage"].GetText(),
	)
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