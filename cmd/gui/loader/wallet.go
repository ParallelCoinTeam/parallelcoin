package loader

import (
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/log"
	"image/color"
)

var (
	createWalletbutton = new(widget.Button)
)

// START OMIT
func DuoUIloaderCreateWallet(duo *models.DuoUI, cx *conte.Xt) {

	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(1, func() {
			cs := duo.DuoUIcontext.Constraints

			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			// START View <<<
			widgets := []func(){
				func() {
					bale := duo.DuoUItheme.H3(passPhrase)
					bale.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					bale.Layout(duo.DuoUIcontext)
				},

				func() {
					balr := duo.DuoUItheme.H3(confirmPassPhrase)

					balr.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					balr.Layout(duo.DuoUIcontext)
				},
				func() {
					bal := duo.DuoUItheme.H3("Enter the private passphrase for your new wallet:")

					bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					bal.Layout(duo.DuoUIcontext)

					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
					ln.Layout(duo.DuoUIcontext, func() {
						helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("fff4f4f4"), [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
						in.Layout(duo.DuoUIcontext, func() {
							e := duo.DuoUItheme.DuoUIeditor("Enter Passpharse", "Enter Passpharse")
							e.Font.Style = text.Regular
							e.Font.Size = unit.Dp(24)
							e.Layout(duo.DuoUIcontext, passEditor)
							for _, e := range passEditor.Events(duo.DuoUIcontext) {
								if e, ok := e.(widget.SubmitEvent); ok {
									passPhrase = e.Text
									passEditor.SetText("")
								}
							}
						})
					})
				},
				func() {

					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
					ln.Layout(duo.DuoUIcontext, func() {
						helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("fff4f4f4"), [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
						in.Layout(duo.DuoUIcontext, func() {
							e := duo.DuoUItheme.DuoUIeditor("Repeat Passpharse", "Repeat Passpharse")
							e.Font.Style = text.Regular
							e.Font.Size = unit.Dp(24)
							e.Layout(duo.DuoUIcontext, confirmPassEditor)
							for _, e := range confirmPassEditor.Events(duo.DuoUIcontext) {
								if e, ok := e.(widget.SubmitEvent); ok {
									confirmPassPhrase = e.Text
									confirmPassEditor.SetText("")
								}
							}
						})
					})
				},
				func() {
					encryptionCheckBox := duo.DuoUItheme.DuoUIcheckBox("Do you want to add an additional layer of encryption for public data?")
					encryptionCheckBox.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					encryptionCheckBox.Layout(duo.DuoUIcontext, encryption)
				},
				func() {
					seedCheckBox := duo.DuoUItheme.DuoUIcheckBox("Do you have an existing wallet seed you want to use?")
					seedCheckBox.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					seedCheckBox.Layout(duo.DuoUIcontext, seed)
				},
				func() {
					var createWalletbuttonComp components.DuoUIbutton
					createWalletbuttonComp = duo.DuoUItheme.DuoUIbutton("Create wallet", "ff303030", "ffcfcfcf", "ff303030", 0, 125, 32, 4, 4, nil)
					for createWalletbutton.Clicked(duo.DuoUIcontext) {
							if passPhrase != "" && passPhrase == confirmPassPhrase {
								CreateWallet(cx, passPhrase, "", "", "")
								log.INFO("WOIKOS!")
							}
					}
					createWalletbuttonComp.Layout(duo.DuoUIcontext, createWalletbutton)
					//for buttonCreateWallet.Clicked(duo.DuoUIcontext) {
					//	if passPhrase != "" && passPhrase == confirmPassPhrase {
					//		//CreateWallet(ldr, passPhrase, "", "", "")
					//		log.INFO("WOIKOS!")
					//	}
					//
					//}
					////duo.DuoUItheme.DuoUIbutton("Create wallet").Layout(duo.DuoUIcontext, buttonCreateWallet)
					//duo.DuoUItheme.DuoUIbutton("Create wallet", "ff303030",  "ff989898", "ff303030", 0, 125, 32, 4, 4, nil)
				},
			}
			list.Layout(duo.DuoUIcontext, len(widgets), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, widgets[i])
			})
		}),
	)
}

// END OMIT
