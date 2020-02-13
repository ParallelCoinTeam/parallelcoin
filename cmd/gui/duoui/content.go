package duoui

func  (ui *DuoUI)DuoUIcontent() func() {
	return func() {
		ui.ly.Pages[ui.rc.ShowPage].Layout(ui.ly.Context)
	}
}
