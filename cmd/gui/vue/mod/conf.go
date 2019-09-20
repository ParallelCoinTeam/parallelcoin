package mod

type DisplayConfig struct {
	Multi      bool                   `json:"multi"`
	Options    map[string]interface{} `json:"options`
	Responsive bool                   `json:"responsive"`
	Screens    map[string]EJSlayout   `json:"screens"`
}
