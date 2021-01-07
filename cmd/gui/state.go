package gui

import (
	"crypto/cipher"
	"encoding/json"
	"io/ioutil"
	"time"
	
	l "gioui.org/layout"
	uberatomic "go.uber.org/atomic"
	
	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/coding/gcm"
	"github.com/p9c/pod/pkg/comm/transport"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/atom"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const ZeroAddress = "1111111111111111111114oLvT2"

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

type AddressEntry struct {
	Address, Comment  string
	Created, Modified time.Time
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
	isAddress               *atom.Bool
	activePage              *uberatomic.String
	addressBook             []AddressEntry
}

func GetNewState(params *netparams.Params,
	activePage *uberatomic.String) *State {
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
		filter:        CategoryFilter{},
		filterChanged: fc,
		currentReceivingAddress: atom.NewAddress(&util.AddressPubKeyHash{},
			params),
		isAddress:  &atom.Bool{Bool: uberatomic.NewBool(false)},
		activePage: activePage,
	}
}

func (s *State) BumpLastUpdated() {
	s.lastUpdated.Store(time.Now())
}

func (s *State) SetReceivingAddress(addr util.Address) {
	s.currentReceivingAddress.Store(addr)
}

func (s *State) IsReceivingAddress() bool {
	addr := s.currentReceivingAddress.String.Load()
	if addr == ZeroAddress {
		s.isAddress.Store(false)
	} else {
		s.isAddress.Store(true)
	}
	return s.isAddress.Load()
}

func (s *State) Save(filename string, pass *string) (err error) {
	Debug("saving state...")
	marshalled := s.Marshal()
	var j []byte
	if j, err = json.MarshalIndent(marshalled, "", "  "); Check(err) {
		return
	}
	// Debug(string(j))
	var ciph cipher.AEAD
	if ciph, err = gcm.GetCipher(*pass); Check(err) {
		return
	}
	var nonce []byte
	if nonce, err = transport.GetNonce(ciph); Check(err) {
		return
	}
	crypted := append(nonce, ciph.Seal(nil, nonce, j, nil)...)
	var b []byte
	_ = b
	if b, err = ciph.Open(nil, nonce, crypted[len(nonce):], nil); Check(err) {
		interrupt.Request()
		return
	}
	if err = ioutil.WriteFile(filename, crypted, 0700); Check(err) {
	}
	if err = ioutil.WriteFile(filename+".clear", j, 0700); Check(err) {
	}
	return
}

func (s *State) Load(filename string, pass *string) (err error) {
	Debug("loading state...")
	var data []byte
	var ciph cipher.AEAD
	if data, err = ioutil.ReadFile(filename); Check(err) {
		return
	}
	// Debug("cipher:", *pass)
	if ciph, err = gcm.GetCipher(*pass); Check(err) {
		return
	}
	ns := ciph.NonceSize()
	// Debug("nonce size:", ns)
	nonce := data[:ns]
	data = data[ns:]
	var b []byte
	if b, err = ciph.Open(nil, nonce, data, nil); Check(err) {
		interrupt.Request()
		return
	}
	// yay, right password, now unmarshal
	ss := &Marshalled{}
	if err = json.Unmarshal(b, ss); Check(err) {
		return
	}
	// Debug(string(b))
	ss.Unmarshal(s)
	return
}

type Marshalled struct {
	LastUpdated        time.Time
	BestBlockHeight    int32
	BestBlockHash      chainhash.Hash
	Balance            float64
	BalanceUnconfirmed float64
	AllTxs             []btcjson.ListTransactionsResult
	Filter             CategoryFilter
	ReceivingAddress   string
	ActivePage         string
	AddressBook        []AddressEntry
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
		ReceivingAddress:   s.currentReceivingAddress.Load().EncodeAddress(),
		ActivePage:         s.activePage.Load(),
		AddressBook:        s.addressBook,
	}
	return
}

func (m *Marshalled) Unmarshal(s *State) {
	s.lastUpdated.Store(m.LastUpdated)
	s.bestBlockHeight.Store(m.BestBlockHeight)
	s.bestBlockHash.Store(m.BestBlockHash)
	s.balance.Store(m.Balance)
	s.balanceUnconfirmed.Store(m.BalanceUnconfirmed)
	if len(s.allTxs.Load()) < len(m.AllTxs) {
		s.allTxs.Store(m.AllTxs)
	}
	s.filter = m.Filter
	
	if m.ReceivingAddress != "1111111111111111111114oLvT2" {
		var err error
		var ra util.Address
		if ra, err = util.DecodeAddress(m.ReceivingAddress, s.currentReceivingAddress.ForNet); Check(err) {
		}
		s.currentReceivingAddress.Store(ra)
	}
	s.SetActivePage(m.ActivePage)
	s.addressBook = m.AddressBook
	return
}

func (s *State) Goroutines() []l.Widget {
	return s.goroutines
}

func (s *State) SetGoroutines(gr []l.Widget) {
	s.goroutines = gr
}

func (s *State) SetAllTxs(atxs []btcjson.ListTransactionsResult) {
	s.allTxs.Store(atxs)
	// generate filtered state
	filteredTxs := make([]btcjson.ListTransactionsResult, 0, len(s.allTxs.Load()))
	for i := range atxs {
		if s.filter.Filter(atxs[i].Category) {
			filteredTxs = append(filteredTxs, atxs[i])
		}
	}
	s.filteredTxs.Store(filteredTxs)
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
