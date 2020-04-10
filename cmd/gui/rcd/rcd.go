package rcd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	scribble "github.com/nanobox-io/golang-scribble"
	config2 "github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	blockchain "github.com/p9c/pod/pkg/chain"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/wallet"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	"golang.org/x/text/unicode/norm"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	NewBlock uint32 = iota
	// Add new events here
	EventCount
)

// System Ststus
type (
	AddBook struct {
		Address string `json:"address"`
		Label   string `json:"label"`
	}
	Address struct {
		Index   int     `json:"num"`
		Label   string  `json:"label"`
		Account string  `json:"account"`
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}
	AddressSlice []model.DuoUIaddress
	Boot         struct {
		IsBoot     bool   `json:"boot"`
		IsFirstRun bool   `json:"firstrun"`
		IsBootMenu bool   `json:"menu"`
		IsBootLogo bool   `json:"logo"`
		IsLoading  bool   `json:"loading"`
		IsScreen   string `json:"screen"`
	}
	// type rcVar interface {
	//	GetDuoUItransactions(sfrom, count int, cat string)
	//	GetDuoUIbalance()
	//	GetDuoUItransactionsExcerpts()
	//	DuoSend(wp string, ad string, am float64)
	//	GetDuoUItatus()
	//	PushDuoUIalert(t string, m interface{}, at string)
	//	GetDuoUIblockHeight()
	//	GetDuoUIblockCount()
	//	GetDuoUIconnectionCount()
	// }
	DbAddress string
	Ddb       interface {
		DbReadAllTypes()
		DbRead(folder, name string)
		DbReadAll(folder string) DuoUIitems
		DbWrite(folder, name string, data interface{})
	}
	DuoUIcommand struct {
		Command func()
	}
	CommandEvent struct {
		Command DuoUIcommand
	}
	DuoUIcommands struct {
		Events  chan DuoUIcommandsEvent
		History []DuoUIcommand
	}
	DuoUIcommandsEvent interface {
		isCommandsEvent()
	}
	DB struct {
		DB     *scribble.Driver
		Folder string      `json:"folder"`
		Name   string      `json:"name"`
		Data   interface{} `json:"data"`
	}
	DuoUIitem struct {
		Enabled  bool        `json:"enabled"`
		Name     string      `json:"name"`
		Slug     string      `json:"slug"`
		Version  string      `json:"ver"`
		CompType string      `json:"comptype"`
		SubType  string      `json:"subtype"`
		Data     interface{} `json:"data"`
	}
	DuoUIitems struct {
		Slug  string               `json:"slug"`
		Items map[string]DuoUIitem `json:"items"`
	}
	DuoUItemplates struct {
		App  map[string][]byte            `json:"app"`
		Data map[string]map[string][]byte `json:"data"`
	}
	ErrorEvent struct {
		Err error
	}
	Event struct {
		Type    uint32
		Payload []byte
	}
	RcVar struct {
		cx             *conte.Xt
		db             *DB
		Boot           *Boot
		Events         chan Event
		UpdateTrigger  chan struct{}
		Status         *model.DuoUIstatus
		Dialog         *model.DuoUIdialog
		Log            *model.DuoUIlog
		ConsoleHistory *model.DuoUIconsoleHistory
		Commands       *DuoUIcommands
		Settings       *model.DuoUIsettings
		Sent           bool
		Toasts         []model.DuoUItoast
		Localhost      model.DuoUIlocalHost
		Uptime         int
		AddressBook    *model.DuoUIaddressBook
		ShowPage       string
		CurrentPage    *gelook.DuoUIpage
		// NodeChan   chan *rpc.Server
		// WalletChan chan *wallet.Wallet
		Explorer *model.DuoUIexplorer
		History  *model.DuoUIhistory
		Network  *model.DuoUInetwork
		Quit     chan struct{}
		Ready    chan struct{}
		IsReady  bool
	}
)

var (
	// MaxLogLength is a var so it can be changed dynamically
	MaxLogLength     = 16384
	_            Ddb = &DB{}
	EventsChan       = make(chan Event, 1)
	safe             = []*unicode.RangeTable{
		unicode.Letter,
		unicode.Number,
	}
	skip = []*unicode.RangeTable{
		unicode.Mark,
		unicode.Sk,
		unicode.Lm,
	}
)

