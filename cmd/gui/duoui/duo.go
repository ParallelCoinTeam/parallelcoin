package duoui

import (
	"github.com/p9c/gio-parallel/app"
	"github.com/p9c/gio-parallel/layout"
	"github.com/p9c/gio-parallel/unit"
	"github.com/p9c/gio-parallel/widget"
	"github.com/p9c/gio-parallel/widget/material"
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/fonts"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

func DuOuI(cx *conte.Xt) (duo *DuoUI) {
	duo = &DuoUI{
		cx: cx,
		rc: RcInit(),
		ww: app.NewWindow(),
	}
	fonts.Register()

	duo.gc = layout.NewContext(duo.ww.Queue())
	duo.cs = &duo.gc.Constraints

	// Layouts
	view := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical},
	}
	header := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween},
	}
	logo := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical},
	}
	body := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Horizontal},
	}
	sidebar := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical},
		i: layout.UniformInset(unit.Dp(8)),
	}
	content := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical},
		i: layout.UniformInset(unit.Dp(30)),
	}
	overview := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical},
	}
	overviewTop := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Horizontal},
	}
	sendReceive := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical},
		i: layout.UniformInset(unit.Dp(15)),
	}
	sendReceiveButtons := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Horizontal},
	}
	overviewBottom := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Horizontal},
	}
	status := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical},
		i: layout.UniformInset(unit.Dp(15)),
	}
	menu := DuoUIcomponent{
		l: layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle},
	}

	duo.comp = &DuoUIcomponents{
		view:               view,
		header:             header,
		logo:               logo,
		body:               body,
		sidebar:            sidebar,
		content:            content,
		overview:           overview,
		overviewTop:        overviewTop,
		sendReceive:        sendReceive,
		sendReceiveButtons: sendReceiveButtons,
		overviewBottom:     overviewBottom,
		status:             status,
		menu:               menu,
	}

	// Navigation
	duo.menu = &DuoUInav{
		current: "overview",
		//icoBackground: color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
		icoColor:    color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30},
		icoPadding:  unit.Dp(8),
		icoSize:     unit.Dp(64),
		overview:    *new(widget.Button),
		history:     *new(widget.Button),
		addressbook: *new(widget.Button),
		explorer:    *new(widget.Button),
		settings:    *new(widget.Button),
	}

	// Icons

	var err error
	ics := &DuoUIicons{}

	ics.Logo, err = material.NewIcon(ico.ParallelCoin)
	if err != nil {
		log.FATAL(err)
	}
	ics.Overview, err = material.NewIcon(icons.ActionHome)
	if err != nil {
		log.FATAL(err)
	}
	ics.History, err = material.NewIcon(icons.ActionHistory)
	if err != nil {
		log.FATAL(err)
	}
	ics.AddressBook, err = material.NewIcon(icons.ActionBook)
	if err != nil {
		log.FATAL(err)
	}
	ics.Explorer, err = material.NewIcon(icons.ActionExplore)
	if err != nil {
		log.FATAL(err)
	}
	ics.Network, err = material.NewIcon(icons.ActionFingerprint)
	if err != nil {
		log.FATAL(err)
	}
	ics.Settings, err = material.NewIcon(icons.ActionSettings)
	if err != nil {
		log.FATAL(err)
	}
	duo.ico = ics

	duo.th = material.NewTheme()
	return
}
