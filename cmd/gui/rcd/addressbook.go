package rcd

type DuoUItemplates struct {
	App  map[string][]byte            `json:"app"`
	Data map[string]map[string][]byte `json:"data"`
}

type DbAddress string

type Address struct {
	Index   int     `json:"num"`
	Label   string  `json:"label"`
	Account string  `json:"account"`
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}
type Send struct {
	// Phrase string  `json:"phrase"`
	// Addr   string  `json:"addr"`
	// Amount float64 `json:"amount"`
	//exit
}

type AddBook struct {
	Address string `json:"address"`
	Label   string `json:"label"`
}
