package loader

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/log"
	"image/color"
)

// START OMIT
func DuoUIloaderCreateWallet(ldr *DuoUIload) layout.FlexChild {
	return ldr.comp.View.Layout.Flex(ldr.gc, 1, func() {
		helpers.DuoUIdrawRect(ldr.gc, ldr.cs.Width.Max, ldr.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 0, 0, 0, 0)
		// START View <<<
		widgets := []func(){
			func() {
				bal := ldr.th.H3("Enter the private passphrase for your new wallet:")

				bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
				bal.Layout(ldr.gc)

				helpers.DuoUIdrawRect(ldr.gc, ldr.cs.Width.Max, ldr.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
				ln.Layout(ldr.gc, func() {
					helpers.DuoUIdrawRect(ldr.gc, ldr.cs.Width.Max, ldr.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
					in.Layout(ldr.gc, func() {
						e := ldr.th.Editor("Enter Passpharse")
						e.Font.Style = text.Italic
						e.Font.Size = unit.Dp(24)
						e.Layout(ldr.gc, passEditor)
						for _, e := range passEditor.Events(ldr.gc) {
							if e, ok := e.(widget.SubmitEvent); ok {
								passPhrase = e.Text
								passEditor.SetText("")
							}
						}
					})
				})
			},
			func() {

				helpers.DuoUIdrawRect(ldr.gc, ldr.cs.Width.Max, ldr.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
				ln.Layout(ldr.gc, func() {
					helpers.DuoUIdrawRect(ldr.gc, ldr.cs.Width.Max, ldr.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
					in.Layout(ldr.gc, func() {
						e := ldr.th.Editor("Repeat Passpharse")
						e.Font.Style = text.Italic
						e.Font.Size = unit.Dp(24)
						e.Layout(ldr.gc, confirmPassEditor)
						for _, e := range confirmPassEditor.Events(ldr.gc) {
							if e, ok := e.(widget.SubmitEvent); ok {
								confirmPassPhrase = e.Text
								confirmPassEditor.SetText("")
							}
						}
					})
				})
			},
			func() {
				encryptionCheckBox := ldr.th.CheckBox("Do you want to add an additional layer of encryption for public data?")
				encryptionCheckBox.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
				encryptionCheckBox.Layout(ldr.gc, encryption)
			},
			func() {
				seedCheckBox := ldr.th.CheckBox("Do you have an existing wallet seed you want to use?")
				seedCheckBox.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
				seedCheckBox.Layout(ldr.gc, seed)
			},
			func() {

				for buttonCreateWallet.Clicked(ldr.gc) {
					if passPhrase != "" && confirmPassPhrase == confirmPassPhrase {
						CreateWallet(ldr, passPhrase, "", "", "")
						log.INFO("WOIKOS!")
					}

				}
				ldr.th.Button("Click me!").Layout(ldr.gc, buttonCreateWallet)

			},
		}
		list.Layout(ldr.gc, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(16)).Layout(ldr.gc, widgets[i])
		})
	})
}

// END OMIT
