package main

import (
	"github.com/p9c/qu"

	l "github.com/p9c/gio/layout"

	"github.com/p9c/gel"
	"github.com/p9c/gel/clipboard"
)

type State struct {
	*gel.Window
	evKey              int
	showClicker        *gel.Clickable
	showPrimaryClicker *gel.Clickable
	showText           *string
	editor             *gel.Input
}

func NewState(quit qu.C) (s *State) {
	s = &State{
		Window: gel.NewWindowP9(quit),
	}
	txt := ""
	s.showText = &txt
	s.showClicker = s.WidgetPool.GetClickable()
	s.showPrimaryClicker = s.WidgetPool.GetClickable()
	s.editor = s.Input(
		"",
		"type something to test editor things",
		"DocText",
		"PanelBg",
		"DocBg",
		func(amt string) {},
		func(string) {},
	)
	return
}

func main() {
	quit := qu.T()
	state := NewState(quit)
	*state.showText = "clipboard demo/test"
	var e error
	if e = state.Window.
		Size(48, 32).
		Title("hello world").
		Open().
		Run(state.rootWidget, quit.Q, quit); E.Chk(e) {
	}
}

func (s *State) rootWidget(gtx l.Context) l.Dimensions {
	return s.Direction().Center().
		Embed(
			s.Inset(0.5,
				s.Border().Color("DocText").Embed(
					s.VFlex().
						Flexed(0.5,
							s.Inset(0.5,
								s.Border().Color("DocText").Embed(
									s.VFlex().
										Rigid(
											s.Direction().Center().Embed(
												s.Flex().
													Rigid(
														s.Inset(0.5,
															s.ButtonLayout(
																s.showClicker.
																	SetClick(func() {
																		I.Ln("user clicked show clipboard button")
																		s.ClipboardReadReqs <- func(cs string) {
																			*s.showText = cs
																			I.Ln("clipboard contents:", cs)
																		}
																	}),
															).CornerRadius(0.25).Corners(^0).
																Embed(
																	s.Border().CornerRadius(0.25).Color("DocText").Embed(
																		s.Inset(0.5,
																			s.H6("show clipboard").
																				Fn,
																		).Fn,
																	).Fn,
																).Fn,
														).Fn,
													).
													Rigid(
														s.Inset(0.5,
															s.ButtonLayout(
																s.showPrimaryClicker.
																	SetClick(func() {
																		*s.showText = clipboard.GetPrimary()
																		I.Ln("clipboard contents:", *s.showText)
																	}),
															).CornerRadius(0.25).Corners(^0).
																Embed(
																	s.Border().CornerRadius(0.25).Color("DocText").Embed(
																		s.Inset(0.5,
																			s.H6("show primary").
																				Fn,
																		).Fn,
																	).Fn,
																).Fn,
														).Fn,
													).Fn,
											).Fn,
										).
										Rigid(
											s.Inset(0.5,
												s.editor.Fn,
											).Fn,
										).Fn,
								).Fn,
							).Fn,
						).
						Flexed(0.5,
							s.Inset(0.5,
								s.Border().Color("DocText").Embed(
									s.Flex().Flexed(1,
										s.Direction().NW().Embed(
											s.Inset(0.5,
												s.Body1(*s.showText).
													Fn,
											).Fn,
										).Fn,
									).
										Fn,
								).Fn,
							).Fn,
						).Fn,
				).Fn,
			).Fn,
		).
		Fn(gtx)
}
