package gui

import (
	"sync"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

type State struct {
	mutex              sync.Mutex
	lastUpdated        time.Time
	bestBlockHeight    int
	bestBlockHash      *chainhash.Hash
	balance       float64
	balanceUnconfirmed float64
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
