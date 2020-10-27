package gui

import "github.com/p9c/pod/app/save"

func (m *GuiAppModel) FlipTheme() {
	m.DarkTheme = !m.DarkTheme
	Debug("dark theme:", m.DarkTheme)
	m.SetTheme(m.DarkTheme)
}

func (m *GuiAppModel) SetTheme(dark bool) {
	m.Theme.Colors.SetTheme(dark)
	*m.Cx.Config.DarkTheme = dark
	save.Pod(m.Cx.Config)
}
