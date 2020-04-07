package gui

func (s *State) FlipTheme(isDark *bool, hook func()) {
	*isDark = !*isDark
	//Debug(s.Config.DarkTheme)
	s.SetTheme(*isDark)
	hook()
}

func (s *State) SetTheme(dark bool) {
	if dark {
		darkTheme(&s.Theme.Colors)
	} else {
		lightTheme(&s.Theme.Colors)
	}
}

func darkTheme(colors *map[string]string) {
	c := *colors
	c["DocText"] = c["dark"]
	c["DocBg"] = c["light"]
	c["PanelText"] = c["dark"]
	c["PanelBg"] = c["white"]
	c["PanelTextDim"] = c["dark-grayii"]
	c["PanelBgDim"] = c["dark-grayi"]
	c["DocTextDim"] = c["light-grayi"]
	c["DocBgDim"] = c["dark-grayi"]
	c["Warning"] = c["light-orange"]
	c["Success"] = c["dark-green"]
	c["Check"] = c["orange"]
	c["DocBgHilite"] = c["dark-white"]

}

func lightTheme(colors *map[string]string) {
	c := *colors
	c["DocText"] = c["light"]
	c["DocBg"] = c["black"]
	c["PanelText"] = c["light"]
	c["PanelBg"] = c["dark"]
	c["PanelTextDim"] = c["light-grayii"]
	c["PanelBgDim"] = c["light-gray"]
	c["DocTextDim"] = c["light-gray"]
	c["DocBgDim"] = c["light-grayii"]
	c["Warning"] = c["yellow"]
	c["Success"] = c["green"]
	c["Check"] = c["orange"]
	c["DocBgHilite"] = c["light-black"]
}
