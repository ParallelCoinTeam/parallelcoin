// SPDX-License-Identifier: Unlicense OR MIT

package gelook

func NewDuoUIcolors() (c map[string]string) {
	c = map[string]string{
		"black":                 "ff000000",
		"light-black":           "ff222222",
		"blue":                  "ff3030cf",
		"blue-lite-blue":        "ff3080cf",
		"blue-orange":           "ff80a830",
		"blue-red":              "ff803080",
		"dark":                  "ff303030",
		"dark-blue":             "ff303080",
		"dark-blue-lite-blue":   "ff305880",
		"dark-blue-orange":      "ff584458",
		"dark-blue-red":         "ff583058",
		"dark-gray":             "ff656565",
		"dark-grayi":            "ff535353",
		"dark-grayii":           "ff424242",
		"dark-green":            "ff308030",
		"dark-green-blue":       "ff305858",
		"dark-green-lite-blue":  "ff308058",
		"dark-green-orange":     "ff586c30",
		"dark-green-red":        "ff585830",
		"dark-green-yellow":     "ff588030",
		"dark-lite-blue":        "ff308080",
		"dark-orange":           "ff805830",
		"dark-purple":           "ff803080",
		"dark-red":              "ff803030",
		"dark-yellow":           "ff808030",
		"gray":                  "ff808080",
		"green":                 "ff30cf30",
		"green-blue":            "ff308080",
		"green-lite-blue":       "ff30cf80",
		"green-orange":          "ff80a830",
		"green-red":             "ff808030",
		"green-yellow":          "ff80cf30",
		"light":                 "ffcfcfcf",
		"light-blue":            "ff8080cf",
		"light-blue-lite-blue":  "ff80a8cf",
		"light-blue-orange":     "ffa894a8",
		"light-blue-red":        "ffa880a8",
		"light-gray":            "ff888888",
		"light-grayi":           "ff9a9a9a",
		"light-grayii":          "ffacacac",
		"light-grayiii":         "ffbdbdbd",
		"light-green":           "ff80cf80",
		"light-green-blue":      "ff80a8a8",
		"light-green-lite-blue": "ff80cfa8",
		"light-green-orange":    "ffa8bc80",
		"light-green-red":       "ffa8a880",
		"light-green-yellow":    "ffa8cf80",
		"light-lite-blue":       "ff80cfcf",
		"light-orange":          "ffcfa880",
		"light-purple":          "ffcf80cf",
		"light-red":             "ffcf8080",
		"light-yellow":          "ffcfcf80",
		"lite-blue":             "ff30cfcf",
		"orange":                "ffcf8030",
		"purple":                "ffcf30cf",
		"red":                   "ffcf3030",
		"white":                 "ffffffff",
		"dark-white":            "ffdddddd",
		"yellow":                "ffcfcf30",
	}

	c["Black"] = c["black"]
	c["White"] = c["white"]
	c["Gray"] = c["gray"]
	c["Light"] = c["light"]
	c["LightGray"] = c["light-grayiii"]
	c["LightGrayI"] = c["light-grayii"]
	c["LightGrayII"] = c["light-grayi"]
	c["LightGrayIII"] = c["light-gray"]
	c["Dark"] = c["dark"]
	c["DarkGray"] = c["dark-grayii"]
	c["DarkGrayI"] = c["dark-grayi"]
	c["DarkGrayII"] = c["dark-gray"]
	c["DarkGrayIII"] = c["dark"]
	c["Primary"] = c["green-blue"]
	c["Secondary"] = c["dark-purple"]
	c["Success"] = c["green"]
	c["Danger"] = c["red"]
	c["Warning"] = c["yellow"]
	c["Info"] = c["blue-lite-blue"]
	c["Check"] = c["orange"]
	c["Hint"] = c["light-gray"]
	c["InvText"] = c["light"]
	c["ButtonText"] = c["light"]
	c["ButtonBg"] = c["blue-lite-blue"]
	c["PanelText"] = c["light"]
	c["PanelBg"] = c["dark"]
	c["DocText"] = c["dark"]
	c["DocBg"] = c["light"]
	c["ButtonTextDim"] = c["light-grayii"]
	c["ButtonBgDim"] = "ff30809a"
	c["PanelTextDim"] = c["LightGrayI"]
	c["PanelBgDim"] = c["light-grayi"]
	c["DocTextDim"] = c["light-grayi"]
	c["DocBgDim"] = c["light-grayii"]
	c["Transparent"] = c["00000000"]
	c["Fatal"] = "ff880000"
	return c
}
