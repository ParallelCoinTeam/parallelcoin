package bnd

type sveBundle map[string][]byte

type DuOSasset struct {
	Name        string `json:"name"`
	Sub         string `json:"sub"`
	DataRaw     []byte `json:"dataraw"`
	DataZip     string `json:"datazip"`
	ContentType string `json:"type"`
}

type DuOSassets map[string]DuOSasset
