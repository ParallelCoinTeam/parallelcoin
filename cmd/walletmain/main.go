package walletmain

import (
	"fmt"
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/logg"
	"github.com/p9c/pod/pkg/podcfg"
	
	// This enables pprof
	// _ "net/http/pprof"
	"sync"
	
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	"github.com/p9c/pod/pkg/wallet/chain"
)

// Main is a work-around main function that is required since deferred functions
// (such as log flushing) are not called with calls to os.Exit. Instead, main
// runs this function and checks for a non-nil error, at point any defers have
// already run, and if the error is non-nil, the program can be exited with an
// error exit status.
func Main(cx *pod.State) (e error) {
	// cx.WaitGroup.Add(1)
	cx.WaitAdd()
	// if *config.Profile != "" {
	//	go func() {
	//		listenAddr := net.JoinHostPort("127.0.0.1", *config.Profile)
	//		I.Ln("profile server listening on", listenAddr)
	//		profileRedirect := http.RedirectHandler("/debug/pprof",
	//			http.StatusSeeOther)
	//		http.Handle("/", profileRedirect)
	//		fmt.Println(http.ListenAndServe(listenAddr, nil))
	//	}()
	// }
	loader := wallet.NewLoader(cx.ActiveNet, *cx.Config.WalletFile, 250)
	// Create and start HTTP server to serve wallet client connections. This will be updated with the wallet and chain
	// server RPC client created below after each is created.
	D.Ln("starting RPC servers")
	var legacyServer *legacy.Server
	if legacyServer, e = startRPCServers(cx, loader); E.Chk(e) {
		E.Ln("unable to create RPC servers:", e)
		return
	}
	loader.RunAfterLoad(
		func(w *wallet.Wallet) {
			D.Ln("starting wallet RPC services", w != nil)
			startWalletRPCServices(w, legacyServer)
			// cx.WalletChan <- w
		},
	)
	if !*cx.Config.NoInitialLoad {
		go func() {
			D.Ln("loading wallet", *cx.Config.WalletPass)
			if e = LoadWallet(loader, cx, legacyServer); E.Chk(e) {
			}
		}()
	}
	interrupt.AddHandler(
		func() {
			cx.WalletKill.Q()
		},
	)
	select {
	case <-cx.WalletKill.Wait():
		D.Ln("wallet killswitch activated")
		if legacyServer != nil {
			D.Ln("stopping wallet RPC server")
			legacyServer.Stop()
			I.Ln("stopped wallet RPC server")
		}
		I.Ln("wallet shutdown from killswitch complete")
		// cx.WaitGroup.Done()
		cx.WaitDone()
		return
		// <-legacyServer.RequestProcessShutdownChan()
	case <-cx.KillAll.Wait():
		D.Ln("killall")
		cx.WalletKill.Q()
	case <-interrupt.HandlersDone.Wait():
	}
	I.Ln("wallet shutdown complete")
	// cx.WaitGroup.Done()
	cx.WaitDone()
	return
}

// LoadWallet ...
func LoadWallet(loader *wallet.Loader, cx *pod.State, legacyServer *legacy.Server) (e error) {
	D.Ln("starting rpc client connection handler", *cx.Config.WalletPass)
	// Create and start chain RPC client so it's ready to connect to the wallet when
	// loaded later. Load the wallet database. It must have been created already or
	// this will return an appropriate error.
	var w *wallet.Wallet
	D.Ln("^^^^^^^^^^^^^ opening existing wallet")
	if w, e = loader.OpenExistingWallet([]byte(*cx.Config.WalletPass), true, cx.Config, nil); E.Chk(e) {
		D.Ln("^^^^^^^^^^^^^ failed to open existing wallet")
		return
	}
	D.Ln("^^^^^^^^^^^^^ opened existing wallet")
	// go func() {
	// W.Ln("refilling mining addresses", cx.Config, cx.StateCfg)
	// addresses.RefillMiningAddresses(w, cx.Config, cx.StateCfg)
	// W.Ln("done refilling mining addresses")
	// D.S(*cx.Config.MiningAddrs)
	// save.Save(cx.Config)
	// }()
	loader.Wallet = w
	// D.Ln("^^^^^^^^^^^ sending back wallet")
	// cx.WalletChan <- w
	D.Ln("starting rpcClientConnectLoop")
	go rpcClientConnectLoop(cx, legacyServer, loader)
	D.Ln("^^^^^^^^^^^^^ adding interrupt handler to unload wallet")
	// Add interrupt handlers to shutdown the various process components before
	// exiting. Interrupt handlers run in LIFO order, so the wallet (which should be
	// closed last) is added first.
	interrupt.AddHandler(
		func() {
			D.Ln("wallet.Main interrupt")
			e := loader.UnloadWallet()
			if e != nil && e != wallet.ErrNotLoaded {
				E.Ln("failed to close wallet:", e)
			}
		},
	)
	if legacyServer != nil {
		interrupt.AddHandler(
			func() {
				D.Ln("stopping wallet RPC server")
				legacyServer.Stop()
				D.Ln("wallet RPC server shutdown")
			},
		)
	}
	go func() {
		select {
		case <-cx.KillAll.Wait():
		case <-legacyServer.RequestProcessShutdownChan().Wait():
		}
		interrupt.Request()
	}()
	return
}

