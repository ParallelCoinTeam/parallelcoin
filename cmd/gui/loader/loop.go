package loader

import (
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/pkg/gio/io/system"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/gio/widget/material"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"image/color"
)

var (
	passPhrase        = ""
	confirmPassPhrase = ""
	passEditor        = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	confirmPassEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	encryption         = new(widget.CheckBox)
	seed               = new(widget.CheckBox)
	buttonCreateWallet = new(widget.Button)
	list               = &layout.List{
		Axis: layout.Vertical,
	}
	ln = layout.UniformInset(unit.Dp(1))
	in = layout.UniformInset(unit.Dp(8))
)

type DuoUIload struct {
	//Boot *Boot
	//rc   *RcVar
	cx   *conte.Xt
	ww   *app.Window
	gc   *layout.Context
	th   *material.Theme
	cs   *layout.Constraints
	ico  *models.DuoUIicons
	comp *models.DuoUIcomponents
}

func DuoUIloaderLoop(cx *conte.Xt) error {
	//fonts.Register()
	ldr := &DuoUIload{
		cx: cx,
		ww: app.NewWindow(),
	}

	//duo.ww.Width =  unit.Dp(800)
	//	duo.ww.Height=  unit.Dp(600)
	//		duo.ww.Title =  "ParallelCoin - True Story"

	ldr.gc = layout.NewContext(ldr.ww.Queue())
	ldr.cs = &ldr.gc.Constraints

	// Layouts
	view := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
	}
	header := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween},
	}
	logo := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
	}
	body := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Horizontal},
	}

	ldr.comp = &models.DuoUIcomponents{
		View:   view,
		Header: header,
		Logo:   logo,
		Body:   body,
	}

	// Icons

	var err error
	ics := &models.DuoUIicons{}

	ics.Logo, err = material.NewIcon(ico.ParallelCoin)
	if err != nil {
		log.FATAL(err)
	}
	ldr.ico = ics

	ldr.th = material.NewTheme()

	ldr.gc = layout.NewContext(ldr.ww.Queue())
	ldr.cs = &ldr.gc.Constraints
	for {
		e := <-ldr.ww.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ldr.gc.Reset(e.Config, e.Size)
			DuoUIloaderGrid(ldr)
			e.Frame(ldr.gc.Ops)
		}
	}
}

// START OMIT
func DuoUIloaderGrid(ldr *DuoUIload) {
	// START View <<<
	ldr.comp.View.Layout.Layout(ldr.gc, DuoUIloaderBody(ldr))
	// END View >>>
}

// END OMIT
// START OMIT
func DuoUIloaderBody(ldr *DuoUIload) layout.FlexChild {
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