// sorter implementation for AddressSlice
func (a AddressSlice) Len() int           { return len(a) }
func (a AddressSlice) Less(i, j int) bool { return a[i].Index < a[j].Index }
func (a AddressSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (d *DuoUIcommands) Run() {
	d.Events <- CommandEvent{
		Command: DuoUIcommand{},
	}
}

func (d *DB) DbRead(folder, name string) {
	item := DuoUIitem{}
	if err := d.DB.Read(folder, name, &item); !Check(err) {
		d.Data = item
	}
}

func (d *DB) DbReadAddressBook() (addressbook map[string]string) {
	addressbook = make(map[string]string)
	if err := d.DB.Read("user", "addressbook", &addressbook); !Check(err) {
		return addressbook
	}
	return
}

func (d *DB) DbReadAll(folder string) DuoUIitems {
	items := make(map[string]DuoUIitem)
	if itemsRaw, err := d.DB.ReadAll(folder); Check(err) {
		if err != nil {
			fmt.Println("Error", err)
		}
		for _, bt := range itemsRaw {
			item := DuoUIitem{}
			if err := json.Unmarshal([]byte(bt), &item); !Check(err) {
				items[item.Slug] = item
			}
		}
	}
	return DuoUIitems{
		Slug:  folder,
		Items: items,
	}

}

func (d *DB) DbReadAllTypes() {
	items := make(map[string]DuoUIitems)
	types := []string{"assets", "config", "apps"}
	for _, t := range types {
		items[t] = d.DbReadAll(t)
	}
	d.Data = items
	Debug("ooooooooooooooooooooooooooooodaaa", d.Data)

}

func (d *DB) DbReadTypeAll(f string) {
	d.Data = d.DbReadAll(f)
}

func (d *DB) DbWrite(folder, name string, data interface{}) {
	d.DB.Write(folder, name, data)
}

func (d *DB) DuoUIdbInit(dataDir string) {
	db, err := scribble.New(dataDir+"/gui", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	d.DB = db
}

func (e CommandEvent) isCommandsEvent() {}

func (e ErrorEvent) isCommandsEvent() {}

func (r *RcVar) ConsoleCmd(com string) (out chan string) {
	out = make(chan string)
	go func() {
		var o string
		split := strings.Split(com, " ")
		method, args := split[0], split[1:]
		ws, cc := r.cx.WalletServer, r.cx.ChainClient
		rpcSrv, lrpcHnd := r.cx.RPCServer, legacy.RPCHandlers
		var cmd, res interface{}
		var err error
		var errString, prev string
		if method == "help" {
			if len(args) < 1 {
				method = ""
				cmd = &btcjson.HelpCmd{Command: &method}
				if res, err = chainrpc.RPCHandlers["help"].Fn(rpcSrv, cmd, nil); Check(err) {
					errString += fmt.Sprintln(err)
				}
				o += fmt.Sprintln(res)
				if res, err = lrpcHnd["help"].Handler(cmd, ws, cc); Check(err) {
					errString += fmt.Sprintln(err)
				}
				o += fmt.Sprintln(res)
				splitted := strings.Split(o, "\n")
				sort.Strings(splitted)
				var dedup []string
				for i := range splitted {
					if i > 0 {
						if splitted[i] != prev {
							dedup = append(dedup, splitted[i])
						}
					}
					prev = splitted[i]
				}
				o = strings.Join(dedup, "\n")
				if errString != "" {
					o += "BTCJSONError:\n"
					o += errString
				}
			} else {
				method = args[0]
				Debug("finding help for command", method)
				if help, err := rpcSrv.HelpCacher.RPCMethodHelp(method); Check(err) {
					o += err.Error() + "\n"
					o += fmt.Sprintln(res)
					cmd = &btcjson.HelpCmd{Command: &method}
					if res, err = lrpcHnd["help"].
						Handler(cmd, ws, cc); Check(err) {
						errString += fmt.Sprintln(err)
					}
					o += fmt.Sprintln(res)
				} else {
					o += help
				}
				// if _, ok := lrpcHnd[method]; ok {
				// 	o += "wallet server:\n"
				// 	o += legacy.HelpDescsEnUS()[method]
				// }
				// if _, ok := rpc.RPCHandlers[method]; ok {
				// 	o += "chain server:\n"
				// 	o += rpc.HelpDescsEnUS[method]
				// }
			}
			out <- o
		} else {
			params := make([]interface{}, 0, len(split[1:]))
			for _, arg := range args {
				params = append(params, arg)
			}
			if cmd, err = btcjson.NewCmd(method, params...); Check(err) {
				o += fmt.Sprintln(err)
			}
			if x, ok := chainrpc.RPCHandlers[method]; !ok {
				if x, ok := lrpcHnd[method]; ok {
					if res, err = x.Handler(cmd, ws, cc); Check(err) {
						o += err.Error()
					}
				}
			} else {
				if res, err = x.Fn(rpcSrv, cmd, nil); Check(err) {
					o += err.Error()
				}
			}
			if res != nil {
				if j, err := json.MarshalIndent(res, "",
					"  "); !Check(err) {
					o += string(j)
				}
			}
			out <- o
		}
	}()
	return
}

func (r *RcVar) CreateNewAddress(acctName string) string {
	if account, err := r.cx.WalletServer.AccountNumber(waddrmgr.
		KeyScopeBIP0044, acctName); !Check(err) {
		if addr, err := r.cx.WalletServer.NewAddress(account,
			waddrmgr.KeyScopeBIP0044, true); !Check(err) {
			if addr == nil {
				return ""
			}
			Debug("low", addr.EncodeAddress())
			return addr.EncodeAddress()
		}
	}
	return "error"
}

func (r *RcVar) CreateWallet(privPassphrase, duoSeed, pubPassphrase, walletDir string) {
	var err error
	var seed []byte
	if walletDir == "" {
		walletDir = *r.cx.Config.WalletFile
	}
	l := wallet.NewLoader(r.cx.ActiveNet, *r.cx.Config.WalletFile, 250)
	if duoSeed != "" {
		seed, err = hex.DecodeString(duoSeed)
		if err != nil {
			// Need to make JS invocation to embed
			Error(err)
		}
	} else if seed, err = hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen); Check(err) {
	}
	if _, err = l.CreateNewWallet([]byte(pubPassphrase),
		[]byte(privPassphrase), seed, time.Now(), true,
		r.cx.Config); Check(err) {
	}

	r.Boot.IsFirstRun = false
	*r.cx.Config.WalletPass = pubPassphrase
	*r.cx.Config.WalletFile = walletDir
	save.Pod(r.cx.Config)
}

func (r *RcVar) DuoNodeService() error {
	r.cx.NodeKill = make(chan struct{})
	r.cx.Node.Store(false)
	var err error
	if !*r.cx.Config.NodeOff {
		go func() {
			Info(r.cx.Language.RenderText("goApp_STARTINGNODE"))
			// utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
			err = node.Main(r.cx, nil)
			if err != nil {
				Info("error running node:", err)
				os.Exit(1)
			}
		}()
	}
	interrupt.AddHandler(func() {
		Warn("interrupt received, shutting down node")
		close(r.cx.NodeKill)
	})
	return err
}

func (r *RcVar) DuoSend(wp string, ad string, am float64) func() {
	return func() {
		pass := legacy.RPCHandlers["walletpassphrase"].Result()
		if _, err := pass.WalletPassphraseWait(&btcjson.WalletPassphraseCmd{
			Passphrase: wp,
			Timeout:    int64(time.Second * 2),
		}); !Check(err) {
			send := legacy.RPCHandlers["sendtoaddress"].Result()
			if _, err = send.SendToAddressWait(&btcjson.SendToAddressCmd{
				Address:   ad,
				Amount:    am,
				Comment:   nil,
				CommentTo: nil,
			}); Check(err) {
			}
		}
	}
}

func (r *RcVar) DuoUIloggerController() {
	logChan := logi.L.AddLogChan()
	logi.L.SetLevel(*r.cx.Config.LogLevel, true, "pod")
	go func() {
	out:
		for {
			select {
			case n := <-logChan:
				le, ok := r.Log.LogMessages.Load().([]logi.Entry)
				if ok {
					le = append(le, n)
					// Once length exceeds MaxLogLength we trim off the start
					// to keep it the same size
					ll := len(le)
					if ll > MaxLogLength {
						le = le[ll-MaxLogLength:]
					}
					r.Log.LogMessages.Store(le)
				} else {
					r.Log.LogMessages.Store([]logi.Entry{n})
				}
			case <-r.Log.StopLogger:
				defer func() {
					r.Log.StopLogger = make(chan struct{})
				}()
				r.Log.LogMessages.Store([]logi.Entry{})
				logi.L.LogChan = nil
				break out
			}
		}
		close(logChan)
	}()
}

func (r *RcVar) GetAddressBook() func() {
	return func() {
		ab := r.db.DbReadAddressBook()
		addressbook := &model.DuoUIaddressBook{}
		minConf := 1
		// Intermediate data for each address.
		type AddrData struct {
			// Total amount received.
			amount util.Amount
			// tx     []string
			// Account which the address belongs to
			// account string
			index int
		}
		syncBlock := r.cx.WalletServer.Manager.SyncedTo()
		// Intermediate data for all addresses.
		allAddrData := make(map[string]AddrData)
		// Create an AddrData entry for each active address in the account.
		// Otherwise we'll just get addresses from transactions later.
		sortedAddrs, err := r.cx.WalletServer.SortedActivePaymentAddresses()
		if err != nil {
		}
		idx := 0
		for _, address := range sortedAddrs {
			// There might be duplicates, just overwrite them.
			allAddrData[address] = AddrData{
				index: idx,
			}
			idx++
		}
		var endHeight int32
		if minConf == 0 {
			endHeight = -1
		} else {
			endHeight = syncBlock.Height - int32(minConf) + 1
		}
		err = wallet.ExposeUnstableAPI(r.cx.WalletServer).RangeTransactions(
			0, endHeight, func(details []wtxmgr.TxDetails) (bool, error) {
				for _, tx := range details {
					for _, cred := range tx.Credits {
						pkScript := tx.MsgTx.TxOut[cred.Index].PkScript
						_, addrs, _, err := txscript.ExtractPkScriptAddrs(
							pkScript, r.cx.WalletServer.ChainParams())
						if err != nil {
							// Non standard script, skip.
							continue
						}
						for _, addr := range addrs {
							addrStr := addr.EncodeAddress()
							addrData, ok := allAddrData[addrStr]
							if ok {
								addrData.amount += cred.Amount
							} else {
								addrData = AddrData{
									amount: cred.Amount,
								}
							}
							allAddrData[addrStr] = addrData
						}
					}
				}
				return false, nil
			})
		if err != nil {
		}
		var addrs AddressSlice
		// Massage address data into output format.
		addressbook.Num = len(allAddrData)
		for address, addrData := range allAddrData {
			addr := btcjson.ListReceivedByAddressResult{
				Address: address,
				Amount:  addrData.amount.ToDUO(),
			}
			if r.AddressBook.ShowMiningAddresses == false &&
				ab[addr.Address] == "Mining" {
			} else {
				addrs = append(addrs, model.DuoUIaddress{
					Index:   addrData.index,
					Label:   ab[addr.Address],
					Account: addr.Account,
					Address: addr.Address,
					Amount:  addr.Amount,
					Copy:    new(gel.Button),
					QrCode:  new(gel.Button),
				})
			}
		}
		sort.Sort(addrs)
		addressbook.Addresses = addrs
		r.AddressBook = addressbook
		return
	}
}

func (r *RcVar) GetBlock(hash string) btcjson.GetBlockVerboseResult {
	verbose, verbosetx := true, true
	bcmd := btcjson.GetBlockCmd{hash, &verbose, &verbosetx}
	if bl, err := chainrpc.HandleGetBlock(r.cx.RPCServer, &bcmd, nil); !Check(err) {
		if gbvr, ok := bl.(btcjson.GetBlockVerboseResult); ok {
			return gbvr
		}
	}
	return btcjson.GetBlockVerboseResult{}
}

func (r *RcVar) GetBlockCount() {
	go func() {
		if getBlockCount, err := chainrpc.HandleGetBlockCount(r.cx.RPCServer,
			nil, nil); Check(err) {
			r.Status.Node.BlockCount.Store(uint64(getBlockCount.(int64)))
		}
	}()
	return
}

func (r *RcVar) GetBlockExcerpt(height int) (b model.DuoUIblock) {
	b = *new(model.DuoUIblock)
	hashHeight, err := r.cx.RPCServer.Cfg.Chain.BlockHashByHeight(int32(height))
	if err != nil {
		Error("Block Hash By Height:", err)
	}
	verbose, verbosetx := true, true
	bcmd := btcjson.GetBlockCmd{hashHeight.String(), &verbose, &verbosetx}
	if bl, err := chainrpc.HandleGetBlock(r.cx.RPCServer, &bcmd,
		nil); Check(err) {
		block := bl.(btcjson.GetBlockVerboseResult)
		b.Height = block.Height
		b.BlockHash = block.Hash
		b.Confirmations = block.Confirmations
		b.TxNum = block.TxNum
		// t := time.Unix(0, block.Time)
		// b.Time = t.Format("02/01/2006, 15:04:05")
		b.Time = block.Time
		b.Link = &gel.Button{}
	}
	return
}

func (r *RcVar) GetBlockHash(blockHeight int) string {
	hcmd := btcjson.GetBlockHashCmd{Index: int64(blockHeight)}
	if hash, err := chainrpc.HandleGetBlockHash(r.cx.RPCServer, &hcmd, nil); !Check(err) {
		return hash.(string)
	} else {
		return err.Error()
	}
}

func (r *RcVar) GetBlocksExcerpts() func() {
	re := r.Explorer
	return func() {
		re.Page.To = int(r.Status.Node.BlockCount.Load()) /
			re.PerPage.Value
		startBlock := re.Page.Value * re.PerPage.Value
		endBlock := re.Page.Value*re.PerPage.Value +
			re.PerPage.Value
		height := int(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Height)
		Debug("GetBlocksExcerpts", startBlock, endBlock, height)
		if endBlock > height {
			endBlock = height
		}
		blocks := *new([]model.DuoUIblock)
		for i := startBlock; i < endBlock; i++ {
			blocks = append(blocks, r.GetBlockExcerpt(i))
			// Info("trazo")
			// Info(r.Status.Node.BlockHeight)
		}
		re.Blocks = blocks
		return
	}
}

func (r *RcVar) GetConnectionCount() {
	go r.Status.Node.ConnectionCount.Store(r.cx.RealNode.ConnectedCount())
	return
}

func (r *RcVar) GetDifficulty() {
	go func() {
		c := btcjson.GetDifficultyCmd{}
		diff, err := chainrpc.HandleGetDifficulty(r.cx.RPCServer, c, nil)
		if err != nil {
			// dv.PushDuoVUEalert("BTCJSONError", err.BTCJSONError(), "error")
		}
		r.Status.Node.Difficulty.Store(diff.(float64))
	}()
	return
}

// func (v *DuoVUEnode) Gethashespersec() {
// 	r, err := v.r.cx.RPCServer.HandleGetHashesPerSec(v.r.cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }
// func (v *DuoVUEnode) Getheaders(a *btcjson.GetHeadersCmd) {
// 	r, err := v.r.cx.RPCServer.HandleGetHeaders(v.r.cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuoVUEnode) Getinfo() {
// 	r, err := v.r.cx.RPCServer.HandleGetInfo(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.InfoChainResult{}
// 	return
// }
// func (v *DuoVUEnode) Getmempoolinfo() {
// 	r, err := v.r.cx.RPCServer.HandleGetMempoolInfo(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.GetMempoolInfoResult{}
// 	return
// }
// func (v *DuoVUEnode) Getmininginfo() {
// 	r, err := v.r.cx.RPCServer.HandleGetMiningInfo(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.GetMiningInfoResult{}
// 	return
// }
// func (v *DuoVUEnode) Getnettotals() {
// 	r, err := v.r.cx.RPCServer.HandleGetNetTotals(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.GetNetTotalsResult{}
// 	return
// }
// func (v *DuoVUEnode) Getnetworkhashps(a *btcjson.GetNetworkHashPSCmd) {
// 	r, err := v.r.cx.RPCServer.HandleGetNetworkHashPS(v.r.cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }

// func (v *DuoVUEnode) Stop() {
// 	r, err := v.r.cx.RPCServer.HandleStop(v.r.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
func (r *RcVar) GetDuoUIbalance() {
	go func() {
		Trace("getting balance")
		acct := "default"
		minconf := 0
		if getBalance, err := legacy.GetBalance(
			&btcjson.GetBalanceCmd{Account: &acct, MinConf: &minconf},
			r.cx.WalletServer); !Check(err) {
			// r.PushDuoUIalert("Error", err.Error(), "error")
			if gb, ok := getBalance.(float64); ok {
				bb := fmt.Sprintf("%0.8f", gb)
				r.Status.Wallet.Balance.Store(bb)
			}
		}
	}()
}

func (r *RcVar) GetDuoUIbestBlockHash() {
	go r.Status.Node.BestBlock.Store(r.cx.RPCServer.Cfg.Chain.BestSnapshot().
		Hash.String())
	return
}

func (r *RcVar) GetDuoUIblockCount() {
	go func() {
		if getBlockCount, err := chainrpc.HandleGetBlockCount(r.cx.RPCServer, nil,
			nil); !Check(err) {
			r.Status.Node.BlockCount.Store(uint64(getBlockCount.(int64)))
		} else {
			// r.PushDuoUIalert("BTCJSONError", err.BTCJSONError(), "error")
		}
	}()
	// Info(getBlockCount)
}

func (r *RcVar) GetDuoUIblockHeight() {
	go r.Status.Node.BlockHeight.
		Store(uint64(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Height))
}

func (r *RcVar) GetDuoUIconnectionCount() {
	go r.Status.Node.ConnectionCount.Store(r.cx.RealNode.ConnectedCount())
}

func (r *RcVar) GetDuoUIdifficulty() {
	go r.Status.Node.Difficulty.Store(
		chainrpc.GetDifficultyRatio(r.cx.RPCServer.Cfg.Chain.BestSnapshot().Bits,
			r.cx.RPCServer.Cfg.ChainParams, 2))
}

func (r *RcVar) GetDuoUIhashesPerSec() {
	// r.Status.Wallet.Hashes = int64(r.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	go func() {
		Debug("centralise hash function stuff here") // cpuminer
		r.Status.Kopach.Hashrate = r.cx.Hashrate.Load()
	}()
}

