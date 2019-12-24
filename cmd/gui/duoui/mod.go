package duoui



type Boot struct {
	IsBoot     bool `json:"boot"`
	IsFirstRun bool `json:"firstrun"`
	IsBootMenu bool `json:"menu"`
	IsBootLogo bool `json:"logo"`
	IsLoading  bool `json:"loading"`
}
