package duoui

import (
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/theme"
	"github.com/p9c/pod/pkg/conte"
	"image/color"

	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/widget"
	"github.com/p9c/pod/pkg/fonts"
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuOuI(rc *rcd.RcVar, cx *conte.Xt) (duo *models.DuoUI, err error) {
	duo = &models.DuoUI{
		//cx: cx,
		CurrentPage: "overview",
		//rc: RcInit(),
		DuoUIwindow: app.NewWindow(
			app.Size(unit.Dp(900), unit.Dp(556)),
			app.Title("ParallelCoin"),
		),
		Quit:  make(chan struct{}),
		Ready: make(chan struct{}, 1),
	}
	fonts.Register()

	duo.DuoUIico = DuoIcons()
	// duo.ww.Width =  unit.Dp(800)
	//	duo.ww.Height=  unit.Dp(600)
	//		duo.ww.Title =  "ParallelCoin - True Story"

	duo.DuoUIcontext = layout.NewContext(duo.DuoUIwindow.Queue())
	duo.DuoUIconstraints = &duo.DuoUIcontext.Constraints

	navigations := make(map[string]*theme.DuoUIthemeNav)

	//	NavList: *new(widget.Enum),
	//
	navigations["mainMenu"] = mainMenu()

	//navButtons := make(map[string]*widget.Button)

	//for navItemKey, _ := range navItems {
	//	navButtons[navItemKey] = new(widget.Button)
	//}

	duo.DuoUIconfiguration = &models.DuoUIconfiguration{
		Abbrevation:        "DUO",
		PrimaryBgColor:     color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30},
		SecondaryBgColor:   color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
		PrimaryTextColor:   color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
		SecondaryTextColor: color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30},
		Navigations:        navigations,
	}

	//
	//navItems := map[string]theme.DuoUIbutton{
	//	"Overview":    theme.DuoUIbutton{
	//		Text:         "",
	//		TxColor:     color.RGBA{A: 0xff, R: 0x38, G: 0x11, B: 0x88},
	//		Font:         text.Font{},
	//		BgColor:      color.RGBA{A: 0xff, R: 0xff, G: 0x33, B: 0x33},
	//		CornerRadius: unit.Value{},
	//		Icon:         nil,
	//		Size:         unit.Value{},
	//		Padding:      unit.Value{},
	//	},
	//	"History": theme.DuoUIbutton{
	//		Text:         "",
	//		TxColor:      color.RGBA{A: 0xff, R: 0x38, G: 0x11, B: 0x88},
	//		Font:         text.Font{},
	//		BgColor:      color.RGBA{A: 0xff, R: 0x38, G: 0x30, B: 0x55},
	//		CornerRadius: unit.Value{},
	//		Icon:         nil,
	//		Size:         unit.Value{},
	//		Padding:      unit.Value{},
	//	},
	//	//"AddressBook": "AddressBook",
	//	//"Explorer":    "Explorer",
	//	//"Console":     "Console",
	//	//"Settings":    "Settings",
	//}

	//navMenus := make(map[string]*models.DuoUInav)

	// Layouts
	view := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
	}
	header := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Horizontal},
	}
	footer := models.DuoUIcomponent{
		Layout: layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly},
	}
	//logo := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Vertical},
	//}
	body := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Horizontal},
	}
	sidebar := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
		Inset:  layout.UniformInset(unit.Dp(8)),
	}
	content := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
		Inset:  layout.UniformInset(unit.Dp(30)),
	}
	overview := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Horizontal},
	}
	//sendReceive := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Vertical},
	//	Inset:  layout.UniformInset(unit.Dp(15)),
	//}
	//sendReceiveButtons := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Horizontal},
	//}
	status := models.DuoUIcomponent{
		Layout: layout.Flex{Axis: layout.Vertical},
		Inset:  layout.UniformInset(unit.Dp(15)),
	}
	statusItem := models.DuoUIcomponent{
		Layout: layout.Flex{Spacing: layout.SpaceBetween},
	}
	menu := models.DuoUIcomponent{
		Layout: layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly},
	}
	//
	//console := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Vertical},
	//	Inset:  layout.UniformInset(unit.Dp(15)),
	//}
	//consoleInput := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Horizontal},
	//}
	//consoleOutput := models.DuoUIcomponent{
	//	Layout: layout.Flex{Axis: layout.Horizontal},
	//}

	duo.DuoUIcomponents = &models.DuoUIcomponents{
		View:   view,
		Header: header,
		Footer: footer,
		//Logo:               logo,
		Body:     body,
		Sidebar:  sidebar,
		Content:  content,
		Overview: overview,
		//SendReceive:        sendReceive,
		//SendReceiveButtons: sendReceiveButtons,
		Status:     status,
		StatusItem: statusItem,
		Menu:       menu,
		//Console:            console,
		//ConsoleInput:       consoleInput,
		//ConsoleOutput:      consoleOutput,
	}

	// Navigation

	//duo.DuoUImenu = &models.DuoUInav{
	//	// icoBackground: color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
	//	IcoColor:   color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
	//	IcoPadding: unit.Dp(4),
	//	IcoSize:    unit.Dp(32),
	//	NavButtons: navButtons,
	//}

	// Icons

	// Settings tabs
	confTabs := make(map[string]*widget.Button)
	settingsFields := make(map[string]interface{})
	for _, group := range rcd.GetCoreSettings(cx).Schema.Groups {
		confTabs[group.Legend] = new(widget.Button)
		for _, field := range group.Fields {
			switch field.Type {
			case "array":
				settingsFields[field.Name] = new(widget.Button)
			case "input":
				settingsFields[field.Name] = &widget.DuoUIeditor{
					SingleLine: true,
					Submit:     true,
				}
			case "switch":
				settingsFields[field.Name] = new(widget.CheckBox)
			case "radio":
				settingsFields[field.Name] = new(widget.Enum)
			default:
				settingsFields[field.Name] = new(widget.Button)

			}

		}
	}
	duo.DuoUIconfiguration.Tabs = models.DuoUIconfTabs{
		Current:  "wallet",
		TabsList: confTabs,
	}
	duo.DuoUIconfiguration.Settings.Daemon.Widgets = settingsFields
	duo.DuoUItheme = theme.NewDuoUItheme()
	return
}