func (r *RcVar) GetDuoUIhashesPerSecList() {
	// // Create a new ring of size 5
	// hps := ring.New(3)
	// //GetDuoUIhashesPerSec
	// // Get the length of the ring
	// n := hps.Len()
	//
	// // Initialize the ring with some integer values
	// for i := 0; i < n; i++ {
	go r.GetDuoUIhashesPerSec()
	// hps.Value = r.Status.Kopach.Hashrate
	//	hps = hps.Next()
	// }
	//
	// // Iterate through the ring and print its contents
	// hps.Do(func(p interface{}) {
	//	r.Status.Kopach.Hps = append(r.Status.Kopach.Hps, p.(float64))
	//
	//fmt.Println(r.Status.Kopach.Hashrate)
	// })

}

func (r *RcVar) GetDuoUIlocalHost() {
	r.Localhost = *new(model.DuoUIlocalHost)
	// sm, _ := mem.VirtualMemory()
	// sc, _ := cpu.Info()
	// sp, _ := cpu.Percent(0, true)
	// sd, _ := disk.Usage("/")
	// r.Localhost.Cpu = sc
	// r.Localhost.CpuPercent = sp
	// r.Localhost.Memory = *sm
	// r.Localhost.Disk = *sd
	return
}

func (r *RcVar) GetDuoUInetworkHashesPerSec() {
	go func() {
		if networkHashesPerSecIface, err :=
			chainrpc.HandleGetNetworkHashPS(
				r.cx.RPCServer,
				btcjson.NewGetNetworkHashPSCmd(nil, nil), nil,
			); !Check(err) {
			if networkHashesPerSec, ok := networkHashesPerSecIface.(int64); ok {
				r.Status.Node.NetHash.Store(uint64(networkHashesPerSec))
			}
		}
	}()
}

