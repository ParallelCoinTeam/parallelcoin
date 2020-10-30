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
	balanceTotal       float64
	balanceUnconfirmed float64
	balanceConfirmed   float64
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

func (s *State) BalanceTotal() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.balanceTotal
}

func (s *State) BalanceUnconfirmed() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.balanceUnconfirmed
}

func (s *State) BalanceConfirmed() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.balanceConfirmed
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

func (s *State) SetBalanceTotal(total float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.balanceTotal = total
}

func (s *State) SetBalanceUnconfirmed(unconfirmed float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.balanceUnconfirmed = unconfirmed
}

func (s *State) SetBalanceConfirmed(confirmed float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastUpdated = time.Now()
	s.balanceConfirmed = confirmed
}
