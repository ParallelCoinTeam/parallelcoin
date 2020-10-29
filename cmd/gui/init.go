package gui

import (
	l "gioui.org/layout"
)

// InitWallet renders a wallet initialization input form
func (wg *WalletGUI) InitWallet() func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		x := gtx.Constraints.Max.X
		wg.Size = &x
		// TODO: put the root stack in here
		return wg.Flex().Rigid(
			wg.loaderCreateWallet(),
		).Fn(gtx)
	}
}

func (wg *WalletGUI) loaderCreateWallet() l.Widget {
	createWalletLayoutList := []l.Widget{
		wg.Inset(0.0, wg.Fill("DocText", wg.Inset(0.5, wg.H6("Enter the private passphrase for your new wallet:").Color("DocBg").Fn).Fn).Fn).Fn,

		wg.Inset(0.25,
			wg.Border().Embed(
				wg.Inset(0.25,
					wg.SimpleInput(wg.editors["passEditor"].
						SetText("Enter Passphrase")).Fn,
				).Fn,
			).Fn,
		).Fn,

		wg.Inset(0.25,
			wg.Border().Embed(
				wg.Inset(0.25,
					wg.SimpleInput(wg.editors["confirmPassEditor"].
						SetText("Repeat Passphrase")).Fn,
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
			wg.clickables["createwallet"].SetClick(func() {
				Info("clicked customised button")
			})).
			CornerRadius(3).
			Background("Secondary").
			Color("Dark").
			Font("bariol bold").
			TextScale(2).
			Text("CREATE WALLET").
			Inset(1.5).
			Fn,
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return createWalletLayoutList[index](gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.lists["createwallet"].Vertical().Length(len(createWalletLayoutList)).ListElement(le).Fn(gtx)
	}
}
