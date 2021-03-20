package gui

import (
	l "gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	p9icons "github.com/p9c/pod/pkg/gui/ico/svg"
)

func (wg *WalletGUI) getLoadingPage() (a *gui.App) {
	a = wg.App(&wg.Window.Width, wg.State.activePage, wg.invalidate, Break1).SetMainDirection(l.Center + 1)
	a.SetThemeHook(
		func() {
			D.Ln("theme hook")
			// D.Ln(wg.bools)
			*wg.cx.Config.DarkTheme = *wg.Dark
			a := wg.configs["config"]["DarkTheme"].Slot.(*bool)
			*a = *wg.Dark
			if wgb, ok := wg.config.Bools["DarkTheme"]; ok {
				wgb.Value(*wg.Dark)
			}
			save.Pod(wg.cx.Config)
		},
	)
	a.Pages(
		map[string]l.Widget{
			"home": wg.Page(
				"home", gui.Widgets{
					gui.WidgetSize{
						Widget:
						func(gtx l.Context) l.Dimensions {
							return a.Flex().Flexed(1, a.Direction().Center().Embed(a.H1("loading").Fn).Fn).Fn(gtx)
						},
					},
				},
			),
		},
	)
	a.ButtonBar(
		[]l.Widget{
			wg.PageTopBarButton(
				"home", 4, &icons.ActionLock, func(name string) {
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage, "Danger",
			),
			// wg.Flex().Rigid(wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn).Fn,
		},
	)
	a.StatusBar(
		[]l.Widget{
			wg.RunStatusPanel,
		},
		[]l.Widget{
			wg.StatusBarButton(
				"console", 2, &p9icons.Terminal, func(name string) {
					wg.MainApp.ActivePage(name)
				}, a,
			),
			wg.StatusBarButton(
				"log", 4, &icons.ActionList, func(name string) {
					D.Ln("click on button", name)
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage,
			),
			wg.StatusBarButton(
				"settings", 5, &icons.ActionSettings, func(name string) {
					wg.unlockPage.ActivePage(name)
				}, wg.unlockPage,
			),
			// wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
		},
	)
	// a.PushOverlay(wg.toasts.DrawToasts())
	// a.PushOverlay(wg.dialog.DrawDialog())
	return
}
