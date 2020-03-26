package duoui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/util/interrupt"
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
	gelook.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max,
		cs.Height.Max, ui.ly.Theme.Colors["Light"],
		[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Center.Layout(ui.ly.Context, func() {
		controllers := []func(){
			func() {
				bal := ui.ly.Theme.H5("Enter the private passphrase for your new wallet:")
				bal.Font.Typeface = ui.ly.Theme.Fonts["Primary"]
				bal.Color = ui.ly.Theme.Colors["Dark"]
				bal.Layout(ui.ly.Context)
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
					e := ui.ly.Theme.DuoUIeditor("Enter Passphrase")
					e.Font.Typeface = ui.ly.Theme.Fonts["Primary"]
					e.Font.Style = text.Regular
					e.Layout(ui.ly.Context, passEditor)
					for _, e := range passEditor.Events(ui.ly.Context) {
						switch e.(type) {
						case gel.ChangeEvent:
							passPhrase = passEditor.Text()
						}
					}
				})
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
					e := ui.ly.Theme.DuoUIeditor("Repeat Passphrase")
					e.Font.Typeface = ui.ly.Theme.Fonts["Primary"]
					e.Font.Style = text.Regular
					e.Layout(ui.ly.Context, confirmPassEditor)
					for _, e := range confirmPassEditor.Events(ui.ly.Context) {
						switch e.(type) {
						case gel.ChangeEvent:
							confirmPassPhrase = confirmPassEditor.Text()
						}
					}
				})
			},
			func() {
				encryptionCheckBox := ui.ly.Theme.DuoUIcheckBox(
					"Do you want to add an additional layer of encryption"+
						" for public data?", ui.ly.Theme.Colors["Dark"], ui.ly.Theme.Colors["Dark"])
				encryptionCheckBox.Font.Typeface = ui.ly.Theme.Fonts["Primary"]
				encryptionCheckBox.Color = gelook.HexARGB(ui.ly.Theme.Colors["Dark"])
				encryptionCheckBox.Layout(ui.ly.Context, encryption)
			},
			func() {
				seedCheckBox := ui.ly.Theme.DuoUIcheckBox(
					"Do you have an existing wallet seed you want to use?",
					ui.ly.Theme.Colors["Dark"], ui.ly.Theme.Colors["Dark"])
				seedCheckBox.Font.Typeface = ui.ly.Theme.Fonts["Primary"]
				seedCheckBox.Color = gelook.HexARGB(ui.ly.Theme.Colors["Dark"])
				seedCheckBox.Layout(ui.ly.Context, seed)
			},
			func() {
				testnetCheckBox := ui.ly.Theme.DuoUIcheckBox(
					"Use testnet?", ui.ly.Theme.Colors["Dark"], ui.ly.Theme.Colors["Dark"])
				testnetCheckBox.Font.Typeface = ui.ly.Theme.Fonts["Primary"]
				testnetCheckBox.Color = gelook.HexARGB(ui.ly.Theme.Colors["Dark"])
				testnetCheckBox.Layout(ui.ly.Context, testnet)
			},
			func() {
				var createWalletbuttonComp gelook.DuoUIbutton
				createWalletbuttonComp = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Fonts["Secondary"], "CREATE WALLET", ui.ly.Theme.Colors["Dark"], ui.ly.Theme.Colors["Light"], ui.ly.Theme.Colors["Light"], ui.ly.Theme.Colors["Dark"], "", ui.ly.Theme.Colors["Dark"], 16, 0, 125, 32, 4, 4, 4, 4)
				for buttonCreateWallet.Clicked(ui.ly.Context) {
					if passPhrase != "" && passPhrase == confirmPassPhrase {
						if testnet.Checked(ui.ly.Context) {
							ui.rc.UseTestnet()
						}
						ui.rc.CreateWallet(passPhrase, "", "", "")
						if testnet.Checked(ui.ly.Context) {
							interrupt.RequestRestart()
						}
					}
				}
				createWalletbuttonComp.Layout(ui.ly.Context, buttonCreateWallet)
			},
		}
		listWallet.Layout(ui.ly.Context, len(controllers), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(ui.ly.Context, controllers[i])
		})
	})
}
