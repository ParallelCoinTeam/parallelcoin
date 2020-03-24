package monitor

func (s *State) FlipTheme() {
	L.Debug(s.Config.DarkTheme)
	s.Config.DarkTheme=!s.Config.DarkTheme
	L.Debug(s.Config.DarkTheme)
	s.SetTheme(s.Config.DarkTheme)
	s.SaveConfig()
}

func (s *State) SetTheme(dark bool) {
	if dark {
		s.Theme.Colors["DocText"] = s.Theme.Colors["Dark"]
		s.Theme.Colors["DocBg"] = s.Theme.Colors["Light"]
		s.Theme.Colors["PanelText"] = s.Theme.Colors["Dark"]
		s.Theme.Colors["PanelBg"] = s.Theme.Colors["White"]
		// s.Theme.Colors["Primary"] = s.Theme.Colors["Gray"]
		// s.Theme.Colors["Secondary"] = s.Theme.Colors["White"]
	} else {
		s.Theme.Colors["DocText"] = s.Theme.Colors["Light"]
		s.Theme.Colors["DocBg"] = s.Theme.Colors["Black"]
		s.Theme.Colors["PanelText"] = s.Theme.Colors["Light"]
		s.Theme.Colors["PanelBg"] = s.Theme.Colors["Dark"]
		// s.Theme.Colors["Primary"] = s.Theme.Colors["Dark"]
		// s.Theme.Colors["Secondary"] = s.Theme.Colors["Black"]
	}
}
