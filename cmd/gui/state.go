package gui

import (
	"time"
	
	l "gioui.org/layout"
	
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
	lastUpdated             atom.Time
	bestBlockHeight         atom.Int32
	bestBlockHash           *chainhash.Hash
	balance                 atom.Float64
	balanceUnconfirmed      atom.Float64
	goroutines              []l.Widget
	AllTxs                  []btcjson.ListTransactionsResult
	FilteredTxs             []btcjson.ListTransactionsResult
	Filter                  CategoryFilter
	FilterChanged           bool
	CurrentReceivingAddress util.Address
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
}

func (s *State) Marshal() (out *Marshalled) {
	out = &Marshalled{
		LastUpdated:        s.lastUpdated.Load(),
		BestBlockHeight:    s.bestBlockHeight.Load(),
		BestBlockHash:      *s.bestBlockHash,
		Balance:            s.balance.Load(),
		BalanceUnconfirmed: s.balanceUnconfirmed.Load(),
		AllTxs:             s.AllTxs,
		Filter:             s.Filter,
	}
	return
}

func (m *Marshalled) Unmarshal(s *State) {
	s.lastUpdated.Store(m.LastUpdated)
	s.bestBlockHeight.Store(m.BestBlockHeight)
	*s.bestBlockHash = m.BestBlockHash
	s.balance.Store(m.Balance)
	s.balanceUnconfirmed.Store(m.BalanceUnconfirmed)
	s.AllTxs = m.AllTxs
	s.Filter = m.Filter
	return
}

func (s *State) Goroutines() []l.Widget {
	return s.goroutines
}

func (s *State) SetGoroutines(gr []l.Widget) {
	s.goroutines = gr
}

func (s *State) SetAllTxs(allTxs []btcjson.ListTransactionsResult) {
	s.AllTxs = allTxs
	// generate filtered state
	s.FilteredTxs = make([]btcjson.ListTransactionsResult, 0, len(s.AllTxs))
	for i := range s.AllTxs {
		if s.Filter.Filter(s.AllTxs[i].Category) {
			s.FilteredTxs = append(s.FilteredTxs, s.AllTxs[i])
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
	return s.bestBlockHash
}

func (s *State) Balance() float64 {
	return s.balance.Load()
}

func (s *State) BalanceUnconfirmed() float64 {
	return s.balanceUnconfirmed.Load()
}

func (s *State) SetBestBlockHeight(height int32) {
	s.BumpLastUpdated()
	s.bestBlockHeight.Store(height)
}

func (s *State) SetBestBlockHash(h *chainhash.Hash) {
	s.BumpLastUpdated()
	s.bestBlockHash = h
}

func (s *State) SetBalance(total float64) {
	s.BumpLastUpdated()
	s.balance.Store(total)
}

func (s *State) SetBalanceUnconfirmed(unconfirmed float64) {
	s.BumpLastUpdated()
	s.balanceUnconfirmed.Store(unconfirmed)
}
