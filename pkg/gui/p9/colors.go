package p9

import (
	"fmt"
	"image/color"
)

// Colors is a map of names to hex strings specifying colors
type Colors map[string]string

// HexARGB converts a 32 bit hex string into a color specification
func HexARGB(s string) (c color.RGBA) {
	_, _ = fmt.Sscanf(s, "%02x%02x%02x%02x", &c.A, &c.R, &c.G, &c.B)
	return
}

// HexNRGB converts a 32 bit hex string into a color specification
func HexNRGB(s string) (c color.NRGBA) {
	_, _ = fmt.Sscanf(s, "%02x%02x%02x%02x", &c.A, &c.R, &c.G, &c.B)
	return
}

// Get returns the named color from the map
func (c Colors) Get(co string) color.NRGBA {
	if col, ok := c[co]; ok {
		return HexNRGB(col)
	}
	return color.NRGBA{}
}

// NewColors creates the base palette for the theme
func NewColors() (c Colors) {
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
		"halfdim":               "88000000",
		"halfbright":            "88888888",
	}

	c["Black"] = c["black"]
	c["ButtonBg"] = c["blue-lite-blue"]
	c["ButtonBgDim"] = "ff30809a"
	c["ButtonText"] = c["White"]
	c["ButtonTextDim"] = c["light-grayii"]
	c["Check"] = c["orange"]
	c["Check"] = c["orange"]
	c["Danger"] = c["red"]
	c["Dark"] = c["dark"]
	c["DarkGray"] = c["dark-grayii"]
	c["DarkGrayI"] = c["dark-grayi"]
	c["DarkGrayII"] = c["dark-gray"]
	c["DarkGrayIII"] = c["dark"]
	c["DocBg"] = c["white"]
	c["DocBgDim"] = c["light-grayii"]
	c["DocBgHilite"] = c["dark-white"]
	c["DocText"] = c["dark"]
	c["DocTextDim"] = c["light-grayi"]
	c["Fatal"] = "ff880000"
	c["Gray"] = c["gray"]
	c["Hint"] = c["light-gray"]
	c["Info"] = c["blue-lite-blue"]
	c["InvText"] = c["light"]
	c["Light"] = c["light"]
	c["LightGray"] = c["light-grayiii"]
	c["LightGrayI"] = c["light-grayii"]
	c["LightGrayII"] = c["light-grayi"]
	c["LightGrayIII"] = c["light-gray"]
	c["PanelBg"] = c["light"]
	c["PanelBgDim"] = c["dark-grayi"]
	c["PanelText"] = c["dark"]
	c["PanelTextDim"] = c["light-grayii"]
	c["PrimaryLight"] = c["green-blue"]
	c["Primary"] = c["PrimaryLight"]
	c["PrimaryDim"] = c["dark-green-blue"]
	c["SecondaryLight"] = c["purple"]
	c["Secondary"] = c["purple"]
	c["SecondaryDim"] = c["dark-purple"]
	c["Success"] = c["green"]
	c["Transparent"] = c["00000000"]
	c["Warning"] = c["light-orange"]
	c["White"] = c["white"]
	c["scrim"] = c["halfdim"]

	c["Primary"] = c["PrimaryLight"]
	c["Secondary"] = c["SecondaryLight"]

	c["DocText"] = c["dark"]
	c["DocBg"] = c["white"]

	c["PanelText"] = c["dark"]
	c["PanelBg"] = c["light"]

	c["PanelTextDim"] = c["dark-grayii"]
	c["PanelBgDim"] = c["dark-grayi"]
	c["DocTextDim"] = c["light-grayi"]
	c["DocBgDim"] = c["dark-grayi"]
	c["Warning"] = c["light-orange"]
	c["Success"] = c["dark-green"]
	c["Check"] = c["orange"]
	c["DocBgHilite"] = c["dark-white"]
	c["scrim"] = c["halfbright"]
	return c
}

func (c Colors) SetTheme(dark bool) {
	if !dark {
		c["Primary"] = c["PrimaryLight"]
		c["Secondary"] = c["SecondaryLight"]

		c["DocText"] = c["dark"]
		c["DocBg"] = c["white"]

		c["PanelText"] = c["dark"]
		c["PanelBg"] = c["light"]

		c["PanelTextDim"] = c["dark-grayii"]
		c["PanelBgDim"] = c["dark-grayi"]
		c["DocTextDim"] = c["light-grayi"]
		c["DocBgDim"] = c["dark-grayi"]
		c["Warning"] = c["light-orange"]
		c["Success"] = c["dark-green"]
		c["Check"] = c["orange"]
		c["DocBgHilite"] = c["dark-white"]
		c["scrim"] = c["halfdim"]
	} else {
		c["Primary"] = c["PrimaryDim"]
		c["Secondary"] = c["SecondaryDim"]

		c["DocText"] = c["light"]
		c["DocBg"] = c["dark"]

		c["PanelText"] = c["light"]
		c["PanelBg"] = c["black"]

		c["PanelTextDim"] = c["light-grayii"]
		c["PanelBgDim"] = c["light-gray"]
		c["DocTextDim"] = c["light-gray"]
		c["DocBgDim"] = c["light-grayii"]
		c["Warning"] = c["yellow"]
		c["Success"] = c["green"]
		c["Check"] = c["orange"]
		c["DocBgHilite"] = c["light-black"]
		c["scrim"] = c["halfbright"]
	}
}