func mainMenu() *theme.DuoUIthemeNav {
	//mainNavButtons := make(map[string]*theme.DuoUIbutton)
	//
	//overviewIcon, _ := theme.NewDuoUIicon(icons.ActionHome)
	//overview := theme.DuoUIbutton{
	//	Text:   "Overview",
	//	Button: *new(widget.Button),
	//	Icon:   overviewIcon,
	//	Order:  0,
	//}
	//mainNavButtons["Overview"] = &overview
	//
	//historyIcon, _ := theme.NewDuoUIicon(icons.ActionHistory)
	//history := theme.DuoUIbutton{
	//	Text:   "History",
	//	Button: *new(widget.Button),
	//	Icon:   historyIcon,
	//	Order:  1,
	//}
	//mainNavButtons["History"] = &history

	//addressBookIcon, _ := theme.NewDuoUIicon(icons.ActionBook)
	//addressBook := theme.DuoUIbutton{
	//	Text:   "AddressBook",
	//	Button: new(widget.Button),
	//	Icon:   addressBookIcon,
	//}
	//mainNavButtons["AddressBook"] = &addressBook
	//
	//explorerIcon, _ := theme.NewDuoUIicon(icons.ActionExplore)
	//explorer := theme.DuoUIbutton{
	//	Text:   "Settings",
	//	Button: new(widget.Button),
	//	Icon:   explorerIcon,
	//}
	//mainNavButtons["Settings"] = &explorer
	//
	//consoleIcon, _ := theme.NewDuoUIicon(icons.ActionInput)
	//console := theme.DuoUIbutton{
	//	Text:   "Console",
	//	Button: new(widget.Button),
	//	Icon:   consoleIcon,
	//}
	//mainNavButtons["Console"] = &console
	//
	//networkIcon, _ := theme.NewDuoUIicon(icons.ActionFingerprint)
	//network := theme.DuoUIbutton{
	//	Text:   "Network",
	//	Button: new(widget.Button),
	//	Icon:   networkIcon,
	//}
	//mainNavButtons["Network"] = &network

	return &theme.DuoUIthemeNav{
		Title:         "mainMenu",
		IcoBackground: color.RGBA{A: 0xff, R: 0x38, G: 0x11, B: 0x88},
		IcoColor:      color.RGBA{A: 0xff, R: 0x38, G: 0x11, B: 0x88},
		IcoPadding:    unit.Dp(16),
		IcoSize:       unit.Dp(32),
		TxColor:       color.RGBA{A: 0xff, R: 0x38, G: 0x11, B: 0x88},
		BgColor:       color.RGBA{A: 0xff, R: 0xff, G: 0x33, B: 0x33},
		//NavButtons:    mainNavButtons,
	}
}
