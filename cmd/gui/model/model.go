package model

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	log "github.com/p9c/pod/pkg/util/logi"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type DuoUIconsoleHistory struct {
	Commands       []DuoUIconsoleCommand `json:"coms"`
	CommandsNumber int                   `json:"comnumber"`
}

type DuoUIconsoleCommand struct {
	Com      interface{}
	ComID    string
	Category string
	Out      string
	Time     time.Time
}

type DuoUIconsoleCommandsNumber struct {
	CommandsNumber int `json:"comnumber"`
}

// Items

type DuOSitem struct {
	Enabled  bool        `json:"enabled"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Version  string      `json:"ver"`
	CompType string      `json:"comptype"`
	SubType  string      `json:"subtype"`
	Data     interface{} `json:"data"`
}

type DuOSitems struct {
	Slug  string              `json:"slug"`
	Items map[string]DuOSitem `json:"items"`
}

type DuOScomps []DuOScomp

//  Vue App Model
type DuOScomp struct {
	IsApp       bool   `json:"isapp"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Version     string `json:"ver"`
	Description string `json:"desc"`
	State       string `json:"state"`
	Image       string `json:"img"`
	URL         string `json:"url"`
	CompType    string `json:"comtype"`
	SubType     string `json:"subtype"`
	Js          string `json:"js"`
	Template    string `json:"template"`
	Css         string `json:"css"`
}

type DuoUIdialog struct {
	Show        bool
	Green       func()
	GreenLabel  string
	Orange      func()
	OrangeLabel string
	Red         func()
	RedLabel    string
	CustomField func()
	Title       string
	Text        string
}

type DuOScomponent struct {
	Name       string
	Version    string
	Model      interface{}
	View       func()
	Controller func()
}

type DuoUI struct {
	Window     *app.Window
	Context    *layout.Context
	Theme      *gelook.DuoUItheme
	Pages      *DuoUIpages
	Navigation *DuoUInav
	// Configuration *DuoUIconfiguration
	Viewport int
	IsReady  bool
}

type DuoUIpages struct {
	CurrentPage *gelook.DuoUIpage
	Controller  map[string]*gel.DuoUIpage
	Theme       map[string]*gelook.DuoUIpage
}

type DuoUIlog struct {
	Mx          sync.Mutex
	LogMessages atomic.Value // []log.Entry
	LogChan     chan log.Entry
	StopLogger  chan struct{}
}

// type DuoUIconfiguration struct {
//	Abbrevation        string
//	PrimaryTextColor   color.RGBA
//	SecondaryTextColor color.RGBA
//	PrimaryBgColor     color.RGBA
//	SecondaryBgColor   color.RGBA
//	Navigations        map[string]*view.DuoUIthemeNav
// }

type DuoUIconfTabs struct {
	Current  string
	TabsList map[string]*gel.Button
}

// type DuoUIalert struct {
//	Time      time.Time   `json:"time"`
//	Title     string      `json:"title"`
//	Message   interface{} `json:"message"`
//	AlertType string      `json:"type"`
// }

type DuoUIsettings struct {
	Abbrevation string
	Tabs        *DuoUIconfTabs
	Daemon      *DaemonConfig `json:"daemon"`
}

type DaemonConfig struct {
	// Config  *pod.Config `json:"config"`
	Config  map[string]interface{} `json:"config"`
	Schema  pod.Schema             `json:"schema"`
	Widgets map[string]interface{}
}

type DuoUIblock struct {
	Height        int64   `json:"height"`
	BlockHash     string  `json:"hash"`
	PowAlgoID     uint32  `json:"pow"`
	Difficulty    float64 `json:"diff"`
	Amount        float64 `json:"amount"`
	TxNum         int     `json:"txnum"`
	Confirmations int64
	Time          int64 `json:"time"`
	Link          *gel.Button
}

type DuoUItoast struct {
	Title   string
	Message string
}

type DuoUInav struct {
	Items             map[string]*gelook.DuoUIthemeNav
	Width             int
	Height            int
	TextSize          int
	IconSize          int
	PaddingVertical   int
	PaddingHorizontal int
}

type DuoUIexplorer struct {
	Page        *gel.DuoUIcounter
	PerPage     *gel.DuoUIcounter
	Blocks      []DuoUIblock
	SingleBlock btcjson.GetBlockVerboseResult
}

