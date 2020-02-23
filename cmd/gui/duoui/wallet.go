package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var (
	passPhrase        string
	confirmPassPhrase string
	passEditor        = &controller.Editor{
		SingleLine: true,
		//Submit:     true,
	}
	confirmPassEditor = &controller.Editor{
		SingleLine: true,
		//Submit:     true,
	}

	encryption         = new(controller.CheckBox)
	seed               = new(controller.CheckBox)
	testnet            = new(controller.CheckBox)
	buttonCreateWallet = new(controller.Button)
)

func (ui *DuoUI) DuoUIloaderCreateWallet() {
	cs := ui.ly.Context.Constraints
	theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Center.Layout(ui.ly.Context, func() {
		controllers := []func(){
			func() {
				bal := ui.ly.Theme.H5("Enter the private passphrase for your new wallet:")
				bal.Font.Typeface = ui.ly.Theme.Font.Primary
				bal.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				bal.Layout(ui.ly.Context)
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
					e := ui.ly.Theme.DuoUIeditor("Enter Passpharse")
					e.Font.Typeface = ui.ly.Theme.Font.Primary
					e.Font.Style = text.Regular
					e.Layout(ui.ly.Context, passEditor)
					for _, e := range passEditor.Events(ui.ly.Context) {
						switch e.(type) {
						case controller.ChangeEvent:
							passPhrase = passEditor.Text()
						}
					}
				})
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
					e := ui.ly.Theme.DuoUIeditor("Repeat Passpharse")
					e.Font.Typeface = ui.ly.Theme.Font.Primary
					e.Font.Style = text.Regular
					e.Layout(ui.ly.Context, confirmPassEditor)
					for _, e := range confirmPassEditor.Events(ui.ly.Context) {
						switch e.(type) {
						case controller.ChangeEvent:
							confirmPassPhrase = confirmPassEditor.Text()
						}
					}
				})
			},
			func() {
				encryptionCheckBox := ui.ly.Theme.DuoUIcheckBox("Do you want to add an additional layer of encryption for public data?", "ff303030", "ff303030")
				encryptionCheckBox.Font.Typeface = ui.ly.Theme.Font.Primary
				encryptionCheckBox.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				encryptionCheckBox.Layout(ui.ly.Context, encryption)
			},
			func() {
				seedCheckBox := ui.ly.Theme.DuoUIcheckBox("Do you have an existing wallet seed you want to use?", "ff303030", "ff303030")
				seedCheckBox.Font.Typeface = ui.ly.Theme.Font.Primary
				seedCheckBox.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				seedCheckBox.Layout(ui.ly.Context, seed)
			},
			func() {
				testnetCheckBox := ui.ly.Theme.DuoUIcheckBox("Use testnet?", "ff303030", "ff303030")
				testnetCheckBox.Font.Typeface = ui.ly.Theme.Font.Primary
				testnetCheckBox.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				testnetCheckBox.Layout(ui.ly.Context, testnet)
			},
			func() {
				var createWalletbuttonComp theme.DuoUIbutton
				createWalletbuttonComp = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "CREATE WALLET", "ff303030", "ffcfcfcf", "", "ff303030", 16, 0, 125, 32, 4, 4)
				for buttonCreateWallet.Clicked(ui.ly.Context) {
					if passPhrase != "" && passPhrase == confirmPassPhrase {
						if testnet.Checked(ui.ly.Context) {
							ui.rc.UseTestnet()
						}
						ui.rc.CreateWallet(passPhrase, "", "", "")
						if testnet.Checked(ui.ly.Context) {
							interrupt.RequestRestart()
						}
						log.INFO("WOIKOS!")
					}
					log.INFO("confirmPassPhrase: ", confirmPassPhrase)
					log.INFO("passPhrase: ", passPhrase)
					log.INFO("posleWOIKOS!")
				}
				createWalletbuttonComp.Layout(ui.ly.Context, buttonCreateWallet)
			},
		}
		list.Layout(ui.ly.Context, len(controllers), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(ui.ly.Context, controllers[i])
		})
	})
}
