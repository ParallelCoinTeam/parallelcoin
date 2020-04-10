package old

//
//import (
//	"gioui.org/layout"
//	"github.com/p9c/pod/app/save"
//	"github.com/p9c/pod/pkg/gui"
//	"github.com/p9c/pod/pkg/util/logi"
//	"github.com/p9c/pod/pkg/util/logi/consume"
//)
//
//func (s *State) Sidebar(headless bool) layout.FlexChild {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	return gui.Rigid(func() {
//		//if !(s.Config.BuildOpen || s.Config.SettingsOpen) {
//		//	gtx.Constraints.Width.Max /= 2
//		//} else {
//		//	gtx.Constraints.Width.Max -= 340
//		//}
//		gtx.Constraints.Width.Min = 332
//		gtx.Constraints.Width.Max = 332
//		//if gtx.Constraints.Width.Max > 360 {
//		//	gtx.Constraints.Width.Max = 360
//		//}
//		cs := gtx.Constraints
//		if s.Config.FilterOpen {
//			s.FlexV(
//				gui.Rigid(func() {
//					s.Rectangle(cs.Width.Max, cs.Height.Max,
//						"DocBg", "ff")
//					s.Inset(4, func() {})
//				}),
//				gui.Flexed(1, func() {
//					cs := gtx.Constraints
//					s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
//					s.Inset(8, func() {
//						s.FlexV(
//							gui.Flexed(1, func() {
//								//gtx.Constraints.Width.Min = 240
//								s.FlexV(gui.Flexed(1, func() {
//									s.Loggers.GetWidget(s, headless)
//								}),
//								)
//							}),
//						)
//					})
//				}), gui.Rigid(func() {
//					gtx.Constraints.Height.Max = 48
//					gtx.Constraints.Height.Min = 48
//					cs := gtx.Constraints
//					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
//					s.LevelsButtons(headless)
//				}), gui.Rigid(func() {
//					gtx.Constraints.Height.Max = 48
//					gtx.Constraints.Height.Min = 48
//					cs := gtx.Constraints
//					s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg",
//						"ff")
//					s.FlexH(
//						gui.Rigid(func() {
//							b := s.Buttons["FilterClear"]
//							s.ButtonArea(func() {
//								s.Inset(8, func() {
//									s.Icon("Delete", "DocText", "DocBg", 32)
//								})
//							}, b)
//							for b.Clicked(gtx) {
//								Debug("clear all")
//								s.EntryBuf.Clear()
//							}
//						}),
//						gui.Flexed(1, func() {
//						}),
//						gui.Rigid(func() {
//							b := s.Buttons["FilterHide"]
//							s.ButtonArea(func() {
//								s.Inset(8, func() {
//									s.Icon("HideAll", "DocText", "DocBg", 32)
//								})
//							}, b)
//							for b.Clicked(gtx) {
//								Debug("hide all")
//								s.Loggers.CloseAllItems(s)
//								s.SaveConfig()
//							}
//						}),
//						gui.Rigid(func() {
//							b := s.Buttons["FilterShow"]
//							//s.IconButton("ShowAll", "DocText", "DocBg", b)
//							s.ButtonArea(func() {
//								s.Inset(8, func() {
//									s.Icon("ShowAll", "DocText", "DocBg", 32)
//								})
//							}, b)
//							for b.Clicked(gtx) {
//								Debug("show all")
//								s.Loggers.OpenAllItems(s)
//								s.SaveConfig()
//							}
//						}),
//						gui.Rigid(func() {
//							b := s.Buttons["FilterAll"]
//							//s.IconButton("ShowItem", "DocText", "DocBg", b)
//							s.ButtonArea(func() {
//								s.Inset(8, func() {
//									s.Icon("ShowItem", "DocText", "DocBg", 32)
//								})
//							}, b)
//							for b.Clicked(gtx) {
//								Debug("filter all")
//								s.Loggers.ShowAllItems(s)
//								consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
//								s.SaveConfig()
//							}
//						}),
//						gui.Rigid(func() {
//							b := s.Buttons["FilterNone"]
//							s.ButtonArea(func() {
//								s.Inset(8, func() {
//									s.Icon("HideItem", "DocText", "DocBg", 32)
//								})
//							}, b)
//							//s.IconButton("HideItem", "DocText", "DocBg", b)
//							for b.Clicked(gtx) {
//								Debug("filter none")
//								s.Loggers.HideAllItems(s)
//								consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
//								s.SaveConfig()
//							}
//						}),
//
//						//Rigid(func() {
//						//	s.IconButton("Filter", "DocBg", "DocText",
//						//		&s.FilterButton)
//						//	for s.FilterButton.Clicked(gtx) {
//						//		Debug("filter header clicked")
//						//		if !s.Config.FilterOpen {
//						//			s.Config.BuildOpen = false
//						//			s.Config.SettingsOpen = false
//						//		}
//						//		s.Config.FilterOpen = !s.Config.FilterOpen
//						//		s.SaveConfig()
//						//	}
//						//}),
//					)
//				}),
//			)
//		}
//	})
//}
//
//func (s *State) LevelsButtons(headless bool) {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	s.Lists["FilterLevel"].Layout(gtx, len(logi.Tags)-1, func(a int) {
//		bn := logi.Tags[logi.Levels[a+1]]
//		color, bg := "PanelBg", "DocBg"
//		if s.Config.FilterLevel > a {
//			switch a + 1 {
//			case 1:
//				bg, color = "DocBgHilite", "Fatal"
//			case 2:
//				bg, color = "DocBgHilite", "Danger"
//			case 3:
//				bg, color = "DocBgHilite", "Check"
//			case 4:
//				bg, color = "DocBgHilite", "Warning"
//			case 5:
//				bg, color = "DocBgHilite", "Success"
//			case 6:
//				bg, color = "DocBgHilite", "Info"
//			case 7:
//				bg, color = "DocBgHilite", "Secondary"
//			}
//		}
//		bb := &s.FilterLevelsButtons[a]
//		s.ButtonArea(func() {
//			cs := gtx.Constraints
//			cs.Width.Max = 48
//			cs.Height.Max = 48
//			s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
//			s.Inset(8, func() {
//				//cs := gtx.Constraints
//				s.Icon(bn, color, bg, 32)
//			})
//		}, bb)
//		for bb.Clicked(gtx) {
//			s.Config.FilterLevel = a + 1
//			*s.Ctx.Config.LogLevel = logi.Levels[s.Config.FilterLevel]
//			consume.SetLevel(s.Worker, logi.Levels[s.Config.FilterLevel])
//			Debug("filter level", logi.Tags[logi.Levels[a+1]])
//			s.W.Invalidate()
//			save.Pod(s.Ctx.Config)
//			s.SaveConfig()
//		}
//	})
//}