func (r *RcVar) GetDuoUInetworkLastBlock() {
	go func() {
		for _, g := range r.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
			l := g.ToPeer().StatsSnapshot().LastBlock
			if l > r.Status.Node.NetworkLastBlock.Load() {
				r.Status.Node.NetworkLastBlock.Store(l)
			}
		}
	}()
	return
}

func (r *RcVar) GetDuoUIstatus() {
	go func() {
		if v, err := chainrpc.HandleVersion(r.cx.RPCServer, nil, nil); Check(err) {
			r.Status.Version = "0.0.1"
			r.Status.Wallet.WalletVersion = v.(map[string]btcjson.VersionResult)
			r.Status.StartTime = time.Unix(0, r.cx.RPCServer.Cfg.StartupTime)
			r.Status.CurrentNet = r.cx.RPCServer.Cfg.ChainParams.Net.String()
			r.Status.Chain = r.cx.RPCServer.Cfg.ChainParams.Name
		}
	}()
}

func (r *RcVar) GetDuoUItransactions() func() {
	rh := r.History
	return func() {
		rh.Page.To = int(r.Status.Wallet.TxsNumber.Load()) /
			rh.PerPage.Value
		startTx := rh.Page.Value * rh.PerPage.Value
		// endTx := rh.Page.Value*rh.PerPage.Value + rh.PerPage.Value
		Debug("getting transactions")
		// account, txcount, startnum, watchonly := "*", n, f, false
		// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{
		// Account: &account, Count: &txcount, From: &startnum,
		// IncludeWatchOnly: &watchonly}, v.ws)
		lt, err := r.cx.WalletServer.ListTransactions(startTx, rh.PerPage.Value)
		if err != nil {
			Info(err)
		}
		rh.Txs.TxsNumber = len(lt)
		txsArray := *new([]model.DuoUItransactionExcerpt)
		// lt := listTransactionsz.([]json.ListTransactionsResult)
		switch rh.Category {
		case "received":
			for _, tx := range lt {
				if tx.Category == "received" {
					txsArray = append(txsArray, txs(tx))
				}
			}
		case "sent":
			for _, tx := range lt {
				if tx.Category == "sent" {
					txsArray = append(txsArray, txs(tx))
				}
			}
		case "immature":
			for _, tx := range lt {
				if tx.Category == "immature" {
					txsArray = append(txsArray, txs(tx))
				}
			}
		case "generate":
			for _, tx := range lt {
				if tx.Category == "generate" {
					txsArray = append(txsArray, txs(tx))
				}
			}
		default:
			for _, tx := range lt {
				txsArray = append(txsArray, txs(tx))
			}
		}
		rh.Txs.Txs = txsArray
		return
	}
}

