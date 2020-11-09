package gui

import (
	"fmt"
	"sync"
	"time"

	l "gioui.org/layout"
	"github.com/kofoworola/godate"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type State struct {
	mutex              sync.Mutex
	lastUpdated        time.Time
	bestBlockHeight    int
	bestBlockHash      *chainhash.Hash
	balance            float64
	balanceUnconfirmed float64
	txs                []tx
	lastTxs            []btcjson.ListTransactionsResult
	lastTimeStrings    []string
	goroutines         []l.Widget
}

type tx struct {
	time       string
	data       btcjson.ListTransactionsResult
	clickTx    *p9.Clickable
	clickBlock *p9.Clickable
	list       *p9.List
}

func (s *State) Goroutines() []l.Widget {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.goroutines
}

func (s *State) SetGoroutines(gr []l.Widget) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.goroutines = gr
}

func (s *State) SetLastTxs(lastTxs []btcjson.ListTransactionsResult) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastTxs = lastTxs
	// if s.lastTimeStrings == nil {
	s.lastTimeStrings = make([]string, len(s.lastTxs))
	// }
	for i := range s.lastTxs {
		s.lastTimeStrings[i] =
			fmt.Sprintf("%v", godate.Now(time.Local).DifferenceForHumans(
				godate.Create(time.Unix(s.lastTxs[i].BlockTime, 0))))
	}
}

func (s *State) Txs() []tx {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.txs
}

func (s *State) SetTxs(txs []tx) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if txs == nil {
		return
	}
	s.txs = txs
}

func (s *State) LastUpdated() time.Time {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastUpdated
}

func (s *State) BestBlockHeight() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.bestBlockHeight
}

func (s *State) BestBlockHash() *chainhash.Hash {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.bestBlockHash
}

func (s *State) Balance() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.balance
}

func (s *State) BalanceUnconfirmed() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.balanceUnconfirmed
}

func (s *State) SetBestBlockHeight(height int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.bestBlockHeight = height
}

func (s *State) SetBestBlockHash(h *chainhash.Hash) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.bestBlockHash = h
}

func (s *State) SetBalance(total float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.balance = total
}

func (s *State) SetBalanceUnconfirmed(unconfirmed float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.balanceUnconfirmed = unconfirmed
}
