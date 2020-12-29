package gui

import (
	"fmt"
	"sync/atomic"
	"time"
	
	uberatomic "go.uber.org/atomic"
	
	l "gioui.org/layout"
	
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
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
	// mutex                   sync.Mutex
	lastUpdated        time.Time
	bestBlockHeight    int
	bestBlockHash      *chainhash.Hash
	balance            uberatomic.Float64
	balanceUnconfirmed uberatomic.Float64
	txs                []tx
	goroutines         []l.Widget
	// AllMutex, FilteredMutex sync.Mutex
	AllTxs                  []btcjson.ListTransactionsResult
	AllTimeStrings          atomic.Value
	FilteredTxs             []btcjson.ListTransactionsResult
	FilteredTimeStrings     []string
	Filter                  CategoryFilter
	FilterChanged           bool
	CurrentReceivingAddress util.Address
}

type Marshalled struct {
	LastUpdated        time.Time
	BestBlockHeight    int
	BestBlockHash      chainhash.Hash
	Balance            float64
	BalanceUnconfirmed float64
	AllTxs             []btcjson.ListTransactionsResult
	AllTimeStrings     []string
	Filter             CategoryFilter
}

func (s *State) Marshal() (out *Marshalled) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	out = &Marshalled{
		LastUpdated:        s.lastUpdated,
		BestBlockHeight:    s.bestBlockHeight,
		BestBlockHash:      *s.bestBlockHash,
		Balance:            s.balance.Load(),
		BalanceUnconfirmed: s.balanceUnconfirmed.Load(),
		AllTxs:             s.AllTxs,
		AllTimeStrings:     s.AllTimeStrings.Load().([]string),
		Filter:             s.Filter,
	}
	return
}

func (m *Marshalled) Unmarshal(s *State) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	s.lastUpdated = m.LastUpdated
	s.bestBlockHeight = m.BestBlockHeight
	s.bestBlockHash = &m.BestBlockHash
	s.balance.Store(m.Balance)
	s.balanceUnconfirmed.Store(m.BalanceUnconfirmed)
	s.AllTxs = m.AllTxs
	s.AllTimeStrings.Store(m.AllTimeStrings)
	s.Filter = m.Filter
	return
}

type tx struct {
	time       string
	data       btcjson.ListTransactionsResult
	clickTx    *p9.Clickable
	clickBlock *p9.Clickable
	list       *p9.List
}

func (s *State) Goroutines() []l.Widget {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	return s.goroutines
}

func (s *State) SetGoroutines(gr []l.Widget) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	s.goroutines = gr
}

func (s *State) SetAllTxs(allTxs []btcjson.ListTransactionsResult) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	// s.AllMutex.Lock()
	// defer s.AllMutex.Unlock()
	s.AllTxs = allTxs
	// if s.AllTimeStrings == nil {
	s.AllTimeStrings.Store(make([]string, len(s.AllTxs)))
	// }
	for i := range s.AllTxs {
		s.AllTimeStrings.Load().([]string)[i] =
		// fmt.Sprintf("%v", godate.Now(time.Local).DifferenceForHumans(
		// 	godate.Create(time.Unix(s.AllTxs[i].BlockTime, 0))))
			fmt.Sprintf(
				"%ds ago",
				time.Now().Unix()-s.AllTxs[i].BlockTime,
			)
	}
	// generate filtered state
	// s.FilteredMutex.Lock()
	// defer s.FilteredMutex.Unlock()
	s.FilteredTxs = make([]btcjson.ListTransactionsResult, 0, len(s.AllTxs))
	s.FilteredTimeStrings = make([]string, 0, len(s.AllTxs))
	for i := range s.AllTxs {
		if s.Filter.Filter(s.AllTxs[i].Category) {
			s.FilteredTxs = append(s.FilteredTxs, s.AllTxs[i])
			s.FilteredTimeStrings = append(s.FilteredTimeStrings, s.AllTimeStrings.Load().([]string)[i])
		}
	}
	// // reverse the filtered tx's because they are in reverse chronological order and prepend rather than append, and
	// // the history view needs them to be immutable but they can grow
	// lf := len(s.AllTxs) - 1
	// if lf > 0 {
	// 	for i := 0; i < lf/2; i++ {
	// 		s.FilteredTxs[i], s.FilteredTxs[lf-i] = s.FilteredTxs[lf-i], s.FilteredTxs[i]
	// 	}
	// }
}

func (s *State) Txs() []tx {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	return s.txs
}

func (s *State) SetTxs(txs []tx) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	if txs == nil {
		return
	}
	s.txs = txs
}

func (s *State) LastUpdated() time.Time {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	return s.lastUpdated
}

func (s *State) BestBlockHeight() int {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	return s.bestBlockHeight
}

func (s *State) BestBlockHash() *chainhash.Hash {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	return s.bestBlockHash
}

func (s *State) Balance() float64 {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	return s.balance.Load()
}

func (s *State) BalanceUnconfirmed() float64 {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	return s.balanceUnconfirmed.Load()
}

func (s *State) SetBestBlockHeight(height int) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.bestBlockHeight = height
}

func (s *State) SetBestBlockHash(h *chainhash.Hash) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.bestBlockHash = h
}

func (s *State) SetBalance(total float64) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.balance.Store(total)
}

func (s *State) SetBalanceUnconfirmed(unconfirmed float64) {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.balanceUnconfirmed.Store(unconfirmed)
}