func (r *RcVar) GetDuoUItransactionsNumber() {
	go func() {
		Debug("getting transaction count")
		// account, txcount, startnum, watchonly := "*", n, f, false
		// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{
		// Account: &account, Count: &txcount, From: &startnum,
		// IncludeWatchOnly: &watchonly,
		// }, v.ws)
		if lt, err := r.cx.WalletServer.ListTransactions(0,
			999999999); Check(err) {
			r.Status.Wallet.TxsNumber.Store(uint64(len(lt)))
		}
	}()
}

func (r *RcVar) GetDuoUIunconfirmedBalance() {
	go func() {
		Trace("getting unconfirmed balance")
		acct := "default"
		if getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(
			&btcjson.GetUnconfirmedBalanceCmd{Account: &acct},
			r.cx.WalletServer); !Check(err) {
			if ub, ok := getUnconfirmedBalance.(float64); ok {
				ubb := fmt.Sprintf("%0.8f", ub)
				r.Status.Wallet.Unconfirmed.Store(ubb)
			}
		}
	}()
}

func (r *RcVar) GetLatestTransactions() {
	go func() {
		ltx := r.Status.Wallet.LastTxs
		Trace("getting latest transactions")
		lt, err := r.cx.WalletServer.ListTransactions(0, 10)
		if err != nil {
			// //r.PushDuoUIalert("BTCJSONError", err.BTCJSONError(), "error")
		}
		ltx.TxsNumber = len(lt)
		// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
		//	lt[i], lt[j] = lt[j], lt[i]
		// }
		balanceHeight := 0.0
		txseRaw := []model.DuoUItransactionExcerpt{}
		for _, txRaw := range lt {
			unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
			txseRaw = append(txseRaw, model.DuoUItransactionExcerpt{
				// Balance:       txse.Balance + txRaw.Amount,
				Comment:       txRaw.Comment,
				Amount:        txRaw.Amount,
				Category:      txRaw.Category,
				Confirmations: txRaw.Confirmations,
				Time:          unixTimeUTC.Format(time.RFC3339),
				TxID:          txRaw.TxID,
				Link:          new(gel.Button),
			})
		}
		var balance float64
		txs := *new([]model.DuoUItransactionExcerpt)
		for _, tx := range txseRaw {
			balance = balance + tx.Amount
			tx.Balance = balance
			txs = append(txs, tx)
			if ltx.Balance > balanceHeight {
				balanceHeight = ltx.Balance
			}
		}
		ltx.Txs = txs
		ltx.BalanceHeight = balanceHeight
	}()
}

// func (r *RcVar) GetTransactions() func() {
//	return func() {
//		Debug("getting transactions")
//		lt, err := r.cx.WalletServer.ListTransactions(0, rh.Txs.TxsListNumber)
//		if err != nil {
//			// //r.PushDuoUIalert("BTCJSONError", err.BTCJSONError(), "error")
//		}
//		rh.Txs.TxsNumber = len(lt)
//		// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
//		//	lt[i], lt[j] = lt[j], lt[i]
//		// }
//		balanceHeight := 0.0
//		txseRaw := []model.DuoUItransactionExcerpt{}
//		for _, txRaw := range lt {
//			unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
//			txseRaw = append(txseRaw, model.DuoUItransactionExcerpt{
//				// Balance:       txse.Balance + txRaw.Amount,
//				Comment:       txRaw.Comment,
//				Amount:        txRaw.Amount,
//				Category:      txRaw.Category,
//				Confirmations: txRaw.Confirmations,
//				Time:          unixTimeUTC.Format(time.RFC3339),
//				TxID:          txRaw.TxID,
//			})
//		}
//		var balance float64
//		txs := *new([]model.DuoUItransactionExcerpt)
//		for _, tx := range txseRaw {
//			balance = balance + tx.Amount
//			tx.Balance = balance
//			txs = append(txs, tx)
//			if rh.Txs.Balance > balanceHeight {
//				balanceHeight = rh.Txs.Balance
//			}
//
//		}
//		rh.Txs.Txs = txs
//		rh.Txs.BalanceHeight = balanceHeight
//	}
// }

