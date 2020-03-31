package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/consume"
)

func (s *State) Sidebar() layout.FlexChild {
	return Rigid(func() {
		//if !(s.Config.BuildOpen || s.Config.SettingsOpen) {
		//	s.Gtx.Constraints.Width.Max /= 2
		//} else {
		//	s.Gtx.Constraints.Width.Max -= 340
		//}
		s.Gtx.Constraints.Width.Min = 332
		s.Gtx.Constraints.Width.Max = 332
		//if s.Gtx.Constraints.Width.Max > 360 {
		//	s.Gtx.Constraints.Width.Max = 360
		//}
		if s.Config.FilterOpen {
			s.FlexV(Rigid(func() {
				s.Gtx.Constraints.Height.Max = 48
				s.Gtx.Constraints.Height.Min = 48
				s.FlexH(
					Flexed(1, func() {
						cs := s.Gtx.Constraints
						s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
						//if s.WindowWidth > 640 {
						//	s.Label("Filter")
						//}
					}), Rigid(func() {
						s.IconButton("Send", "DocText", "DocBg",
							&s.FilterSendButton)
						for s.FilterSendButton.Clicked(s.Gtx) {
							Debug("send current log buffer")
							//s.EntryBuf.Clear()
						}
					}), Rigid(func() {
						s.IconButton("Delete", "DocText", "DocBg",
							&s.FilterClearButton)
						for s.FilterClearButton.Clicked(s.Gtx) {
							Debug("clear all")
							s.EntryBuf.Clear()
						}
					}), Rigid(func() {
						s.IconButton("HideAll", "DocText", "DocBg",
							&s.FilterHideButton)
						for s.FilterHideButton.Clicked(s.Gtx) {
							Debug("hide all")
							s.Loggers.CloseAllItems(s)
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("ShowAll", "DocText", "DocBg",
							&s.FilterShowButton)
						for s.FilterShowButton.Clicked(s.Gtx) {
							Debug("show all")
							s.Loggers.OpenAllItems(s)
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("ShowItem", "DocText", "DocBg",
							&s.FilterAllButton)
						for s.FilterAllButton.Clicked(s.Gtx) {
							Debug("filter all")
							s.Loggers.ShowAllItems(s)
							consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("HideItem", "DocText", "DocBg",
							&s.FilterNoneButton)
						for s.FilterNoneButton.Clicked(s.Gtx) {
							Debug("filter none")
							s.Loggers.HideAllItems(s)
							consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
							s.SaveConfig()
						}
					}), Rigid(func() {
						s.IconButton("Filter", "DocText", "DocBg",
							&s.FilterButton)
						for s.FilterButton.Clicked(s.Gtx) {
							Debug("filter header clicked")
							if !s.Config.FilterOpen {
								s.Config.BuildOpen = false
								s.Config.SettingsOpen = false
							}
							s.Config.FilterOpen = !s.Config.FilterOpen
							s.SaveConfig()
						}
					}),
				)
			}), Rigid(func() {
				s.Gtx.Constraints.Height.Max = 48
				s.Gtx.Constraints.Height.Min = 48
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
				s.LevelsButtons()
			}), Flexed(1, func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
				s.Inset(16, func() {
					s.FlexV(
						Flexed(1, func() {
							//s.Gtx.Constraints.Width.Min = 240
							s.FlexV(Flexed(1, func() {
								s.Loggers.GetWidget(s)
							}))
						}),
					)
				})
			}),
			)
		}
	})
}

func (s *State) LevelsButtons() {
	s.FilterLevelList.Layout(s.Gtx, len(logi.Tags)-1, func(a int) {
		bn := logi.Tags[logi.Levels[a+1]]
		bg := "PanelBg"
		color := "DocBg"
		if s.Config.FilterLevel > a {
			switch a + 1 {
			case 1:
				bg, color = "DocBg", "Danger"
			case 2:
				bg, color = "DocBg", "Danger"
			case 3:
				bg, color = "DocBg", "Check"
			case 4:
				bg, color = "DocBg", "Warning"
			case 5:
				bg, color = "DocBg", "Success"
			case 6:
				bg, color = "DocBg", "Info"
			case 7:
				bg, color = "DocBg", "Secondary"
			}
		}
		bb := &s.FilterLevelsButtons[a]
		s.IconButton(bn, color, bg, bb)
		for bb.Clicked(s.Gtx) {
			s.Config.FilterLevel = a + 1
			*s.Ctx.Config.LogLevel = logi.Levels[s.Config.FilterLevel]
			consume.SetLevel(s.Worker, logi.Levels[s.Config.FilterLevel])
			Debug("filter level", logi.Tags[logi.Levels[a+1]])
			s.W.Invalidate()
			save.Pod(s.Ctx.Config)
			s.SaveConfig()
		}
	})
}
