package duoui

func (ui *DuoUI) DuoUIcontent() func() {
	return func() {
		ui.rc.CurrentPage.Layout(ui.ly.Context)
	}
}
