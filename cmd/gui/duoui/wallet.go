package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/log"
	"image/color"
)

var (
	createWalletbutton = new(controller.Button)
	passEditor = &controller.Editor{
		SingleLine: true,
		Submit:     true,
	}
	confirmPassEditor = &controller.Editor{
		SingleLine: true,
		Submit:     true,
	}


	//logOutputList = &layout.List{
	//	Axis:        layout.Vertical,
	//	ScrollToEnd: true,
	//}

	encryption         = new(controller.CheckBox)
	seed               = new(controller.CheckBox)
	buttonCreateWallet = new(controller.Button)
)

//
//func init() {
//	log.L.LogChan = logChan
//	log.L.SetLevel("Info", false)
//	go func() {
//		for {
//			select {
//			case n := <-log.L.LogChan:
//				logMessages = append(logMessages, n)
//			}
//		}
//	}()
//}

func (ui *DuoUI) DuoUIloaderCreateWallet() {
	//const buflen = 9
	layout.Flex{}.Layout(ui.ly.Context,
		layout.Flexed(0.5, func() {
			cs := ui.ly.Context.Constraints
			theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			// START View <<<
			controllers := []func(){
				func() {
					bale := ui.ly.Theme.H3(ui.rc.PassPhrase)
					bale.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
					bale.Layout(ui.ly.Context)
				},
				func() {
					balr := ui.ly.Theme.H3(ui.rc.ConfirmPassPhrase)

					balr.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
					balr.Layout(ui.ly.Context)
				},
				func() {
					bal := ui.ly.Theme.H3("Enter the private passphrase for your new wallet:")
					bal.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
					bal.Layout(ui.ly.Context)
					theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Bg, [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
					layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
						theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, "fff4f4f4", [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
						layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
							e := ui.ly.Theme.DuoUIeditor("Enter Passpharse")
							e.Font.Style = text.Regular
							e.Layout(ui.ly.Context, passEditor)
							for _, e := range passEditor.Events(ui.ly.Context) {
								if e, ok := e.(controller.SubmitEvent); ok {
									ui.rc.PassPhrase = e.Text
									passEditor.SetText("")
								}
							}
						})
					})
				},
				func() {

					theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Bg, [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
					layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
						theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, "fff4f4f4", [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
						layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
							e := ui.ly.Theme.DuoUIeditor("Repeat Passpharse")
							e.Font.Style = text.Regular
							e.Layout(ui.ly.Context, confirmPassEditor)
							for _, e := range confirmPassEditor.Events(ui.ly.Context) {
								if e, ok := e.(controller.SubmitEvent); ok {
									ui.rc.ConfirmPassPhrase = e.Text
									confirmPassEditor.SetText("")
								}
							}
						})
					})
				},
				func() {
					encryptionCheckBox := ui.ly.Theme.DuoUIcheckBox("Do you want to add an additional layer of encryption for public data?")
					encryptionCheckBox.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					encryptionCheckBox.Layout(ui.ly.Context, encryption)
				},
				func() {
					seedCheckBox := ui.ly.Theme.DuoUIcheckBox("Do you have an existing wallet seed you want to use?")
					seedCheckBox.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					seedCheckBox.Layout(ui.ly.Context, seed)
				},
				func() {
					var createWalletbuttonComp theme.DuoUIbutton
					createWalletbuttonComp = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Primary,"Create wallet", "ff303030", "ffcfcfcf","", "ff303030", 0, 125, 32, 4, 4)
					for createWalletbutton.Clicked(ui.ly.Context) {
						if ui.rc.PassPhrase != "" && ui.rc.PassPhrase == ui.rc.ConfirmPassPhrase {
							//CreateWallet(ui.rc.Cx, ui.rc.PassPhrase, "", "", "")
							log.INFO("WOIKOS!")
						}
					}
					createWalletbuttonComp.Layout(ui.ly.Context, createWalletbutton)
					//for buttonCreateWallet.Clicked(ui.ly.Context) {
					//	if passPhrase != "" && passPhrase == confirmPassPhrase {
					//		//CreateWallet(ldr, passPhrase, "", "", "")
					//		log.INFO("WOIKOS!")
					//	}
					//
					//}
					////ui.ly.Theme.DuoUIbutton("Create wallet").Layout(ui.ly.Context, buttonCreateWallet)
					//ui.ly.Theme.DuoUIbutton("Create wallet", "ff303030",  "ff989898", "ff303030", 0, 125, 32, 4, 4, nil)
				},
			}
			list.Layout(ui.ly.Context, len(controllers), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, controllers[i])
			})
		}),
		layout.Flexed(0.5, func() {
			//const n = 1e6
			//logOutputList.Layout(ui.ly.Context, len(rc.Log.LogMessages), func(i int) {
			//	t := ui.rc.Log.LogMessages[i]
			//	cs := ui.ly.Context.Constraints
			//	col := "ff3030cf"
			//
			//	if t.Level == "TRC" {
			//		col = "ff3030cf"
			//	}
			//	if t.Level == "DBG" {
			//		col = "ffcfcf30"
			//	}
			//	if t.Level == "INF" {
			//		col = "ff30cf30"
			//	}
			//	if t.Level == "WRN" {
			//		col = "ffcfcf30"
			//	}
			//	if t.Level == "Error" {
			//		col = "ffcf8030"
			//	}
			//	if t.Level == "FTL" {
			//		col = "ffcf3030"
			//	}
			//
			//	view.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, col, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			//
			//	logText := ui.ly.Theme.H6(fmt.Sprint(i) + "->" + fmt.Sprint(t.Text))
			//	logText.Layout(ui.ly.Context)
			//})
		}),
	)
}
