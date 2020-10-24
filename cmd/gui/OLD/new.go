package gui

import (
	"fmt"

	"gioui.org/unit"
	"gioui.org/widget"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/gui/wallet/dap/mod"
	"github.com/p9c/pod/pkg/gui/wallet/dap/page"
	"github.com/p9c/pod/pkg/gui/wallet/dap/res"
	"github.com/p9c/pod/pkg/gui/wallet/nav"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

func (w *Worker) NewGuiApp() *GuiAppModel {

	// ww := map[string]*win.Window{
	//	"main": &win.Window{
	//		W: app.NewWindow(
	//			app.Size(unit.Dp(1024), unit.Dp(800)),
	//			app.Title("ParallelCoin"),
	//		)},
	// }

	n := &nav.Navigation{
		Name: "Navigacion",
		// Bg:           d.UI.Theme.Colors["NavBg"],
		ItemIconSize: unit.Px(24),
	}

	// s := &mod.Settings{
	//	Dir: appdata.Dir("dap", false),
	// }
	// d.S = s

	// Debug("New DAP", d)

	th := p9.NewTheme(p9fonts.Collection(), w.quit)
	solButtons := make([]*p9.Clickable, 201)
	for i := range solButtons {
		solButtons[i] = th.Clickable()
	}
	lists := map[string]*p9.List{
		"latestTransactions": th.List().Vertical().Start(),
	}
	g := &GuiAppModel{
		Cx:        w.cx,
		worker:    w,
		Theme:     th,
		DarkTheme: *w.cx.Config.DarkTheme,
		logoButton: th.Clickable().SetClick(func() {
			Debug("clicked logo button")
		}),
		mineToggle:      th.Bool(*w.cx.Config.Generate),
		cores:           th.Float().SetValue(float32(*w.cx.Config.GenThreads)),
		solButtons:      solButtons,
		lists:           lists,
		unhideClickable: th.Clickable(),
		modalScrim:      th.Clickable(),
		modalClose:      th.Clickable(),
		threadsMax:      th.Clickable(),
		threadsMin:      th.Clickable(),
		ui: &mod.UserInterface{
			Theme: theme.NewTheme(),
			N:     n,
			R:     res.Resposnsivity(0, 0),
			// W: &win.Windows{
			//	W: ww,
			// },
		},
	}
	g.SetTheme(g.DarkTheme)
	g.pass = th.Editor().Mask('•').SingleLine().Submit(true)
	g.passInput = th.SimpleInput(g.pass).Color("DocText")
	g.unhideButton = th.IconButton(g.unhideClickable).
		Background("").
		Color("Primary").
		Icon(g.Icon().Src(icons2.ActionVisibility))
	showClickableFn := func() {
		g.hide = !g.hide
		if !g.hide {
			g.unhideButton.Color("Primary").Icon(g.Icon().Src(icons2.ActionVisibility))
			g.pass.Mask('•')
			g.passInput.Color("Primary")
		} else {
			g.unhideButton.Color("DocText").Icon(g.Icon().Src(icons2.ActionVisibilityOff))
			g.pass.Mask(0)
			g.passInput.Color("DocText")
		}
	}
	g.unhideClickable.SetClick(showClickableFn)
	g.pass.SetText(*w.cx.Config.MinerPass).Mask('•').SetSubmit(func(txt string) {
		if !g.hide {
			showClickableFn()
		}
		showClickableFn()
		go func() {
			*w.cx.Config.MinerPass = txt
			save.Pod(w.cx.Config)
			w.Stop()
			w.Start()
		}()
	}).SetChange(func(txt string) {
		// send keystrokes to the NSA
	})
	for i := 0; i < 201; i++ {
		g.solButtons[i] = th.Clickable()
	}
	g.logoButton.SetClick(
		func() {
			g.FlipTheme()
			Info("clicked logo button")
			showClickableFn()
			showClickableFn()
		})

	g.ui.N.MenuItems = append(g.getMenuItems(g.ui))
	g.ui.N.CurrentPage = page.Page{
		Title:  "Overview",
		Header: g.overviewHeader(),
		Body:   g.overviewBody(),
		Footer: noReturn,
	}
	g.ui.F = g.footer()

	return g
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}

func (g *GuiAppModel) getMenuItem(hide bool, title string, header, body, footer func(gtx C) D) nav.Item {
	return nav.Item{
		Title: title,
		Icon:  g.ui.Theme.Icons[title],
		Btn:   new(widget.Clickable),
		Page: page.Page{
			Title:  title,
			Header: header,
			Body:   body,
			Footer: footer,
		},
		HideOnMob: false,
	}
}

func (g *GuiAppModel) getMenuItems(ui *mod.UserInterface) []nav.Item {
	return []nav.Item{
		g.getMenuItem(false, "Overview", g.overviewHeader(), g.overviewBody(), noReturn),
		g.getMenuItem(false, "Send", g.sendHeader(), g.sendBody(), g.sendFooter()),
		g.getMenuItem(false, "Receive", g.receiveHeader(), g.receiveBody(), noReturn),
		g.getMenuItem(false, "Transactions", g.transactionsHeader(), g.transactionsBody(), noReturn),
		// g.getMenuItem(true, "Explore", g.exploreHeader(), g.exploreBody(), noReturn),
		// g.getMenuItem(true, "Peers", g.peersHeader(), g.peersBody(), noReturn),
		// g.getMenuItem(true, "Settings", g.settingsHeader(), g.settingsBody(), noReturn),
	}
}
