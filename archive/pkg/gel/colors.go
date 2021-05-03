package gel

import (
	"fmt"
	"image/color"
	"sync"
)

// Colors is a map of names to hex strings specifying colors
type Colors struct {
	sync.Mutex
	m map[string]string
}

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

// GetNRGBAFromName returns the named color from the map
func (c *Colors) GetNRGBAFromName(co string) color.NRGBA {
	c.Lock()
	defer c.Unlock()
	if col, ok := c.m[co]; ok {
		return HexNRGB(col)
	}
	return color.NRGBA{}
}

// newColors creates the base palette for the theme
func newColors() (c *Colors) {
	c = new(Colors)
	c.Lock()
	defer c.Unlock()
	c.m = map[string]string{
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
		"halfdim":               "33000000",
		"halfbright":            "02ffffff",
	}
	
	c.m["Black"] = c.m["black"]
	c.m["ButtonBg"] = c.m["blue-lite-blue"]
	c.m["ButtonBgDim"] = "ff30809a"
	c.m["ButtonText"] = c.m["White"]
	c.m["ButtonTextDim"] = c.m["light-grayii"]
	c.m["Chk"] = c.m["orange"]
	c.m["Chk"] = c.m["orange"]
	c.m["Danger"] = c.m["red"]
	c.m["Dark"] = c.m["dark"]
	c.m["DarkGray"] = c.m["dark-grayii"]
	c.m["DarkGrayI"] = c.m["dark-grayi"]
	c.m["DarkGrayII"] = c.m["dark-gray"]
	c.m["DarkGrayIII"] = c.m["dark"]
	c.m["DocBg"] = c.m["white"]
	c.m["DocBgDim"] = c.m["light-grayiii"]
	c.m["DocBgHilite"] = c.m["dark-white"]
	c.m["DocText"] = c.m["dark"]
	c.m["DocTextDim"] = c.m["light-grayi"]
	c.m["Fatal"] = "ff880000"
	c.m["Gray"] = c.m["gray"]
	c.m["Hint"] = c.m["light-gray"]
	c.m["Info"] = c.m["blue-lite-blue"]
	c.m["InvText"] = c.m["light"]
	c.m["Light"] = c.m["light"]
	c.m["LightGray"] = c.m["light-grayiii"]
	c.m["LightGrayI"] = c.m["light-grayii"]
	c.m["LightGrayII"] = c.m["light-grayi"]
	c.m["LightGrayIII"] = c.m["light-gray"]
	c.m["PanelBg"] = c.m["light"]
	c.m["PanelBgDim"] = c.m["dark-grayi"]
	c.m["PanelText"] = c.m["dark"]
	c.m["PanelTextDim"] = c.m["light-grayii"]
	c.m["PrimaryLight"] = c.m["light-green-blue"]
	c.m["Primary"] = c.m["green-blue"]
	c.m["PrimaryDim"] = c.m["dark-green-blue"]
	c.m["SecondaryLight"] = c.m["purple"]
	c.m["Secondary"] = c.m["purple"]
	c.m["SecondaryDim"] = c.m["dark-purple"]
	c.m["Success"] = c.m["green"]
	c.m["Transparent"] = c.m["00000000"]
	c.m["Warning"] = c.m["light-orange"]
	c.m["White"] = c.m["white"]
	c.m["scrim"] = c.m["halfbright"]
	
	c.m["Primary"] = c.m["PrimaryLight"]
	c.m["Secondary"] = c.m["SecondaryLight"]
	
	c.m["DocText"] = c.m["dark"]
	c.m["DocBg"] = c.m["light"]
	
	c.m["PanelText"] = c.m["dark"]
	c.m["PanelBg"] = c.m["light"]
	
	c.m["PanelTextDim"] = c.m["dark-grayii"]
	c.m["PanelBgDim"] = c.m["dark-grayi"]
	c.m["DocTextDim"] = c.m["light-grayi"]
	c.m["DocBgDim"] = c.m["dark-grayi"]
	c.m["Warning"] = c.m["light-orange"]
	c.m["Success"] = c.m["dark-green"]
	c.m["Chk"] = c.m["orange"]
	c.m["DocBgHilite"] = c.m["dark-white"]
	c.m["scrim"] = c.m["halfbright"]
	return c
}

// SetDarkTheme to dark or light
func (c *Colors) SetDarkTheme(dark bool) {
	c.Lock()
	defer c.Unlock()
	if !dark {
		c.m["Primary"] = c.m["PrimaryLight"]
		c.m["Secondary"] = c.m["SecondaryLight"]
		
		c.m["DocText"] = c.m["dark"]
		c.m["DocBg"] = c.m["light"]
		
		c.m["PanelText"] = c.m["dark"]
		c.m["PanelBg"] = c.m["white"]
		
		c.m["PanelTextDim"] = c.m["dark-grayii"]
		c.m["PanelBgDim"] = c.m["dark-grayi"]
		c.m["DocTextDim"] = c.m["light-grayi"]
		c.m["DocBgDim"] = c.m["light-grayiii"]
		c.m["Warning"] = c.m["light-orange"]
		c.m["Success"] = c.m["dark-green"]
		c.m["Chk"] = c.m["orange"]
		c.m["DocBgHilite"] = c.m["dark-white"]
		c.m["scrim"] = c.m["halfdim"]
		c.m["Fatal"] = c.m["light-red"]
		c.m["Info"] = c.m["light-blue"]
	} else {
		c.m["Primary"] = c.m["PrimaryDim"]
		c.m["Secondary"] = c.m["SecondaryDim"]
		
		c.m["DocText"] = c.m["light"]
		c.m["DocBg"] = c.m["dark"]
		
		c.m["PanelText"] = c.m["light"]
		c.m["PanelBg"] = c.m["black"]
		
		c.m["PanelTextDim"] = c.m["light-grayii"]
		c.m["PanelBgDim"] = c.m["light-gray"]
		c.m["DocTextDim"] = c.m["dark-gray"]
		c.m["DocBgDim"] = c.m["dark-grayii"]
		c.m["Warning"] = c.m["yellow"]
		c.m["Success"] = c.m["green"]
		c.m["Chk"] = c.m["orange"]
		c.m["DocBgHilite"] = c.m["light-black"]
		c.m["scrim"] = c.m["halfbright"]
		c.m["Fatal"] = c.m["red"]
		c.m["Info"] = c.m["blue"]
	}
}
