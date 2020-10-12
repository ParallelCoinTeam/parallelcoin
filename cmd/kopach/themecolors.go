package kopach

import "github.com/p9c/pod/cmd/kopach/gui"

func (s *MinerModel) FlipTheme() {
	gui.Debug("dark theme:", s.DarkTheme)
	s.DarkTheme = !s.DarkTheme
	// Debug(s.Config.DarkTheme)
	s.SetTheme(s.DarkTheme)
	// s.SaveConfig()
}

func (s *MinerModel) SetTheme(dark bool) {
	if !dark {
		s.Theme.Colors["DocText"] = s.Theme.Colors["dark"]
		s.Theme.Colors["DocBg"] = s.Theme.Colors["light"]
		s.Theme.Colors["PanelText"] = s.Theme.Colors["dark"]
		s.Theme.Colors["PanelBg"] = s.Theme.Colors["white"]
		s.Theme.Colors["PanelTextDim"] = s.Theme.Colors["dark-grayii"]
		s.Theme.Colors["PanelBgDim"] = s.Theme.Colors["dark-grayi"]
		s.Theme.Colors["DocTextDim"] = s.Theme.Colors["light-grayi"]
		s.Theme.Colors["DocBgDim"] = s.Theme.Colors["dark-grayi"]
		s.Theme.Colors["Warning"] = s.Theme.Colors["light-orange"]
		s.Theme.Colors["Success"] = s.Theme.Colors["dark-green"]
		s.Theme.Colors["Check"] = s.Theme.Colors["orange"]
		s.Theme.Colors["DocBgHilite"] = s.Theme.Colors["dark-white"]
	} else {
		s.Theme.Colors["DocText"] = s.Theme.Colors["light"]
		s.Theme.Colors["DocBg"] = s.Theme.Colors["black"]
		s.Theme.Colors["PanelText"] = s.Theme.Colors["light"]
		s.Theme.Colors["PanelBg"] = s.Theme.Colors["dark"]
		s.Theme.Colors["PanelTextDim"] = s.Theme.Colors["light-grayii"]
		s.Theme.Colors["PanelBgDim"] = s.Theme.Colors["light-gray"]
		s.Theme.Colors["DocTextDim"] = s.Theme.Colors["light-gray"]
		s.Theme.Colors["DocBgDim"] = s.Theme.Colors["light-grayii"]
		s.Theme.Colors["Warning"] = s.Theme.Colors["yellow"]
		s.Theme.Colors["Success"] = s.Theme.Colors["green"]
		s.Theme.Colors["Check"] = s.Theme.Colors["orange"]
		s.Theme.Colors["DocBgHilite"] = s.Theme.Colors["light-black"]
	}
}
