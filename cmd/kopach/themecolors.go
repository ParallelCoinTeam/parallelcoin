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
	if !dark {
		m.Theme.Colors["DocText"] = m.Theme.Colors["dark"]
		m.Theme.Colors["DocBg"] = m.Theme.Colors["light"]
		m.Theme.Colors["PanelText"] = m.Theme.Colors["dark"]
		m.Theme.Colors["PanelBg"] = m.Theme.Colors["white"]
		m.Theme.Colors["PanelTextDim"] = m.Theme.Colors["dark-grayii"]
		m.Theme.Colors["PanelBgDim"] = m.Theme.Colors["dark-grayi"]
		m.Theme.Colors["DocTextDim"] = m.Theme.Colors["light-grayi"]
		m.Theme.Colors["DocBgDim"] = m.Theme.Colors["dark-grayi"]
		m.Theme.Colors["Warning"] = m.Theme.Colors["light-orange"]
		m.Theme.Colors["Success"] = m.Theme.Colors["dark-green"]
		m.Theme.Colors["Check"] = m.Theme.Colors["orange"]
		m.Theme.Colors["DocBgHilite"] = m.Theme.Colors["dark-white"]
	} else {
		m.Theme.Colors["DocText"] = m.Theme.Colors["light"]
		m.Theme.Colors["DocBg"] = m.Theme.Colors["black"]
		m.Theme.Colors["PanelText"] = m.Theme.Colors["light"]
		m.Theme.Colors["PanelBg"] = m.Theme.Colors["dark"]
		m.Theme.Colors["PanelTextDim"] = m.Theme.Colors["light-grayii"]
		m.Theme.Colors["PanelBgDim"] = m.Theme.Colors["light-gray"]
		m.Theme.Colors["DocTextDim"] = m.Theme.Colors["light-gray"]
		m.Theme.Colors["DocBgDim"] = m.Theme.Colors["light-grayii"]
		m.Theme.Colors["Warning"] = m.Theme.Colors["yellow"]
		m.Theme.Colors["Success"] = m.Theme.Colors["green"]
		m.Theme.Colors["Check"] = m.Theme.Colors["orange"]
		m.Theme.Colors["DocBgHilite"] = m.Theme.Colors["light-black"]
	}
}
