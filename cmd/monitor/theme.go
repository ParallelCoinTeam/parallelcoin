package monitor

func (s *State) FlipTheme() {
	Debug(s.Config.DarkTheme)
	s.Config.DarkTheme=!s.Config.DarkTheme
	//Debug(s.Config.DarkTheme)
	s.SetTheme(s.Config.DarkTheme)
	s.SaveConfig()
}

func (s *State) SetTheme(dark bool) {
	if dark {
		s.Theme.Colors["DocText"] = s.Theme.Colors["Dark"]
		s.Theme.Colors["DocBg"] = s.Theme.Colors["Light"]
		s.Theme.Colors["PanelText"] = s.Theme.Colors["Dark"]
		s.Theme.Colors["PanelBg"] = s.Theme.Colors["White"]
		s.Theme.Colors["PanelTextDim"] = s.Theme.Colors["DarkGray"]
		s.Theme.Colors["PanelBgDim"] = s.Theme.Colors["DarkGrayI"]
		s.Theme.Colors["DocTextDim"] = s.Theme.Colors["LightGrayII"]
		s.Theme.Colors["DocBgDim"] = s.Theme.Colors["DarkGrayI"]
	} else {
		s.Theme.Colors["DocText"] = s.Theme.Colors["Light"]
		s.Theme.Colors["DocBg"] = s.Theme.Colors["Black"]
		s.Theme.Colors["PanelText"] = s.Theme.Colors["Light"]
		s.Theme.Colors["PanelBg"] = s.Theme.Colors["Dark"]
		s.Theme.Colors["PanelTextDim"] = s.Theme.Colors["LightGrayI"]
		s.Theme.Colors["PanelBgDim"] = s.Theme.Colors["LightGrayIII"]
		s.Theme.Colors["DocTextDim"] = s.Theme.Colors["LightGrayIII"]
		s.Theme.Colors["DocBgDim"] = s.Theme.Colors["LightGrayI"]
	}
}
