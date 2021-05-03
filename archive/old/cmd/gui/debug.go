package gui

// func (wg *WalletGUI) goRoutines() {
// 	var e error
// 	if wg.App.ActivePageGet() == "goroutines" || wg.unlockPage.ActivePageGet() == "goroutines" {
// 		D.Ln("updating goroutines data")
// 		var b []byte
// 		buf := bytes.NewBuffer(b)
// 		if e = pprof.Lookup("goroutine").WriteTo(buf, 2); E.Chk(e) {
// 		}
// 		lines := strings.Split(buf.String(), "\n")
// 		var out []l.Widget
// 		var clickables []*p9.Clickable
// 		for x := range lines {
// 			i := x
// 			clickables = append(clickables, wg.Clickable())
// 			var text string
// 			if strings.HasPrefix(lines[i], "goroutine") && i < len(lines)-2 {
// 				text = lines[i+2]
// 				text = strings.TrimSpace(strings.Split(text, " ")[0])
// 				// outString += text + "\n"
// 				out = append(
// 					out, func(gtx l.Context) l.Dimensions {
// 						return wg.ButtonLayout(clickables[i]).Embed(
// 							wg.ButtonInset(
// 								0.25,
// 								wg.Caption(text).
// 									Color("DocText").Fn,
// 							).Fn,
// 						).Background("Transparent").SetClick(
// 							func() {
// 								go func() {
// 									out := make([]string, 2)
// 									split := strings.Split(text, ":")
// 									if len(split) > 2 {
// 										out[0] = strings.Join(split[:len(split)-1], ":")
// 										out[1] = split[len(split)-1]
// 									} else {
// 										out[0] = split[0]
// 										out[1] = split[1]
// 									}
// 									D.Ln("path", out[0], "line", out[1])
// 									goland := "goland64.exe"
// 									if runtime.GOOS != "windows" {
// 										goland = "goland"
// 									}
// 									launch := exec.Command(goland, "--line", out[1], out[0])
// 									if e = launch.Start(); E.Chk(e) {
// 									}
// 								}()
// 							},
// 						).
// 							Fn(gtx)
// 					},
// 				)
// 			}
// 		}
// 		// D.Ln(outString)
// 		wg.State.SetGoroutines(out)
// 		wg.invalidate <- struct{}{}
// 	}
// }
