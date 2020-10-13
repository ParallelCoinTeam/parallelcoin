package kopach

import "github.com/p9c/pod/cmd/kopach/gui"

func (m *MinerModel) FlipTheme() {
	gui.Debug("dark theme:", m.DarkTheme)
	m.DarkTheme = !m.DarkTheme
	// Debug(s.Config.DarkTheme)
	m.SetTheme(m.DarkTheme)
	// s.SaveConfig()
}

func (m *MinerModel) SetTheme(dark bool) {
	m.Theme.Colors.SetTheme(dark)
}
