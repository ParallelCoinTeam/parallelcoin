package mod

type RcVar struct {
	//db            *DuoUIdb
	Boot *Boot
	//Events        chan Event
	//UpdateTrigger chan struct{}
	//Status         *model.DuoUIstatus
	Dialog *Dialog
	//Log            *model.DuoUIlog
	//ConsoleHistory *model.DuoUIconsoleHistory

	Commands *Commands

	//Settings  *model.DuoUIsettings
	//Sent      bool
	//Toasts    []model.DuoUItoast
	//Localhost model.DuoUIlocalHost
	Uptime int
	// NodeChan   chan *rpc.Server
	// WalletChan chan *wallet.Wallet
	//Quit    chan struct{}
	//Ready   chan struct{}
	//IsReady bool
}

type Boot struct {
	IsBoot     bool
	IsFirstRun bool
	IsBootMenu bool
	IsBootLogo bool
	IsLoading  bool
	IsScreen   string
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