type DuoUIhistory struct {
	TransList  *layout.List
	Category   string
	Categories *DuoUIhistoryCategories
	Page       *gel.DuoUIcounter
	PerPage    *gel.DuoUIcounter
	Txs        *DuoUItransactionsExcerpts
	SingleTx   btcjson.GetTransactionResult
}

type DuoUIhistoryCategories struct {
	AllTxs      *gel.CheckBox
	MintedTxs   *gel.CheckBox
	ImmatureTxs *gel.CheckBox
	SentTxs     *gel.CheckBox
	ReceivedTxs *gel.CheckBox
}

type DuoUInetwork struct {
	PeersList *layout.List
	Page      *gel.DuoUIcounter
	PerPage   *gel.DuoUIcounter
	Peers     []*btcjson.GetPeerInfoResult `json:"peers"`
}

// System Ststus
type DuoUIstatus struct {
	Version    string
	StartTime  time.Time
	CurrentNet string
	Chain      string
	Node       *NodeStatus
	Wallet     *WalletStatus
	Kopach     *KopachStatus
}

type NodeStatus struct {
	// updates on new block and init
	NetHash          atomic.Uint64
	BlockHeight      atomic.Uint64
	BestBlock        atomic.String
	Difficulty       atomic.Float64
	BlockCount       atomic.Uint64
	NetworkLastBlock atomic.Int32
	// update every 5 seconds
	ConnectionCount atomic.Int32
}

type KopachStatus struct {
	// update every second
	Hashrate uint64
	Hps      *ring.BufferFloat64
}

type WalletStatus struct {
	// unchanging
	WalletVersion map[string]btcjson.VersionResult `json:"walletver"`
	// update on new block and at start
	Balance     atomic.String
	Unconfirmed atomic.String
	TxsNumber   atomic.Uint64
	// components
	LastTxs *DuoUItransactionsExcerpts
}

// slots of user interface
type DuoUIhashes struct{ int64 }

type DuoUInetworkHash struct{ int64 }

type DuoUIheight struct{ int32 }

type DuoUIbestBlockHash struct{ string }

type DuoUIdifficulty struct{ float64 }

// type
// MempoolInfo      struct { string}
type DuoUIblockCount struct{ int64 }

type DuoUInetLastBlock struct{ int32 }

type DuoUIconnections struct{ int32 }

type DuoUIlocalHost struct {
	// Cpu        []cpu.InfoStat        `json:"cpu"`
	// CpuPercent []float64             `json:"cpupercent"`
	// Memory     mem.VirtualMemoryStat `json:"mem"`
	// Disk       disk.UsageStat        `json:"disk"`
}

type DuoUIbalance struct {
	Balance string `json:"balance"`
}

type DuoUIunconfirmed struct {
	Unconfirmed string `json:"unconfirmed"`
}

type DuoUItransactionsNumber struct {
	TxsNumber int `json:"txsnumber"`
}

type DuoUItransactionsExcerpts struct {
	ModelTxsListNumber int
	TxsListNumber      int
	Txs                []DuoUItransactionExcerpt `json:"txs"`
	TxsNumber          int                       `json:"txsnumber"`
	Balance            float64                   `json:"balance"`
	BalanceHeight      float64                   `json:"balanceheight"`
}

type DuoUItransactionExcerpt struct {
	Balance       float64 `json:"balance"`
	Amount        float64 `json:"amount"`
	Category      string  `json:"category"`
	Confirmations int64   `json:"confirmations"`
	Time          string  `json:"time"`
	TxID          string  `json:"txid"`
	Comment       string  `json:"comment,omitempty"`
	Link          *gel.Button
}

type DuoUIaddress struct {
	Index   int     `json:"num"`
	Label   string  `json:"label"`
	Account string  `json:"account"`
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
	Copy    *gel.Button
	QrCode  *gel.Button
}

type DuoUIaddressBook struct {
	ShowMiningAddresses bool
	Num                 int            `json:"num"`
	Addresses           []DuoUIaddress `json:"addresses"`
}

type DuoUIqrCode struct {
	AddrQR  paint.ImageOp
	PubAddr string
}
