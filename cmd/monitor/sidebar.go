package monitor

import (
	"gioui.org/layout"

	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

func (s *State) Sidebar() layout.FlexChild {
	return gui.Rigid(func() {
		s.Gtx.Constraints.Width.Min = 332
		s.Gtx.Constraints.Width.Max = 332
		cs := s.Gtx.Constraints
		if s.Config.FilterOpen {
			s.FlexV(
				gui.Rigid(func() {
					s.Rectangle(cs.Width.Max, cs.Height.Max,
						"DocBg")
					s.Inset(4, func() {})
				}),
				gui.Flexed(1, func() {
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
					s.FlexV(
						gui.Flexed(1, func() {
							s.FlexV(gui.Flexed(1, func() {
								s.Loggers.GetWidget(s)
							}),
							)
						}),
					)
				}), gui.Rigid(func() {
					// if s.Config.FilterMode {
					// 	return
					// }
					s.Gtx.Constraints.Height.Max = 48
					s.Gtx.Constraints.Height.Min = 48
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					s.LevelsButtons()
				}), gui.Rigid(func() {
					s.Gtx.Constraints.Height.Max = 48
					s.Gtx.Constraints.Height.Min = 48
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					s.FlexH(
						gui.Rigid(func() {
							b := s.Buttons["FilterClear"]
							s.ButtonArea(func() {
								s.Inset(8, func() {
									s.Icon("Delete", "DocText", "DocBg", 32)
								})
							}, b)
							for b.Clicked(s.Gtx) {
								Debug("clear all")
								s.EntryBuf.Clear()
								s.FilterBuf.Clear()
							}
						}),
						gui.Flexed(1, func() {
						}),
						gui.Rigid(func() {
							b := s.Buttons["FilterHide"]
							s.ButtonArea(func() {
								s.Inset(8, func() {
									s.Icon("HideAll", "DocText", "DocBg", 32)
								})
							}, b)
							for b.Clicked(s.Gtx) {
								Debug("hide all")
								s.Loggers.CloseAllItems(s)
								s.SaveConfig()
							}
						}),
						gui.Rigid(func() {
							b := s.Buttons["FilterShow"]
							s.ButtonArea(func() {
								s.Inset(8, func() {
									s.Icon("ShowAll", "DocText", "DocBg", 32)
								})
							}, b)
							for b.Clicked(s.Gtx) {
								Debug("show all")
								s.Loggers.OpenAllItems(s)
								s.SaveConfig()
							}
						}),
						gui.Rigid(func() {
							b := s.Buttons["FilterAll"]
							s.ButtonArea(func() {
								s.Inset(8, func() {
									s.Icon("ShowItem", "DocText", "DocBg", 32)
								})
							}, b)
							for b.Clicked(s.Gtx) {
								Debug("filter all")
								s.Loggers.ShowAllItems(s)
								consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
								s.SaveConfig()
							}
						}),
						gui.Rigid(func() {
							b := s.Buttons["FilterNone"]
							s.ButtonArea(func() {
								s.Inset(8, func() {
									s.Icon("HideItem", "DocText", "DocBg", 32)
								})
							}, b)
							for b.Clicked(s.Gtx) {
								Debug("filter none")
								s.Loggers.HideAllItems(s)
								consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
								s.SaveConfig()
							}
						}),

						// Rigid(func() {
						//	s.IconButton("Filter", "DocBg", "DocText",
						//		&s.FilterButton)
						//	for s.FilterButton.Clicked(s.Gtx) {
						//		Debug("filter header clicked")
						//		if !s.Config.FilterOpen {
						//			s.Config.BuildOpen = false
						//			s.Config.SettingsOpen = false
						//		}
						//		s.Config.FilterOpen = !s.Config.FilterOpen
						//		s.SaveConfig()
						//	}
						// }),
					)
				}),
			)
		}
	})
}

func (s *State) LevelsButtons() {
	s.Lists["FilterLevel"].Layout(s.Gtx, len(logi.Tags)-1, func(a int) {
		bn := logi.Tags[logi.Levels[a+1]]
		color, bg := "PanelBg", "DocBg"
		if s.Config.FilterLevel > a {
			switch a + 1 {
			case 1:
				bg, color = "DocBgHilite", "Fatal"
			case 2:
				bg, color = "DocBgHilite", "Danger"
			case 3:
				bg, color = "DocBgHilite", "Check"
			case 4:
				bg, color = "DocBgHilite", "Warning"
			case 5:
				bg, color = "DocBgHilite", "Success"
			case 6:
				bg, color = "DocBgHilite", "Info"
			case 7:
				bg, color = "DocBgHilite", "Secondary"
			}
		}
		bb := &s.FilterLevelsButtons[a]
		s.ButtonArea(func() {
			cs := s.Gtx.Constraints
			cs.Width.Max = 48
			cs.Height.Max = 48
			s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
			s.Inset(8, func() {
				// cs := s.Gtx.Constraints
				s.Icon(bn, color, bg, 32)
			})
		}, bb)
		for bb.Clicked(s.Gtx) {
			s.Config.FilterLevel = a + 1
			*s.Ctx.Config.LogLevel = logi.Levels[s.Config.FilterLevel]
			if !s.Config.FilterMode {
				consume.SetLevel(s.Worker, logi.Levels[s.Config.FilterLevel])
				save.Pod(s.Ctx.Config)
			}
			Debug("filter level", logi.Tags[logi.Levels[a+1]])
			s.W.Invalidate()
			s.SaveConfig()
		}
	})
}