func (r *RcVar) GetNetworkLastBlock() (out int32) {
	for _, g := range r.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		out := g.ToPeer().StatsSnapshot().LastBlock
		if out > r.Status.Node.NetworkLastBlock.Load() {
			r.Status.Node.NetworkLastBlock.Store(out)
		}
	}
	return
}

func (r *RcVar) GetPeerInfo() func() {
	return func() {
		if getPeers, err := chainrpc.HandleGetPeerInfo(r.cx.RPCServer, nil,
			nil); !Check(err) {
			r.Network.Peers = getPeers.([]*btcjson.GetPeerInfoResult)
		}
	}
}

func (r *RcVar) GetSingleBlock(hash string) func() {
	return func() {
		r.Explorer.SingleBlock = r.GetBlock(hash)
	}
}

func (r *RcVar) GetSingleTx(txid string) func() {
	return func() {
		r.History.SingleTx = r.GetTx(txid)
	}
}

func (r *RcVar) GetTx(txid string) btcjson.GetTransactionResult {
	verbose := 1
	tcmd := btcjson.GetRawTransactionCmd{
		Txid:    txid,
		Verbose: &verbose,
	}

	//lt, err := r.cx.RPCServer .(startTx, r.History.PerPage.Value)
	//if err != nil {
	//	Info(err)
	//}

	if tx, err := chainrpc.HandleGetRawTransaction(r.cx.RPCServer, &tcmd,
		nil); !Check(err) {
		if gbvr, ok := tx.(btcjson.GetTransactionResult); ok {
			Debug("zekr", gbvr)
			Debug(txid)
			Debug("txtxtx", tx)
			return gbvr
		}
	}
	return btcjson.GetTransactionResult{}
}

func (r *RcVar) GetUptime() {
	if rRaw, err := chainrpc.HandleUptime(r.cx.RPCServer, nil,
		nil); !Check(err) {
		// rRaw = int64(0)
		r.Uptime = rRaw.(int)
	}
	return
}

// func (v *DuoVUEnode) Validateaddress(a *btcjson.ValidateAddressCmd) {
// 	r, err := v.r.cx.RPCServer.HandleValidateAddress(v.r.cx.RPCServer, a, nil)
// 	r = btcjson.ValidateAddressChainResult{}
// 	return
// }
// func (v *DuoVUEnode) Verifychain(a *btcjson.VerifyChainCmd) {
// 	r, err := v.r.cx.RPCServer.HandleVerifyChain(v.r.cx.RPCServer, a, nil)
// }
// func (v *DuoVUEnode) Verifymessage(a *btcjson.VerifyMessageCmd) {
// 	r, err := v.r.cx.RPCServer.HandleVerifyMessage(v.r.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }

func (r *RcVar) GetWalletVersion() map[string]btcjson.VersionResult {
	v, err := chainrpc.HandleVersion(r.cx.RPCServer, nil, nil)
	if err != nil {
	}
	return v.(map[string]btcjson.VersionResult)
}

func (r *RcVar) labelMiningAddreses() {
	// r.db.DbReadAddressBook()
	ma := r.cx.ConfigMap["MiningAddrs"].(*cli.StringSlice)
	for _, miningAddress := range *ma {
		// for _, address := range r.AddressBook.Addresses {
		//	if miningAddress != address.Address {
		r.SaveAddressLabel(miningAddress, "Mining")
		// }
		// }
	}
}

func (r *RcVar) ListenInit(trigger chan struct{}) {
	Debug("listeninit")
	r.Events = EventsChan
	r.UpdateTrigger = trigger

	// first time starting up get all of these and trigger update
	update(r)
	r.labelMiningAddreses()

	var ready atomic.Bool
	ready.Store(false)
	r.cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted,
			blockchain.NTBlockConnected,
			blockchain.NTBlockDisconnected:
			if !ready.Load() {
				return
			}
			update(r)
			// go r.toastAdd("New block: "+
			// fmt.Sprint(callback.Data.(*util.Block).Height()),
			// callback.Data.(*util.Block).Hash().String())
		}
	})
	go func() {
		ticker := time.NewTicker(time.Second)
	out:
		for {
			select {
			case <-ticker.C:
				if !ready.Load() {
					if r.cx.IsCurrent() {
						ready.Store(true)
						// 		go func() {
						// 			r.cx.WalletServer.Rescan(nil, nil)
						// 			r.Ready <- struct{}{}
						// 			r.UpdateTrigger <- struct{}{}
						// 		}()
					}
				}
				r.GetDuoUIconnectionCount()
				r.UpdateTrigger <- struct{}{}
			// Warn("GetDuoUIconnectionCount")
			case <-r.cx.WalletServer.Update:
				update(r)
			case <-r.cx.KillAll:
				break out
			}
		}
	}()
	Warn("event update listener started")
	return
}

func (r *RcVar) SaveAddressLabel(address, label string) {
	// hf, err := highwayhash.New64(make([]byte, 32))
	// if err != nil {
	//	panic(err)
	// }
	// addressHash := hex.EncodeToString(hf.Sum([]byte(address)))
	addressbook := r.db.DbReadAddressBook()
	addressbook[address] = label
	r.db.DbWrite("user", "addressbook", addressbook)
}

func (r *RcVar) SaveDaemonCfg() {
	marshalled, _ := json.Marshal(r.Settings.Daemon.Config)
	config, _ := pod.EmptyConfig()
	if err := json.Unmarshal(marshalled, config); err != nil {
	}
	config2.Configure(r.cx, r.cx.AppContext.Command.Name)
	save.Pod(config)
}

func (r *RcVar) ShowAddressBook() func() {
	return func() {

	}
}

