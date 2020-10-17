package gwallet

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/p9c/pod/pkg/gui/wallet/dap/mod"
)

var (
	selected int
	pwd      []string
)

type (
	D = layout.Dimensions
	C = layout.Context
	W = layout.Widget
)

type GioWallet struct {
	//rpc    *rpcclient.Client
	Status Status
	ui     *mod.UserInterface
}

type folderListItem struct {
	Name  string
	Hash  string
	Size  uint64
	Type  uint8
	btn   *widget.Clickable
	check *widget.Bool
}

type Status struct {
	bal *Balances
	txs []Tx
}

type Balances struct {
	available string
	pending   string
	immature  string
	total     string
}
type Tx struct {
	Id            string
	Time          string
	Type          string
	Address       string
	Amount        string
	Verifications int
	Btn           *widget.Clickable
}
