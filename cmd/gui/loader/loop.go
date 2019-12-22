package loader

import (
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/pkg/gio/io/system"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/gio/widget/material"
	"github.com/p9c/pod/pkg/log"
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
	fr   bool
	cx   *conte.Xt
	ww   *app.Window
	gc   *layout.Context
	th   *material.Theme
	cs   *layout.Constraints
	ico  *models.DuoUIicons
	comp *models.DuoUIcomponents
}

func DuoUIloaderLoop(firstRun bool, cx *conte.Xt) error {
	//fonts.Register()
	ldr := &DuoUIload{
		fr: firstRun,
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
		Layout: layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle},
	}
	header := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween},
	}
	intro := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
	}
	logo := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
	}
	logMsgs := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
	}
	body := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Horizontal},
	}

	ldr.comp = &models.DuoUIcomponents{
		View:   view,
		Intro:  intro,
		Header: header,
		Logo:   logo,
		Log:    logMsgs,
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
	if ldr.fr != false {
		DuoUIloaderCreateWallet(ldr)
	} else {
		DuoUIloaderIntro(ldr)
	}
	// END View >>>
}

// END OMIT