func (r *RcVar) StartServices() (err error) {
	Debug("starting up services")
	// Start Node
	if err = r.DuoNodeService(); Check(err) {
	}
	r.cx.RPCServer = <-r.cx.NodeChan
	r.cx.Node.Store(true)
	// Start wallet
	if err = r.StartWallet(); Check(err) {
	}
	r.cx.WalletServer = <-r.cx.WalletChan
	r.cx.Wallet.Store(true)
	r.cx.WalletServer.Rescan(nil, nil)
	<-r.cx.ChainClientReady
	legacy.RunAPI(r.cx.ChainClient, r.cx.WalletServer, r.cx.KillAll)
	r.Ready <- struct{}{}
	return
}

func (r *RcVar) StartWallet() error {
	r.cx.WalletKill = make(chan struct{})
	r.cx.Wallet.Store(false)
	var err error
	if !*r.cx.Config.WalletOff {
		go func() {
			Info("starting wallet")
			// utils.GetBiosMessage(view, "starting wallet")
			err = walletmain.Main(r.cx)
			if err != nil {
				fmt.Println("error running wallet:", err)
				os.Exit(1)
			}
		}()
	}
	interrupt.AddHandler(func() {
		Warn("interrupt received, " +
			"shutting down shell modules")
		close(r.cx.WalletKill)
	})
	return err
}

func (r *RcVar) toastAdd(t, m string) {
	r.Toasts = append(r.Toasts, model.DuoUItoast{
		Title:   t,
		Message: m,
	})
}

func (r *RcVar) UseTestnet() {
	*r.cx.Config.Network = "testnet"
	save.Pod(r.cx.Config)
}

// Items

func RcInit(cx *conte.Xt) (r *RcVar) {
	b := Boot{
		IsBoot:     true,
		IsFirstRun: false,
		IsBootMenu: false,
		IsBootLogo: false,
		IsLoading:  false,
		IsScreen:   "",
	}
	// d := models.DuoUIdialog{
	//	Show:   true,
	//	Ok:     func() { r.Dialog.Show = false },
	//	Cancel: func() { r.Dialog.Show = false },
	//	Title:  "Dialog!",
	//	Text:   "Dialog text",
	// }
	l := new(model.DuoUIlog)

	r = &RcVar{
		cx:          cx,
		db:          new(DB),
		Boot:        &b,
		AddressBook: new(model.DuoUIaddressBook),
		Status: &model.DuoUIstatus{
			Node: &model.NodeStatus{},
			Wallet: &model.WalletStatus{
				WalletVersion: make(map[string]btcjson.VersionResult),
				LastTxs:       &model.DuoUItransactionsExcerpts{},
			},
			Kopach: &model.KopachStatus{},
		},
		Dialog:   &model.DuoUIdialog{},
		Settings: settings(cx),
		Log:      l,
		Commands: new(DuoUIcommands),
		ConsoleHistory: &model.DuoUIconsoleHistory{
			Commands: []model.DuoUIconsoleCommand{
				{
					ComID:    "input",
					Category: "input",
					Time:     time.Now(),

					// Out: input(duo),
				},
			},
			CommandsNumber: 1,
		},
		Sent:      false,
		Localhost: model.DuoUIlocalHost{},
		ShowPage:  "OVERVIEW",
		Explorer: &model.DuoUIexplorer{
			PerPage: &gel.DuoUIcounter{
				Value:        20,
				OperateValue: 1,
				From:         0,
				To:           50,
				CounterInput: &gel.Editor{
					Alignment:  text.Middle,
					SingleLine: true,
				},
				CounterIncrease: new(gel.Button),
				CounterDecrease: new(gel.Button),
				CounterReset:    new(gel.Button),
			},
			Page: &gel.DuoUIcounter{
				Value:        0,
				OperateValue: 1,
				From:         0,
				To:           50,
				CounterInput: &gel.Editor{
					Alignment:  text.Middle,
					SingleLine: true,
				},
				CounterIncrease: new(gel.Button),
				CounterDecrease: new(gel.Button),
				CounterReset:    new(gel.Button),
			},
			Blocks:      []model.DuoUIblock{},
			SingleBlock: btcjson.GetBlockVerboseResult{},
		},

		Network: &model.DuoUInetwork{
			PerPage: &gel.DuoUIcounter{
				Value:        20,
				OperateValue: 1,
				From:         0,
				To:           50,
				CounterInput: &gel.Editor{
					Alignment:  text.Middle,
					SingleLine: true,
				},
				CounterIncrease: new(gel.Button),
				CounterDecrease: new(gel.Button),
				CounterReset:    new(gel.Button),
			},
			Page: &gel.DuoUIcounter{
				Value:        0,
				OperateValue: 1,
				From:         0,
				To:           50,
				CounterInput: &gel.Editor{
					Alignment:  text.Middle,
					SingleLine: true,
				},
				CounterIncrease: new(gel.Button),
				CounterDecrease: new(gel.Button),
				CounterReset:    new(gel.Button),
			},
			PeersList: &layout.List{
				Axis: layout.Vertical,
			},
			// Peers:     []*btcjson.GetPeerInfoResult
		},

		History: &model.DuoUIhistory{
			PerPage: &gel.DuoUIcounter{
				Value:        20,
				OperateValue: 1,
				From:         1,
				To:           50,
				CounterInput: &gel.Editor{
					Alignment:  text.Middle,
					SingleLine: true,
				},
				CounterIncrease: new(gel.Button),
				CounterDecrease: new(gel.Button),
				CounterReset:    new(gel.Button),
			},
			Page: &gel.DuoUIcounter{
				Value:        0,
				OperateValue: 1,
				From:         0,
				To:           50,
				CounterInput: &gel.Editor{
					Alignment:  text.Middle,
					SingleLine: true,
				},
				CounterIncrease: new(gel.Button),
				CounterDecrease: new(gel.Button),
				CounterReset:    new(gel.Button),
			},
			TransList: &layout.List{
				Axis: layout.Vertical,
			},
			Categories: &model.DuoUIhistoryCategories{
				AllTxs:      new(gel.CheckBox),
				MintedTxs:   new(gel.CheckBox),
				ImmatureTxs: new(gel.CheckBox),
				SentTxs:     new(gel.CheckBox),
				ReceivedTxs: new(gel.CheckBox),
			},
			Txs: &model.DuoUItransactionsExcerpts{
				ModelTxsListNumber: 0,
				TxsListNumber:      0,
				Txs:                []model.DuoUItransactionExcerpt{},
				TxsNumber:          0,
				Balance:            0,
				BalanceHeight:      0,
			},
			SingleTx: btcjson.GetTransactionResult{},
		},
		Quit:  make(chan struct{}),
		Ready: make(chan struct{}),
	}
	r.db.DuoUIdbInit(r.cx.DataDir)
	return
}

