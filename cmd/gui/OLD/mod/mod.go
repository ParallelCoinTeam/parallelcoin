package mod

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/wallet/dap/res"
	"github.com/p9c/pod/pkg/gui/wallet/dap/win"
	"github.com/p9c/pod/pkg/gui/wallet/nav"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

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

type Settings struct {
	Dir  string
	File string
}

type UserInterface struct {
	Device string
	W      *win.Windows

	Theme *theme.Theme
	// Ekran   func(gtx C) D
	FontSize float32
	R        *res.Responsive
	// P        Pages
	N   *nav.Navigation
	F   func(gtx layout.Context) layout.Dimensions
	G   layout.Context
	Pop bool
	Ops op.Ops
}

// type Pages map[string]Page
