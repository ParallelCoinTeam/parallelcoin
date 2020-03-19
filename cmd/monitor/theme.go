package monitor

func (st *State) FlipTheme() {
	st.SetTheme(Toggle(&st.Config.DarkTheme))
}

func (st *State) SetTheme(dark bool) {
	if dark {
		st.Theme.Colors["DocText"] = st.Theme.Colors["Dark"]
		st.Theme.Colors["DocBg"] = st.Theme.Colors["Light"]
		st.Theme.Colors["PanelText"] = st.Theme.Colors["Dark"]
		st.Theme.Colors["PanelBg"] = st.Theme.Colors["White"]
		// st.Theme.Colors["Primary"] = st.Theme.Colors["Gray"]
		// st.Theme.Colors["Secondary"] = st.Theme.Colors["White"]
	} else {
		st.Theme.Colors["DocText"] = st.Theme.Colors["Light"]
		st.Theme.Colors["DocBg"] = st.Theme.Colors["Black"]
		st.Theme.Colors["PanelText"] = st.Theme.Colors["Light"]
		st.Theme.Colors["PanelBg"] = st.Theme.Colors["Dark"]
		// st.Theme.Colors["Primary"] = st.Theme.Colors["Dark"]
		// st.Theme.Colors["Secondary"] = st.Theme.Colors["Black"]
	}
}
