package bnd

var files = []string{"svelte.js", "svelte.css"}

type sveBundle map[string][]byte

type DuOSsveBundle map[string][]string

type DuOSasset struct {
	Name        string `json:"name"`
	DataRaw     []byte `json:"dataraw"`
	DataZip     string `json:"datazip"`
	ContentType string `json:"type"`
}

type DuOSassets struct {
	Html   string     `json:"html"`
	Css    DuOScss    `json:"css"`
	Fonts  DuOSfonts  `json:"fonts"`
	Svelte DuOSsvelte `json:"svelte"`
}

type DuOSsvelte struct {
	Js  string `json:"js"`
	Css string `json:"css"`
}

type DuOScss struct {
	Root    string `json:"root"`
	Colors  string `json:"colors"`
	Helpers string `json:"helpers"`
	Grid    string `json:"grid"`
}

type DuOSfonts struct {
	BariolRegular string `json:"bariol"`
	BariolThin    string `json:"bariolthin"`
	BariolBold    string `json:"bariolbold"`
	BariolLight   string `json:"bariollight"`
}
