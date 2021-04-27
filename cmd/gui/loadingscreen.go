package gui

import (
	"golang.org/x/exp/shiny/materialdesign/icons"

	l "github.com/p9c/gio/layout"

	"github.com/p9c/gel"
	"github.com/p9c/matrjoska/pkg/p9icons"
)

func (wg *WalletGUI) getLoadingPage() (a *gel.App) {
	a = wg.App(wg.Window.Width, wg.State.activePage, Break1).SetMainDirection(l.Center + 1)
	a.Pages(
		map[string]l.Widget{
			"loading": wg.Page(
				"loading", gel.Widgets{
					gel.WidgetSize{
						Widget:
						func(gtx l.Context) l.Dimensions {
							return a.Flex().Flexed(1, a.Direction().Center().Embed(a.H1("loading").Fn).Fn).Fn(gtx)
						},
					},
				},
			),
			"unlocking": wg.Page(
				"unlocking", gel.Widgets{
					gel.WidgetSize{
						Widget:
						func(gtx l.Context) l.Dimensions {
							return a.Flex().Flexed(1, a.Direction().Center().Embed(a.H1("unlocking").Fn).Fn).Fn(gtx)
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
			// wg.Flex().Rigid(wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn).Fn,
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
			// wg.Inset(0.5, gel.EmptySpace(0, 0)).Fn,
		},
	)
	// a.PushOverlay(wg.toasts.DrawToasts())
	// a.PushOverlay(wg.dialog.DrawDialog())
	return
}
