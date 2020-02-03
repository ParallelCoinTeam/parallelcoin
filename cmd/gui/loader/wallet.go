package loader

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"github.com/p9c/pod/pkg/log"
	"image/color"
)

var (
	createWalletbutton = new(widget.Button)

	consoleInputField = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	logMessages []log.Entry
	logChan     = make(chan log.Entry, 111)
)

func init() {
	log.L.LogChan = logChan
	log.L.SetLevel("Info", false)
	go func() {
		for {
			select {
			case n := <-log.L.LogChan:
				logMessages = append(logMessages, n)
			}
		}
	}()
}

// START OMIT
func DuoUIloaderCreateWallet(duo *models.DuoUI, cx *conte.Xt) {
	//const buflen = 9
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(0.5, func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ff303030", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
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
					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ff303030", [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
					ln.Layout(duo.DuoUIcontext, func() {
						helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "fff4f4f4", [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
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

					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ff303030", [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
					ln.Layout(duo.DuoUIcontext, func() {
						helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "fff4f4f4", [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})
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
					var createWalletbuttonComp parallel.DuoUIbutton
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
		layout.Flexed(0.5, func() {
			//const n = 1e6
			logOutputList.Layout(duo.DuoUIcontext, len(logMessages), func(i int) {
				t := logMessages[i]
				cs := duo.DuoUIcontext.Constraints
				col := "ff3030cf"

				if t.Level == "TRC" {
					col = "ff3030cf"
				}
				if t.Level == "DBG" {
					col = "ffcfcf30"
				}
				if t.Level == "INF" {
					col = "ff30cf30"
				}
				if t.Level == "WRN" {
					col = "ffcfcf30"
				}
				if t.Level == "Error" {
					col = "ffcf8030"
				}
				if t.Level == "FTL" {
					col = "ffcf3030"
				}

				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, col, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

				logText := duo.DuoUItheme.H6(fmt.Sprint(i) + "->" + fmt.Sprint(t.Text))
				logText.Layout(duo.DuoUIcontext)
			})
		}),
	)
}

// END OMIT

func ttDuoUIloaderCreateWallet(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	//const buflen = 1000
	//log.L.LogChan = make(chan log.Entry)
	//
	//logMessages := make([]string, buflen)
	//var cursor int
	//quit := make(chan struct{})
	//
	//go func() {
	//	for {
	//		select {
	//		case n := <-log.L.LogChan:
	//			log.INFO(n)
	//		logMessages[cursor] = fmt.Sprint(n)
	//		cursor ++
	//		if cursor > buflen {
	//			cursor = 0
	//		}
	//		case <-quit:
	//		}
	//	}
	//}()
	//logMessages := []string{"asasas", "Eeee"}

	return duo.DuoUIcontext, func() {
		layout.Flex{}.Layout(duo.DuoUIcontext,
			layout.Flexed(1, func() {

				layout.Flex{}.Layout(duo.DuoUIcontext,
					layout.Flexed(1, func() {
						duo.DuoUIcomponents.Console.Inset.Layout(duo.DuoUIcontext, func() {
							layout.Flex{
								Axis:    layout.Vertical,
								Spacing: layout.SpaceAround,
							}.Layout(duo.DuoUIcontext,
								layout.Rigid(func() {
									cs := duo.DuoUIcontext.Constraints
									helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ff303030", [4]float32{9, 9, 9, 9}, [4]float32{0, 0, 0, 0})

									///
									//logOutputList.Layout(duo.DuoUIcontext, len(logMessages), func(i int) {
									//	t := logMessages[i]
									//	layout.Flex{
									//		Alignment: layout.End,
									//	}.Layout(duo.DuoUIcontext,
									//		layout.Rigid(func() {
									//			sat := duo.DuoUItheme.Body1(t)
									//			sat.Font.Size = unit.Dp(16)
									//			sat.Layout(duo.DuoUIcontext)
									//		}),
									//	)
									//})
									const n = 1e6
									logOutputList.Layout(duo.DuoUIcontext, n, func(i int) {
										txt := fmt.Sprintf("List element #%d", i)

										duo.DuoUItheme.H3(txt).Layout(duo.DuoUIcontext)
									})
								}))
							// Overview >>>
						})
					}),
				)
			}))
	}
}
