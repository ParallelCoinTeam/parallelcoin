package gui

import (
	l "gioui.org/layout"
	icons "github.com/p9c/pod/pkg/gui/ico/svg"
)

// InitWallet renders a wallet initialization input form
func (wg *WalletGUI) InitWallet() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		x := gtx.Constraints.Max.X
		wg.Size = &x
		// TODO: put the root stack in here
		return wg.Flex().Rigid(
			wg.Fill("DocText", wg.Inset(1,
				wg.Fill("DocBg", wg.Inset(1,
					wg.loaderCreateWallet(),
				).Fn).Fn).Fn).Fn).Fn(gtx)
	}
}

func (wg *WalletGUI) loaderCreateWallet() l.Widget {
	createWalletLayoutList := []l.Widget{

		wg.Icon().Scale(5).Color("DocText").Src(icons.ParallelCoinRound).Fn,

		wg.Inset(0.0, wg.Fill("DocBg", wg.Inset(0.5, wg.H6("Enter the private passphrase for your new wallet:").Color("DocText").Fn).Fn).Fn).Fn,
		wg.Inset(0.25,
			wg.Border().Embed(
				wg.Inset(0.25,
					wg.passwords["passEditor"].Fn,
				).Fn,
			).Fn,
		).Fn,
		wg.Inset(0.25,
			wg.Border().Embed(
				wg.Inset(0.25,
					wg.passwords["confirmPassEditor"].Fn,
				).Fn,
			).Fn,
		).Fn,
		wg.CheckBox(wg.bools["encryption"].SetOnChange(func(b bool) {
			Debug("change state to", b)
		})).
			IconColor("Primary").
			TextColor("DocText").
			// IconScale(0.1).
			Text("Do you want to add an additional layer of encryption for public data?").
			Fn,
		wg.CheckBox(wg.bools["seed"].SetOnChange(func(b bool) {
			Debug("change state to", b)
		})).
			IconColor("Primary").
			TextColor("DocText").
			// IconScale(0.1).
			Text("Do you have an existing wallet seed you want to use?").
			Fn,

		wg.CheckBox(wg.bools["testnet"].SetOnChange(func(b bool) {
			Debug("change state to", b)
		})).
			IconColor("Primary").
			TextColor("DocText").
			// IconScale(0.1).
			Text("Use testnet?").
			Fn,
		wg.Button(
			wg.clickables["createWallet"].SetClick(func() {
				Info("clicked customised button")
			})).
			CornerRadius(0).
			Background("Primary").
			Color("Dark").
			Font("bariol bold").
			TextScale(1).
			Text("CREATE WALLET").
			Inset(0.5).
			Fn,
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return createWalletLayoutList[index](gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.lists["createWallet"].Vertical().Length(len(createWalletLayoutList)).ListElement(le).Fn(gtx)
	}
}
