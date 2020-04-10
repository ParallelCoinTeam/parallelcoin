package monitor

//
//func (s *State) BottomBar(headless bool) layout.FlexChild {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	return gui.Rigid(func() {
//		cs := gtx.Constraints
//		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
//		s.FlexV(
//			s.SettingsPage(headless),
//			s.BuildPage(headless),
//			s.StatusBar(headless),
//		)
//	})
//}
//
//func (s *State) StatusBar(headless bool) layout.FlexChild {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	return gui.Rigid(func() {
//		cs := gtx.Constraints
//		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
//		s.FlexH(
//			s.RunControls(headless),
//			s.RunmodeButtons(headless),
//			//s.Spacer("PanelBg"),
//			gui.Flexed(1, func() {
//				gtx.Constraints.Height.Max = 48
//				cs := gtx.Constraints
//				s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "FF")
//			}),
//			s.BuildButtons(headless),
//			s.SettingsButtons(headless),
//			s.Filter(headless),
//		)
//	})
//}
//
//func (s *State) RunmodeButtons(headless bool) layout.FlexChild {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	return gui.Rigid(func() {
//		fg, bg := "ButtonText", "ButtonBg"
//		s.FlexH(gui.Rigid(func() {
//			if !s.Config.RunModeOpen {
//				txt := s.Config.RunMode
//				if s.Config.Running {
//					cs := gtx.Constraints
//					bg, fg = "DocBg", "DocText"
//					s.Rectangle(cs.Width.Min, 48, bg, "ff")
//					s.Label(txt, fg, bg)
//				} else {
//					b := s.Buttons["RunModeFold"]
//					s.TextButton(txt, "Secondary", 34, fg, bg, b)
//					for b.Clicked(s.Gtx) {
//						if !s.Config.Running {
//							s.Config.RunModeOpen = true
//							s.SaveConfig()
//						}
//					}
//				}
//			} else {
//				modes := []string{
//					"node", "wallet", "shell", "gui", "mon",
//				}
//				s.Lists["Modes"].Layout(s.Gtx, len(modes), func(i int) {
//					mm := modes[i]
//					fg := "DocBg"
//					if modes[i] == s.Config.RunMode {
//						fg = "DocText"
//					}
//					txt := mm
//					if s.WindowWidth <= 880 && s.Config.FilterOpen ||
//						s.WindowWidth <= 640 && !s.Config.FilterOpen {
//						txt = txt[:1]
//					}
//					cs := gtx.Constraints
//					s.Rectangle(cs.Width.Max, cs.Height.Max, "ButtonBg",
//						"ff")
//					s.TextButton(txt, "Secondary", 34, fg,
//						"ButtonBg", s.ModesButtons[modes[i]])
//					for s.ModesButtons[modes[i]].Clicked(s.Gtx) {
//						Debug(mm, "clicked")
//						if s.Config.RunModeOpen {
//							s.Config.RunMode = modes[i]
//							s.Config.RunModeOpen = false
//						}
//						s.SaveConfig()
//					}
//				})
//			}
//		}),
//		)
//	})
//}
//
//func (s *State) Filter(headless bool) layout.FlexChild {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	return gui.Rigid(func() {
//		fg, bg := "PanelText", "PanelBg"
//		if s.Config.FilterOpen {
//			fg, bg = "DocText", "DocBg"
//		}
//		b := s.Buttons["Filter"]
//		s.ButtonArea(func() {
//			cs := gtx.Constraints
//			s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
//			s.Inset(8, func() {
//				s.Icon("Filter", fg, bg, 32)
//			})
//		}, b)
//		for b.Clicked(s.Gtx) {
//			Debug("clicked filter button")
//			if !s.Config.FilterOpen {
//				s.Config.SettingsOpen = false
//				s.Config.BuildOpen = false
//			}
//			s.Config.FilterOpen = !s.Config.FilterOpen
//			s.SaveConfig()
//		}
//		//}
//	})
//}
