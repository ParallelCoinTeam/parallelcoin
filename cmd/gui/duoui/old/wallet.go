package duoui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/stalker-loki/pod/pkg/gui/gel"
	"github.com/stalker-loki/pod/pkg/gui/gelook"
	"github.com/stalker-loki/pod/pkg/util/interrupt"
)

var (
	passPhrase        string
	confirmPassPhrase string
	passEditor        = &gel.Editor{
		SingleLine: true,
		// Submit:     true,
	}
	confirmPassEditor = &gel.Editor{
		SingleLine: true,
		// Submit:     true,
	}
	listWallet = &layout.List{
		Axis: layout.Vertical,
	}
	encryption         = new(gel.CheckBox)
	seed               = new(gel.CheckBox)
	testnet            = new(gel.CheckBox)
	buttonCreateWallet = new(gel.Button)
)

func (ui *DuoUI) DuoUIloaderCreateWallet() {
	cs := ui.ly.Context.Constraints
	th := ui.ly.Theme
	ctx := ui.ly.Context
	gelook.DuoUIdrawRectangle(ctx, cs.Width.Max, cs.Height.Max,
		th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Center.Layout(ctx, func() {
		controllers := []func(){
			func() {
				bal := th.H5("Enter the private passphrase for your new" +
					" wallet:")
				bal.Font.Typeface = th.Fonts["Primary"]
				bal.Color = th.Colors["Dark"]
				bal.Layout(ctx)
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ctx, func() {
					e := th.DuoUIeditor("Enter Passphrase", "Dark", "Light", 32)
					e.Font.Typeface = th.Fonts["Primary"]
					e.Font.Style = text.Regular
					e.Layout(ctx, passEditor)
					for _, e := range passEditor.Events(ctx) {
						switch e.(type) {
						case gel.ChangeEvent:
							passPhrase = passEditor.Text()
						}
					}
				})
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ctx, func() {
					e := th.DuoUIeditor("Repeat Passphrase", "Dark", "Light", 32)
					e.Font.Typeface = th.Fonts["Primary"]
					e.Font.Style = text.Regular
					e.Layout(ctx, confirmPassEditor)
					for _, e := range confirmPassEditor.Events(ctx) {
						switch e.(type) {
						case gel.ChangeEvent:
							confirmPassPhrase = confirmPassEditor.Text()
						}
					}
				})
			},
			func() {
				encryptionCheckBox := th.DuoUIcheckBox(
					"Do you want to add an additional layer of encryption"+
						" for public data?", th.Colors["Dark"],
					th.Colors["Dark"])
				encryptionCheckBox.Font.Typeface = th.Fonts["Primary"]
				encryptionCheckBox.Color = gelook.HexARGB(th.Colors["Dark"])
				encryptionCheckBox.Layout(ctx, encryption)
			},
			func() {
				// TODO: needs input box for seed
				seedCheckBox := th.DuoUIcheckBox(
					"Do you have an existing wallet seed you want to use?",
					th.Colors["Dark"], th.Colors["Dark"])
				seedCheckBox.Font.Typeface = th.Fonts["Primary"]
				seedCheckBox.Color = gelook.HexARGB(th.Colors["Dark"])
				seedCheckBox.Layout(ctx, seed)
			},
			func() {
				testnetCheckBox := th.DuoUIcheckBox(
					"Use testnet?", th.Colors["Dark"], th.Colors["Dark"])
				testnetCheckBox.Font.Typeface = th.Fonts["Primary"]
				testnetCheckBox.Color = gelook.HexARGB(th.Colors["Dark"])
				testnetCheckBox.Layout(ctx, testnet)
			},
			func() {
				var createWalletbuttonComp gelook.DuoUIbutton
				createWalletbuttonComp = th.DuoUIbutton(th.
					Fonts["Secondary"], "CREATE WALLET", th.Colors["Dark"],
					th.Colors["Light"], th.Colors["Light"],
					th.Colors["Dark"], "", th.Colors["Dark"], 16, 0,
					125, 32, 4, 4, 4, 4)
				for buttonCreateWallet.Clicked(ctx) {
					if passPhrase != "" && passPhrase == confirmPassPhrase {
						if testnet.Checked(ctx) {
							ui.rc.UseTestnet()
						}
						ui.rc.CreateWallet(passPhrase, "", "", "")
						if testnet.Checked(ctx) {
							interrupt.RequestRestart()
						}
					}
				}
				createWalletbuttonComp.Layout(ctx, buttonCreateWallet)
			},
		}
		listWallet.Layout(ctx, len(controllers), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(ctx, controllers[i])
		})
	})
}
