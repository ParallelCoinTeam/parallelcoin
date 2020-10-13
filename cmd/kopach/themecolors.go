package kopach

import "github.com/p9c/pod/app/save"

func (m *MinerModel) FlipTheme() {
	m.DarkTheme = !m.DarkTheme
	Debug("dark theme:", m.DarkTheme)
	m.SetTheme(m.DarkTheme)
}

func (m *MinerModel) SetTheme(dark bool) {
	m.Theme.Colors.SetTheme(dark)
	*m.Cx.Config.DarkTheme = m.DarkTheme
	save.Pod(m.Cx.Config)
}
