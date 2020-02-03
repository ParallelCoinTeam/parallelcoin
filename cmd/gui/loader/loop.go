package loader

import (
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
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

//type DuoUIload struct {
//	//Boot *Boot
//	//rc   *RcVar
//	firstRun      bool
//	cx            *conte.Xt
//	loaderWindow  *app.Window
//	loaderContext *layout.Context
//	loaderTheme   *components.DuoUItheme
//	//ico  *models.DuoUIicons
//	comp *models.DuoUIcomponents
//}

//func DuoUIloaderLoop(firstRun bool, cx *conte.Xt) error {
	//fonts.Register()
	//ldr := &DuoUIload{
	//	firstRun:     firstRun,
	//	cx:           cx,
	//	loaderWindow: app.NewWindow(),
	//}

	//duo.ww.Width =  unit.Dp(800)
	//	duo.ww.Height=  unit.Dp(600)
	//		duo.ww.Title =  "ParallelCoin - True Story"

	//ldr.loaderContext = layout.NewContext(ldr.loaderWindow.Queue())
	//ldr.cs = &ldr.gc.Constraints

	// Layouts
	//view := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle},
	//}
	//header := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween},
	//}
	////intro := models.DuoUIcomponent{
	////	Layout: layout.Flex{Axis: layout.Vertical},
	////}
	//logo := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Vertical},
	//}
	////logMsgs := models.DuoUIcomponent{
	////	Layout: layout.Flex{Axis: layout.Vertical},
	////}
	//body := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Horizontal},
	//}
	////
	//ldr.comp = &models.DuoUIcomponents{
	//	View: view,
	//	//Intro:  intro,
	//	Header: header,
	//	Logo:   logo,
	//	//Log:    logMsgs,
	//	Body: body,
	//}

	// Icons

	//var err error
	//ics := &models.DuoUIicons{}

	//ics.Logo, err = components.NewDuoUIicon(ico.ParallelCoin)
	//if err != nil {
	//	log.FATAL(err)
	//}
	//ldr.ico = ics

	//ldr.loaderTheme = components.NewDuoUItheme()
	//
	//ldr.loaderContext = layout.NewContext(ldr.loaderWindow.Queue())
	////ldr.cs = &ldr.gc.Constraints
	//for {
	//	e := <-ldr.loaderWindow.Events()
	//	switch e := e.(type) {
	//	case system.DestroyEvent:
	//		return e.Err
	//	case system.FrameEvent:
	//		ldr.loaderContext.Reset(e.Config, e.Size)
	//		DuoUIloaderGrid(ldr)
	//		e.Frame(ldr.loaderContext.Ops)
	//	}
	//}
//}
//
//// START OMIT
//func DuoUIloaderGrid(duo models.DuoUI, rc rcd.RcVar) {
//	// START View <<<
//	if rc.IsFirstRun != false {
//		DuoUIloaderCreateWallet(ldr)
//	} else {
//		DuoUIloaderIntro(ldr)
//	}
//	// END View >>>
//}
//
//// END OMIT
