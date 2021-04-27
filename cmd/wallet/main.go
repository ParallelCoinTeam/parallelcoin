package wallet

import (
	"fmt"
	// This enables pprof
	// _ "net/http/pprof"
	"sync"

	"github.com/p9c/qu"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pod/config"
	"github.com/p9c/matrjoska/pod/state"

	"github.com/p9c/interrupt"

	"github.com/p9c/matrjoska/pkg/chainclient"
)

// Main is a work-around main function that is required since deferred functions
// (such as log flushing) are not called with calls to os.Exit. Instead, main
// runs this function and checks for a non-nil error, at point any defers have
// already run, and if the error is non-nil, the program can be exited with an
// error exit status.
func Main(cx *state.State) (e error) {
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
	loader := NewLoader(cx.ActiveNet, cx.Config.WalletFile.V(), 250)
	// Create and start HTTP server to serve wallet client connections. This will be updated with the wallet and chain
	// server RPC client created below after each is created.
	D.Ln("starting RPC servers")
	var legacyServer *Server
	if legacyServer, e = startRPCServers(cx, loader); E.Chk(e) {
		E.Ln("unable to create RPC servers:", e)
		return
	}
	loader.RunAfterLoad(
		func(w *Wallet) {
			D.Ln("starting wallet RPC services", w != nil)
			startWalletRPCServices(w, legacyServer)
			// cx.WalletChan <- w
		},
	)
	if !cx.Config.NoInitialLoad.True() {
		go func() {
			D.Ln("loading wallet", cx.Config.WalletPass.V())
			if e = LoadWallet(loader, cx, legacyServer); E.Chk(e) {
			}
		}()
	}
	interrupt.AddHandler(cx.WalletKill.Q)
	select {
	case <-cx.WalletKill.Wait():
		D.Ln("wallet killswitch activated")
		if legacyServer != nil {
			D.Ln("stopping wallet RPC server")
			legacyServer.Stop()
			I.Ln("stopped wallet RPC server")
		}
		I.Ln("wallet shutdown from killswitch complete")
		cx.WaitDone()
		return
	case <-cx.KillAll.Wait():
		D.Ln("killall")
		cx.WalletKill.Q()
	case <-interrupt.HandlersDone.Wait():
	}
	I.Ln("wallet shutdown complete")
	cx.WaitDone()
	return
}

// LoadWallet ...
func LoadWallet(
	loader *Loader, cx *state.State, legacyServer *Server,
) (e error) {
	T.Ln("starting rpc client connection handler", cx.Config.WalletPass.V())
	// Create and start chain RPC client so it's ready to connect to the wallet when
	// loaded later. Load the wallet database. It must have been created already or
	// this will return an appropriate error.
	var w *Wallet
	T.Ln("opening existing wallet, pass:", cx.Config.WalletPass.V())
	if w, e = loader.OpenExistingWallet(cx.Config.WalletPass.Bytes(), true, cx.Config, nil); E.Chk(e) {
		T.Ln("failed to open existing wallet")
		return
	}
	T.Ln("opened existing wallet")
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
	T.Ln("starting rpcClientConnectLoop")
	go rpcClientConnectLoop(cx, legacyServer, loader)
	T.Ln("adding interrupt handler to unload wallet")
	// Add interrupt handlers to shutdown the various process components before
	// exiting. Interrupt handlers run in LIFO order, so the wallet (which should be
	// closed last) is added first.
	interrupt.AddHandler(
		func() {
			D.Ln("wallet.CtlMain interrupt")
			e := loader.UnloadWallet()
			if e != nil && e != ErrNotLoaded {
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
	cx *state.State, legacyServer *Server,
	loader *Loader,
) {
	T.Ln("rpcClientConnectLoop", log.Caller("which was started at:", 2))
	// var certs []byte
	// if !cx.PodConfig.UseSPV {
	certs := cx.Config.ReadCAFile()
	// }
	for {
		var (
			chainClient chainclient.Interface
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
		var cc *chainclient.RPCClient
		T.Ln("starting wallet's ChainClient")
		cc, e = StartChainRPC(cx.Config, cx.ActiveNet, certs, cx.KillAll)
		if e != nil {
			E.Ln(
				"unable to open connection to consensus RPC server:", e,
			)
			continue
		}
		T.Ln("storing chain client")
		cx.ChainClient = cc
		cx.ChainClientReady.Q()
		chainClient = cc
		// Rather than inlining this logic directly into the loader callback, a function
		// variable is used to avoid running any of this after the client disconnects by
		// setting it to nil. This prevents the callback from associating a wallet
		// loaded at a later time with a client that has already disconnected. A mutex
		// is used to make this concurrent safe.
		associateRPCClient := func(w *Wallet) {
			T.Ln("associating chain client")
			if w != nil {
				w.SynchronizeRPC(chainClient)
			}
			if legacyServer != nil {
				legacyServer.SetChainServer(chainClient)
			}
		}
		T.Ln("adding wallet loader hook to connect to chain")
		mu := new(sync.Mutex)
		loader.RunAfterLoad(
			func(w *Wallet) {
				T.Ln("running associate chain client")
				mu.Lock()
				associate := associateRPCClient
				mu.Unlock()
				if associate != nil {
					associate(w)
					T.Ln("wallet is now associated by chain client")
				} else {
					T.Ln("wallet chain client associate function is nil")
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
func StartChainRPC(
	config *config.Config,
	activeNet *chaincfg.Params,
	certs []byte,
	quit qu.C,
) (rpcC *chainclient.RPCClient, e error) {
	D.F(
		"attempting RPC client connection to %v, TLS: %s",
		config.RPCConnect.V(),
		fmt.Sprint(config.ClientTLS.True()),
	)
	if rpcC, e = chainclient.NewRPCClient(
		activeNet,
		config.RPCConnect.V(),
		config.Username.V(),
		config.Password.V(),
		certs,
		config.ClientTLS.True(),
		0,
		quit,
	); E.Chk(e) {
		return nil, e
	}
	e = rpcC.Start()
	return rpcC, e
}