func settings(cx *conte.Xt) *model.DuoUIsettings {
	settings := &model.DuoUIsettings{
		Abbrevation: "DUO",
		Tabs: &model.DuoUIconfTabs{
			Current:  "wallet",
			TabsList: make(map[string]*gel.Button),
		},
		Daemon: &model.DaemonConfig{
			Config: cx.ConfigMap,
			Schema: pod.GetConfigSchema(cx.Config, cx.ConfigMap),
		},
	}
	// Settings tabs
	settingsFields := make(map[string]interface{})
	for _, group := range settings.Daemon.Schema.Groups {
		settings.Tabs.TabsList[group.Legend] = new(gel.Button)
		for _, field := range group.Fields {
			switch field.Type {
			case "stringSlice":
				settingsFields[field.Model] = &gel.Editor{
					SingleLine: false,
				}
				switch field.InputType {
				case "text":
					if settings.Daemon.Config[field.Model] != nil {
						var text string
						for _, str := range *settings.Daemon.
							Config[field.Model].(*cli.StringSlice) {
							text = text + str + "\n"
						}
						(settingsFields[field.Model]).(*gel.Editor).SetText(text)
					}
				}
			case "input":
				settingsFields[field.Model] = &gel.Editor{
					SingleLine: true,
				}
				if settings.Daemon.Config[field.Model] != nil {
					switch field.InputType {
					case "text":
						(settingsFields[field.Model]).(*gel.Editor).SetText(
							fmt.Sprint(*settings.Daemon.Config[field.Model].(*string)))
					case "number":
						(settingsFields[field.Model]).(*gel.Editor).SetText(
							fmt.Sprint(*settings.Daemon.Config[field.Model].(*int)))
					case "decimal":
						(settingsFields[field.Model]).(*gel.Editor).SetText(
							fmt.Sprintf("%0.f", *settings.Daemon.Config[field.
								Model].(*float64)))
					case "time":
						(settingsFields[field.Model]).(*gel.Editor).SetText(
							fmt.Sprint(*settings.
								Daemon.Config[field.Model].(*time.Duration)))
					case "password":
						(settingsFields[field.Model]).(*gel.Editor).SetText(
							fmt.Sprint(*settings.Daemon.Config[field.Model].(*string)))
					}
				}
			case "switch":
				settingsFields[field.Model] = new(gel.CheckBox)
				(settingsFields[field.Model]).(*gel.CheckBox).SetChecked(
					*settings.Daemon.Config[field.Model].(*bool))
			case "radio":
				settingsFields[field.Model] = new(gel.Enum)
			default:
				settingsFields[field.Model] = new(gel.Button)
			}
		}
	}
	settings.Daemon.Widgets = settingsFields
	return settings
}

func slug(text string) string {
	buf := make([]rune, 0, len(text))
	dash := false
	for _, r := range norm.NFKD.String(text) {
		switch {
		case unicode.IsOneOf(safe, r):
			buf = append(buf, unicode.ToLower(r))
			dash = true
		case unicode.IsOneOf(skip, r):
		case dash:
			buf = append(buf, '-')
			dash = false
		}
	}
	if i := len(buf) - 1; i >= 0 && buf[i] == '-' {
		buf = buf[:i]
	}
	return string(buf)
}

func txs(t btcjson.ListTransactionsResult) model.DuoUItransactionExcerpt {
	return model.DuoUItransactionExcerpt{
		TxID:     t.TxID,
		Amount:   t.Amount,
		Category: t.Category,
		Time:     helpers.FormatTime(time.Unix(t.Time, 0)),
		Link:     new(gel.Button),
	}

}

func update(r *RcVar) {
	// Warn("GetDuoUIbalance")
	r.GetDuoUIbalance()
	// Warn("GetDuoUIunconfirmedBalance")
	r.GetDuoUIunconfirmedBalance()
	// Warn("GetDuoUItransactionsNumber")
	r.GetDuoUItransactionsNumber()
	// r.GetTransactions()
	// Warn("GetLatestTransactions")
	r.GetLatestTransactions()
	// Info("")
	// Info("UPDATE")
	// Trace(r.History.PerPage)
	// Info("")
	// r.GetDuoUIstatus()
	// r.GetDuoUIlocalHost()
	// r.GetDuoUIblockHeight()
	// Warn("GetDuoUIblockCount")
	r.GetDuoUIdifficulty()
	r.GetDuoUIblockCount()
	r.GetPeerInfo()
	// Warn("GetDuoUIdifficulty")
	r.UpdateTrigger <- struct{}{}
}

// )
//
// type
// DuOShistory struct {
//	cx *conte.Xt
//	db DB
//	txs model.DuOStransactionsExcerpts
// }
//
// func (d *DuOShistory)GetDuOShistory() {
//	lt, err := d.cx.WalletServer.ListTransactions(0, 99999)
//	if err != nil {
//		//d.PushDuOSalert("Error", err.Error(), "error")
//	}
//	d.txs.TxsNumber = len(lt)
//	// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
//	//	lt[i], lt[j] = lt[j], lt[i]
//	// }
//	balanceHeight := 0.0
//	txseRaw := []DuOStransactionExcerpt{}
//	for _, txRaw := range lt {
//		unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
//		txseRaw = append(txseRaw, DuOStransactionExcerpt{
//			// Balance:       txse.Balance + txRaw.Amount,
//			Comment:       txRaw.Comment,
//			Amount:        txRaw.Amount,
//			Category:      txRaw.Category,
//			Confirmations: txRaw.Confirmations,
//			Time:          unixTimeUTC.Format(time.RFC3339),
//			TxID:          txRaw.TxID,
//		})
//	}
//	var balance float64
//	for _, tx := range txseRaw {
//		balance = balance + tx.Amount
//		tx.Balance = balance
//		d.txs.Txs = append(d.txs.Txs, tx)
//		if d.txs.Balance > balanceHeight {
//			balanceHeight = d.txs.Balance
//		}
//	}
//	d.txs.BalanceHeight = balanceHeight
//	Info("HISTORY-test:VAR->", d.txs)
//
//	return
// }
