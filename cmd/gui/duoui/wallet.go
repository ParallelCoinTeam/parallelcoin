package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/log"
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
	buttonCreateWallet = new(controller.Button)
)

func (ui *DuoUI) DuoUIloaderCreateWallet() {
	cs := ui.ly.Context.Constraints
	theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Center.Layout(ui.ly.Context, func() {
		controllers := []func(){
			func() {
				bale := ui.ly.Theme.Body1(passPhrase)
				bale.Font.Typeface = ui.ly.Theme.Font.Primary
				bale.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				bale.Layout(ui.ly.Context)
			},
			func() {
				balr := ui.ly.Theme.Body1(confirmPassPhrase)
				balr.Font.Typeface = ui.ly.Theme.Font.Primary
				balr.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				balr.Layout(ui.ly.Context)
			},
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
				encryptionCheckBox := ui.ly.Theme.DuoUIcheckBox("Do you want to add an additional layer of encryption for public data?")
				encryptionCheckBox.Font.Typeface = ui.ly.Theme.Font.Primary
				encryptionCheckBox.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				encryptionCheckBox.Layout(ui.ly.Context, encryption)
			},
			func() {
				seedCheckBox := ui.ly.Theme.DuoUIcheckBox("Do you have an existing wallet seed you want to use?")
				seedCheckBox.Font.Typeface = ui.ly.Theme.Font.Primary
				seedCheckBox.Color = theme.HexARGB(ui.ly.Theme.Color.Dark)
				seedCheckBox.Layout(ui.ly.Context, seed)
			},
			func() {
				var createWalletbuttonComp theme.DuoUIbutton
				createWalletbuttonComp = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "CREATE WALLET", "ff303030", "ffcfcfcf", "", "ff303030", 0, 125, 32, 4, 4)
				for buttonCreateWallet.Clicked(ui.ly.Context) {
					if passPhrase != "" && passPhrase == confirmPassPhrase {
						ui.rc.CreateWallet(passPhrase, "", "", "")
						log.INFO("WOIKOS!")
					}
					log.INFO("confirmPassPhrase: ", confirmPassPhrase)
					log.INFO("passPhrase: ", passPhrase)
					log.INFO("posleWOIKOS!")
				}
				createWalletbuttonComp.Layout(ui.ly.Context, buttonCreateWallet)

				//ui.ly.Theme.DuoUIbutton("Create wallet").Layout(ui.ly.Context, buttonCreateWallet)
				//ui.ly.Theme.DuoUIbutton("Create wallet", "ff303030", "ff989898", "ff303030", 0, 125, 32, 4, 4, nil)
				//linkButton = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Mono, "CREATE WALLET", "ffcfcfcf", "ff303030", "", "ffcfcfcf", 0, 60, 24, 0, 0)
				//
				//var blocksMenuItem theme.DuoUIbutton
				//blocksMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Primary, "Blocks: "+fmt.Sprint(ui.rc.Status.Node.BlockHeight), "ffcfcfcf", "", "", "", iconSize, 80, height, paddingVertical, 0)
				//for buttonBlocks.Clicked(ui.ly.Context) {
				//	ui.rc.ShowPage = "EXPLORER"
				//	//ui.rc.ShowToast = true
				//	//ui.toastAdd()
				//}
				//blocksMenuItem.Layout(ui.ly.Context, buttonBlocks)

			},
		}
		list.Layout(ui.ly.Context, len(controllers), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(ui.ly.Context, controllers[i])
		})
	})
}
