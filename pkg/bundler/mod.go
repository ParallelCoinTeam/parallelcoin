package bnd

type DuOSasset struct {
	Name        string `json:"name"`
	Sub         string `json:"sub"`
	Data        string `json:"data"`
	ContentType string `json:"type"`
}

type DuOSassets map[string]DuOSasset