// rpcClientConnectLoop continuously attempts a connection to the consensus RPC
// server. When a connection is established, the client is used to sync the
// loaded wallet, either immediately or when loaded at a later time.
//
// The legacy RPC is optional. If set, the connected RPC client will be
// associated with the server for RPC pass-through and to enable additional
// methods.
func rpcClientConnectLoop(
	cx *pod.State, legacyServer *legacy.Server,
	loader *wallet.Loader,
) {
	D.Ln("^^^^^^^^^^^^^^^ rpcClientConnectLoop", logg.Caller("which was started at:", 2))
	// var certs []byte
	// if !cx.PodConfig.UseSPV {
	certs := podcfg.ReadCAFile(cx.Config)
	// }
	for {
		var (
			chainClient chain.Interface
			e           error
		)
		// if cx.PodConfig.UseSPV {
		// 	var (
		// 		chainService *neutrino.ChainService
		// 		spvdb        walletdb.DB
		// 	)
		// 	netDir := networkDir(cx.PodConfig.AppDataDir.value, ActiveNet.Params)
		// 	spvdb, e = walletdb.Create("bdb",
		// 		filepath.Join(netDir, "neutrino.db"))
		// 	defer spvdb.Close()
		// 	if e != nil  {
		// 		log<-cl.Errorf{"unable to create Neutrino DB: %s", e)
		// 		continue
		// 	}
		// 	chainService, e = neutrino.NewChainService(
		// 		neutrino.Config{
		// 			DataDir:      netDir,
		// 			Database:     spvdb,
		// 			ChainParams:  *ActiveNet.Params,
		// 			ConnectPeers: cx.PodConfig.ConnectPeers,
		// 			AddPeers:     cx.PodConfig.AddPeers,
		// 		})
		// 	if e != nil  {
		// 		log<-cl.Errorf{"couldn't create Neutrino ChainService: %s", e)
		// 		continue
		// 	}
		// 	chainClient = chain.NewNeutrinoClient(ActiveNet.Params, chainService)
		// 	e = chainClient.Start()
		// 	if e != nil  {
		// 		log<-cl.Errorf{"couldn't start Neutrino client: %s", e)
		// 	}
		// } else {
		var cc *chain.RPCClient
		D.Ln("starting wallet's ChainClient")
		cc, e = StartChainRPC(cx.Config, cx.ActiveNet, certs, cx.KillAll)
		if e != nil {
			E.Ln(
				"unable to open connection to consensus RPC server:", e,
			)
			continue
		}
		D.Ln("^^^^^^^^^ storing chain client")
		cx.ChainClient = cc
		cx.ChainClientReady.Q()
		chainClient = cc
		// Rather than inlining this logic directly into the loader callback, a function
		// variable is used to avoid running any of this after the client disconnects by
		// setting it to nil. This prevents the callback from associating a wallet
		// loaded at a later time with a client that has already disconnected. A mutex
		// is used to make this concurrent safe.
		associateRPCClient := func(w *wallet.Wallet) {
			D.Ln(">>>>>>>>>>> associating chain client")
			if w != nil {
				w.SynchronizeRPC(chainClient)
			}
			if legacyServer != nil {
				legacyServer.SetChainServer(chainClient)
			}
		}
		D.Ln("adding wallet loader hook to connect to chain")
		mu := new(sync.Mutex)
		loader.RunAfterLoad(
			func(w *wallet.Wallet) {
				D.Ln("running associate chain client")
				mu.Lock()
				associate := associateRPCClient
				mu.Unlock()
				if associate != nil {
					associate(w)
					D.Ln("wallet is now associated by chain client")
				} else {
					D.Ln("wallet chain client associate function is nil")
				}
			},
		)
		chainClient.WaitForShutdown()
		mu.Lock()
		associateRPCClient = nil
		mu.Unlock()
		loadedWallet, ok := loader.LoadedWallet()
		if ok {
			// Do not attempt a reconnect when the wallet was explicitly stopped.
			if loadedWallet.ShuttingDown() {
				return
			}
			loadedWallet.SetChainSynced(false)
			// TODO: Rework the wallet so changing the RPC client does not
			//  require stopping and restarting everything.
			loadedWallet.Stop()
			loadedWallet.WaitForShutdown()
			loadedWallet.Start()
		}
	}
}

// StartChainRPC opens a RPC client connection to a pod server for blockchain
// services. This function uses the RPC options from the global config and there
// is no recovery in case the server is not available or if there is an
// authentication error. Instead, all requests to the client will simply error.
func StartChainRPC(config *podcfg.Config, activeNet *chaincfg.Params, certs []byte, quit qu.C) (*chain.RPCClient, error) {
	D.Ln(
		">>>>>>>>>>>>>>> attempting RPC client connection to %v, TLS: %s", *config.RPCConnect, fmt.Sprint(*config.TLS),
	)
	rpcC, e := chain.NewRPCClient(
		activeNet,
		*config.RPCConnect,
		*config.Username,
		*config.Password,
		certs,
		*config.TLS,
		0,
		quit,
	)
	if e != nil {
		return nil, e
	}
	e = rpcC.Start()
	return rpcC, e
}
