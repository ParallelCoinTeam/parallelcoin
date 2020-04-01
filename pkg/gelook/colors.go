// SPDX-License-Identifier: Unlicense OR MIT

package gelook

func NewDuoUIcolors() (c map[string]string) {
	c = make(map[string]string)
	c["Black"] = "ff000000"
	c["White"] = "ffffffff"
	c["Gray"] = "ff808080"
	c["Light"] = "ffcfcfcf"
	c["LightGray"] = "ffbdbdbd"
	c["LightGrayI"] = "ffacacac"
	c["LightGrayII"] = "ff9a9a9a"
	c["LightGrayIII"] = "ff888888"
	c["Dark"] = "ff303030"
	c["DarkGray"] = "ff424242"
	c["DarkGrayI"] = "ff535353"
	c["DarkGrayII"] = "ff656565"
	c["DarkGrayIII"] = "ff303030"
	c["Primary"] = "ff308080"
	c["Secondary"] = "ff803080"
	c["Success"] = "ff30cf30"
	c["Danger"] = "ffcf3030"
	c["Warning"] = "ffcfcf30"
	c["Info"] = "ff3080cf"
	c["Check"] = "ffcf8030"
	c["Hint"] = "ff888888"
	c["InvText"] = "ffcfcfcf"
	c["ButtonText"] = "ffcfcfcf"
	c["ButtonBg"] = "ff3080cf"
	c["PanelText"] = "ffcfcfcf"
	c["PanelBg"] = c["Dark"]
	c["DocText"] = c["Dark"]
	c["DocBg"] = c["Light"]
	c["ButtonTextDim"] = c["LightGrayI"]
	c["ButtonBgDim"] = "ff30809a"
	c["PanelTextDim"] = c["LightGrayI"]
	c["PanelBgDim"] = c["LightGrayII"]
	c["DocTextDim"] = c["LightGrayII"]
	c["DocBgDim"] = c["LightGrayI"]
	c["Transparent"] = c["00000000"]
	return c
}
