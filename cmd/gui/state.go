package gui

import (
	"time"
	
	l "gioui.org/layout"
	uberatomic "go.uber.org/atomic"
	
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/atom"
)

// CategoryFilter marks which transactions to omit from the filtered transaction list
type CategoryFilter struct {
	Send     bool
	Generate bool
	Immature bool
	Receive  bool
	Unknown  bool
}

func (c *CategoryFilter) Filter(s string) (include bool) {
	include = true
	if c.Send && s == "send" {
		include = false
	}
	if c.Generate && s == "generate" {
		include = false
	}
	if c.Immature && s == "immature" {
		include = false
	}
	if c.Receive && s == "receive" {
		include = false
	}
	if c.Unknown && s == "unknown" {
		include = false
	}
	return
}

type State struct {
	lastUpdated             *atom.Time
	bestBlockHeight         *atom.Int32
	bestBlockHash           *atom.Hash
	balance                 *atom.Float64
	balanceUnconfirmed      *atom.Float64
	goroutines              []l.Widget
	allTxs                  *atom.ListTransactionsResult
	filteredTxs             *atom.ListTransactionsResult
	filter                  CategoryFilter
	filterChanged           *atom.Bool
	currentReceivingAddress *atom.Address
	activePage              *uberatomic.String
}

func GetNewState() *State {
	fc := &atom.Bool{
		Bool: uberatomic.NewBool(false),
	}
	return &State{
		lastUpdated:     atom.NewTime(time.Now()),
		bestBlockHeight: &atom.Int32{Int32: uberatomic.NewInt32(0)},
		bestBlockHash:   atom.NewHash(chainhash.Hash{}),
		balance:         &atom.Float64{Float64: uberatomic.NewFloat64(0)},
		balanceUnconfirmed: &atom.Float64{Float64: uberatomic.NewFloat64(0),
		},
		goroutines: nil,
		allTxs: atom.NewListTransactionsResult(
			[]btcjson.ListTransactionsResult{}),
		filteredTxs: atom.NewListTransactionsResult(
			[]btcjson.ListTransactionsResult{}),
		filter:                  CategoryFilter{},
		filterChanged:           fc,
		currentReceivingAddress: atom.NewAddress(&util.AddressPubKeyHash{}),
		activePage:              uberatomic.NewString("home"),
	}
}

func (s *State) BumpLastUpdated() {
	s.lastUpdated.Store(time.Now())
}

type Marshalled struct {
	LastUpdated        time.Time
	BestBlockHeight    int32
	BestBlockHash      chainhash.Hash
	Balance            float64
	BalanceUnconfirmed float64
	AllTxs             []btcjson.ListTransactionsResult
	Filter             CategoryFilter
	ActivePage         string
}

func (s *State) Marshal() (out *Marshalled) {
	out = &Marshalled{
		LastUpdated:        s.lastUpdated.Load(),
		BestBlockHeight:    s.bestBlockHeight.Load(),
		BestBlockHash:      s.bestBlockHash.Load(),
		Balance:            s.balance.Load(),
		BalanceUnconfirmed: s.balanceUnconfirmed.Load(),
		AllTxs:             s.allTxs.Load(),
		Filter:             s.filter,
		ActivePage:         s.activePage.Load(),
	}
	return
}

func (m *Marshalled) Unmarshal(s *State) {
	s.lastUpdated.Store(m.LastUpdated)
	s.bestBlockHeight.Store(m.BestBlockHeight)
	s.bestBlockHash.Store(m.BestBlockHash)
	s.balance.Store(m.Balance)
	s.balanceUnconfirmed.Store(m.BalanceUnconfirmed)
	s.allTxs.Store(m.AllTxs)
	s.filter = m.Filter
	s.activePage.Store(m.ActivePage)
	return
}

func (s *State) Goroutines() []l.Widget {
	return s.goroutines
}

func (s *State) SetGoroutines(gr []l.Widget) {
	s.goroutines = gr
}

func (s *State) SetAllTxs(allTxs []btcjson.ListTransactionsResult) {
	s.allTxs.Store(allTxs)
	// generate filtered state
	filteredTxs := make([]btcjson.ListTransactionsResult, 0, len(s.allTxs.Load()))
	atxs := s.allTxs.Load()
	for i := range atxs {
		if s.filter.Filter(atxs[i].Category) {
			filteredTxs = append(filteredTxs, atxs[i])
		}
	}
}

func (s *State) LastUpdated() time.Time {
	return s.lastUpdated.Load()
}

func (s *State) BestBlockHeight() int32 {
	return s.bestBlockHeight.Load()
}

func (s *State) BestBlockHash() *chainhash.Hash {
	o := s.bestBlockHash.Load()
	return &o
}

func (s *State) Balance() float64 {
	return s.balance.Load()
}

func (s *State) BalanceUnconfirmed() float64 {
	return s.balanceUnconfirmed.Load()
}

func (s *State) ActivePage() string {
	return s.activePage.Load()
}

func (s *State) SetActivePage(page string) {
	s.activePage.Store(page)
}

func (s *State) SetBestBlockHeight(height int32) {
	s.BumpLastUpdated()
	s.bestBlockHeight.Store(height)
}

func (s *State) SetBestBlockHash(h *chainhash.Hash) {
	s.BumpLastUpdated()
	s.bestBlockHash.Store(*h)
}

func (s *State) SetBalance(total float64) {
	s.BumpLastUpdated()
	s.balance.Store(total)
}

func (s *State) SetBalanceUnconfirmed(unconfirmed float64) {
	s.BumpLastUpdated()
	s.balanceUnconfirmed.Store(unconfirmed)
}
