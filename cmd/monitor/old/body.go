package old

//
//func (s *State) Body(headless bool) layout.FlexChild {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	return gui.Flexed(1, func() {
//		cs := gtx.Constraints
//		cs.Width.Min = cs.Width.Max / 2
//		s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg", "ff")
//		s.Lists["Log"].Axis = layout.Vertical
//		s.Lists["Log"].ScrollToEnd = true
//		s.Lists["Log"].Layout(gtx, s.EntryBuf.Len(), func(i int) {
//			if s.EntryBuf.Clicked == i {
//				cs := gtx.Constraints
//				//cs.Height.Max = 48
//				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocBgHilite", "ff")
//			}
//			b := s.EntryBuf.Get(i)
//			color := "DocText"
//			//fmt.Println("level", b.Level)
//			switch b.Level {
//			case logi.Trace:
//				color = "Secondary"
//			case logi.Debug:
//				color = "Info"
//			case logi.Info:
//				color = "Success"
//			case logi.Warn:
//				color = "Warning"
//			case logi.Check:
//				color = "Check"
//			case logi.Error:
//				color = "Danger"
//			case logi.Fatal:
//				color = "Fatal"
//			}
//			button := s.EntryBuf.GetButton(i)
//			hider := s.EntryBuf.GetHider(i)
//			ww := s.WindowWidth
//			if s.Config.FilterOpen {
//				ww -= 332
//			}
//			s.FlexH(
//				gui.Flexed(1, func() {
//					s.ButtonArea(func() {
//						s.Inset(8, func() {
//							s.FlexH(
//								gui.Rigid(func() {
//									if ww > 480 {
//										s.Icon(logi.Tags[b.Level], color,
//											"Transparent", 32)
//									}
//								}),
//								gui.Rigid(func() {
//									if ww > 960 {
//										s.FlexH(gui.Rigid(
//											s.Text(b.Time.Format("15:04:05"),
//												color, "Transparent",
//												"Mono", "body2"),
//										))
//									}
//								}),
//								gui.Flexed(1, func() {
//									//cs := gtx.Constraints
//									//s.Rectangle(cs.Width.Max, cs.Height.Max,
//									//	"PanelBg", "ff")
//									tc := "DocText"
//									if ww <= 480 {
//										tc = color
//									}
//									s.FlexH(gui.Rigid(
//										s.Text(b.Text, tc, "Transparent",
//											"Mono", "body2"),
//									))
//								}),
//								//s.Spacer(),
//								gui.Rigid(func() {
//									if ww > 720 {
//										s.FlexH(gui.Rigid(
//											s.Text(b.Package, "PanelBg",
//												"Transparent", "Primary",
//												"body1"),
//										))
//									}
//								}),
//							)
//
//						})
//					}, button)
//					for button.Clicked(gtx) {
//						go func() {
//							if s.Config.ClickCommand == "" {
//								return
//							}
//							s.EntryBuf.Clicked = i
//							split := strings.Split(b.CodeLocation, ":")
//							v1 := split[0]
//							v2 := split[1]
//							c := strings.Replace(s.Config.ClickCommand, "$1", v1, 1)
//							c = strings.Replace(c, "$2", v2, 1)
//							Debug("running command", c)
//							args := strings.Split(c, " ")
//							cmd := exec.Command(args[0], args[1:]...)
//							_ = cmd.Run()
//							//s.Config.ClickCommand
//						}()
//					}
//				}),
//				gui.Rigid(func() {
//					if ww > 640 {
//						s.ButtonArea(func() {
//							s.Inset(4, func() {
//								s.Icon("HideItem", "PanelBg", "DocBg", 32)
//							})
//						}, hider)
//						for hider.Clicked(gtx) {
//							s.Config.FilterNodes[s.EntryBuf.Get(i).
//								Package].Hidden = true
//						}
//					}
//				}),
//			)
//		})
//		s.W.Invalidate()
//	})
//}
