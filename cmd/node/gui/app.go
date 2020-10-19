package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (ng *NodeGUI) GetAppWidget() (a *p9.App) {
	a = ng.th.App()
	a.Pages(map[string]l.Widget{
		"main": a.VFlex().SpaceEvenly().Rigid(
			a.H2("first page").
				Alignment(text.Middle).
				Fn,
		).Fn,
		"second": a.VFlex().SpaceEvenly().Rigid(
			a.H2("second page").
				Alignment(text.Middle).
				Fn,
		).Fn,
		"third": a.VFlex().SpaceEvenly().Rigid(
			a.H2("third page").
				Alignment(text.Middle).
				Fn,
		).Fn,
		"fourth": a.VFlex().SpaceEvenly().Rigid(
			a.H2("fourth page").
				Alignment(text.Middle).
				Fn,
		).Fn,
		"fifth": a.VFlex().SpaceEvenly().Rigid(
			a.H2("fifth page").
				Alignment(text.Middle).
				Fn,
		).Fn,
	})
	a.SideBar([]l.Widget{
		ng.th.ButtonLayout(ng.sidebarButtons[0]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if a.ActivePageGet() == "main" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("first").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("main")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[1]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if a.ActivePageGet() == "second" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("second").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("second")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[2]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if a.ActivePageGet() == "third" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("third").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("third")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[3]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if a.ActivePageGet() == "fourth" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("fourth").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("fourth")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[4]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					if a.ActivePageGet() == "fifth" {
						background = "DocText"
						color = "DocBg"
					}
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("fifth").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("fifth")
				}).
			Fn,
		ng.th.ButtonLayout(ng.sidebarButtons[5]).
			Embed(
				func(gtx l.Context) l.Dimensions {
					background := "Transparent"
					color := "DocText"
					return ng.th.Fill(background,
						ng.th.Flex().Flexed(1,
							ng.th.Inset(0.5,
								ng.th.H6("invalid").
									Color(color).
									Fn,
							).Fn,
						).Fn,
					).Fn(gtx)
				},
			).
			Background("Transparent").
			SetClick(
				func() {
					if a.MenuOpen {
						a.MenuOpen = false
					}
					a.ActivePage("invalid")
				}).
			Fn,
	})
	return
}
