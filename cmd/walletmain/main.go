package walletmain

import (
	"fmt"
	"io/ioutil"
	// This enables pprof
	// _ "net/http/pprof"
	"sync"
	
	"github.com/p9c/pod/pkg/util/logi"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	"github.com/p9c/pod/pkg/wallet/chain"
)

// Main is a work-around main function that is required since deferred functions (such as log flushing) are not called
// with calls to os.Exit. Instead, main runs this function and checks for a non-nil error, at point any defers have
// already run, and if the error is non-nil, the program can be exited with an error exit status.
func Main(cx *conte.Xt) (err error) {
	// cx.WaitGroup.Add(1)
	cx.WaitAdd()
	// if *config.Profile != "" {
	//	go func() {
	//		listenAddr := net.JoinHostPort("127.0.0.1", *config.Profile)
	//		Info("profile server listening on", listenAddr)
	//		profileRedirect := http.RedirectHandler("/debug/pprof",
	//			http.StatusSeeOther)
	//		http.Handle("/", profileRedirect)
	//		fmt.Println(http.ListenAndServe(listenAddr, nil))
	//	}()
	// }
	loader := wallet.NewLoader(cx.ActiveNet, *cx.Config.WalletFile, 250)
	// Create and start HTTP server to serve wallet client connections. This will be updated with the wallet and chain
	// server RPC client created below after each is created.
	Debug("starting RPC servers")
	var legacyServer *legacy.Server
	if legacyServer, err = startRPCServers(cx, loader); Check(err) {
		Error("unable to create RPC servers:", err)
		return
	}
	loader.RunAfterLoad(
		func(w *wallet.Wallet) {
			Warn("starting wallet RPC services", w != nil)
			startWalletRPCServices(w, legacyServer)
			// cx.WalletChan <- w
		},
	)
	if !*cx.Config.NoInitialLoad {
		go func() {
			Debug("loading wallet", *cx.Config.WalletPass)
			if err = LoadWallet(loader, cx, legacyServer); Check(err) {
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
		Warn("wallet killswitch activated")
		if legacyServer != nil {
			Warn("stopping wallet RPC server")
			legacyServer.Stop()
			Info("stopped wallet RPC server")
		}
		Info("wallet shutdown from killswitch complete")
		// cx.WaitGroup.Done()
		cx.WaitDone()
		return
		// <-legacyServer.RequestProcessShutdownChan()
	case <-cx.KillAll.Wait():
		Debug("killall")
		cx.WalletKill.Q()
	case <-interrupt.HandlersDone.Wait():
	}
	Info("wallet shutdown complete")
	// cx.WaitGroup.Done()
	cx.WaitDone()
	return
}

func ReadCAFile(config *pod.Config) []byte {
	// Read certificate file if TLS is not disabled.
	var certs []byte
	if *config.TLS {
		var err error
		certs, err = ioutil.ReadFile(*config.CAFile)
		if err != nil {
			Error("cannot open CA file:", err)
			// If there's an error reading the CA file, continue with nil certs and without the client connection.
			certs = nil
		}
	} else {
		Info("chain server RPC TLS is disabled")
	}
	return certs
}

func LoadWallet(loader *wallet.Loader, cx *conte.Xt, legacyServer *legacy.Server) (err error) {
	Debug("starting rpc client connection handler", *cx.Config.WalletPass)
	// Create and start chain RPC client so it's ready to connect to the wallet when loaded later.
	// Load the wallet database. It must have been created already or this will return an appropriate error.
	var w *wallet.Wallet
	Debug("^^^^^^^^^^^^^ opening existing wallet")
	if w, err = loader.OpenExistingWallet([]byte(*cx.Config.WalletPass), true, cx.Config, nil); Check(err) {
		Debug("^^^^^^^^^^^^^ failed to open existing wallet")
		return
	}
	Debug("^^^^^^^^^^^^^ opened existing wallet")
	// go func() {
	// Warn("refilling mining addresses", cx.Config, cx.StateCfg)
	// addresses.RefillMiningAddresses(w, cx.Config, cx.StateCfg)
	// Warn("done refilling mining addresses")
	// Debugs(*cx.Config.MiningAddrs)
	// save.Pod(cx.Config)
	// }()
	loader.Wallet = w
	// Debug("^^^^^^^^^^^ sending back wallet")
	cx.WalletChan <- w
	Debug("starting rpcClientConnectLoop")
	go rpcClientConnectLoop(cx, legacyServer, loader)
	Debug("^^^^^^^^^^^^^ adding interrupt handler to unload wallet")
	// Add interrupt handlers to shutdown the various process components before exiting. Interrupt handlers run in
	// LIFO order, so the wallet (which should be closed last) is added first.
	interrupt.AddHandler(
		func() {
			Debug("wallet.Main interrupt")
			err := loader.UnloadWallet()
			if err != nil && err != wallet.ErrNotLoaded {
				Error("failed to close wallet:", err)
			}
		},
	)
	if legacyServer != nil {
		interrupt.AddHandler(
			func() {
				Trace("stopping wallet RPC server")
				legacyServer.Stop()
				Trace("wallet RPC server shutdown")
			},
		)
	}
	go func() {
	out:
		for {
			select {
			case <-cx.KillAll.Wait():
				break out
			case <-legacyServer.RequestProcessShutdownChan().Wait():
				interrupt.Request()
			}
		}
	}()
	return
}

// rpcClientConnectLoop continuously attempts a connection to the consensus RPC server. When a connection is
// established, the client is used to sync the loaded wallet, either immediately or when loaded at a later time.
//
// The legacy RPC is optional. If set, the connected RPC client will be associated with the server for RPC pass-through
// and to enable additional methods.
func rpcClientConnectLoop(
	cx *conte.Xt, legacyServer *legacy.Server,
	loader *wallet.Loader,
) {
	Debug("^^^^^^^^^^^^^^^ rpcClientConnectLoop", logi.Caller("which was started at:", 2))
	// var certs []byte
	// if !cx.PodConfig.UseSPV {
	certs := ReadCAFile(cx.Config)
	// }
	for {
		var (
			chainClient chain.Interface
			err         error
		)
		// if cx.PodConfig.UseSPV {
		// 	var (
		// 		chainService *neutrino.ChainService
		// 		spvdb        walletdb.DB
		// 	)
		// 	netDir := networkDir(cx.PodConfig.AppDataDir.value, ActiveNet.Params)
		// 	spvdb, err = walletdb.Create("bdb",
		// 		filepath.Join(netDir, "neutrino.db"))
		// 	defer spvdb.Close()
		// 	if err != nil {
		// 		log<-cl.Errorf{"unable to create Neutrino DB: %s", err)
		// 		continue
		// 	}
		// 	chainService, err = neutrino.NewChainService(
		// 		neutrino.Config{
		// 			DataDir:      netDir,
		// 			Database:     spvdb,
		// 			ChainParams:  *ActiveNet.Params,
		// 			ConnectPeers: cx.PodConfig.ConnectPeers,
		// 			AddPeers:     cx.PodConfig.AddPeers,
		// 		})
		// 	if err != nil {
		// 		log<-cl.Errorf{"couldn't create Neutrino ChainService: %s", err)
		// 		continue
		// 	}
		// 	chainClient = chain.NewNeutrinoClient(ActiveNet.Params, chainService)
		// 	err = chainClient.Start()
		// 	if err != nil {
		// 		log<-cl.Errorf{"couldn't start Neutrino client: %s", err)
		// 	}
		// } else {
		var cc *chain.RPCClient
		Debug("starting wallet's ChainClient")
		cc, err = StartChainRPC(cx.Config, cx.ActiveNet, certs, cx.KillAll)
		if err != nil {
			Error(
				"unable to open connection to consensus RPC server:", err,
			)
			continue
		}
		Debug("^^^^^^^^^ storing chain client")
		cx.ChainClient = cc
		cx.ChainClientReady.Q()
		chainClient = cc
		// Rather than inlining this logic directly into the loader callback, a function variable is used to avoid
		// running any of this after the client disconnects by setting it to nil. This prevents the callback from
		// associating a wallet loaded at a later time with a client that has already disconnected. A mutex is used to
		// make this concurrent safe.
		associateRPCClient := func(w *wallet.Wallet) {
			Debug(">>>>>>>>>>> associating chain client")
			if w != nil {
				w.SynchronizeRPC(chainClient)
			}
			if legacyServer != nil {
				legacyServer.SetChainServer(chainClient)
			}
		}
		Debug("adding wallet loader hook to connect to chain")
		mu := new(sync.Mutex)
		loader.RunAfterLoad(
			func(w *wallet.Wallet) {
				Debug("running associate chain client")
				mu.Lock()
				associate := associateRPCClient
				mu.Unlock()
				if associate != nil {
					associate(w)
					Debug("wallet is now associated by chain client")
				} else {
					Debug("wallet chain client associate function is nil")
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

// StartChainRPC opens a RPC client connection to a pod server for blockchain services. This function uses the RPC
// options from the global config and there is no recovery in case the server is not available or if there is an
// authentication error. Instead, all requests to the client will simply error.
func StartChainRPC(config *pod.Config, activeNet *netparams.Params, certs []byte, quit qu.C) (*chain.RPCClient, error) {
	Debug(">>>>>>>>>>>>>>> attempting RPC client connection to %v, TLS: %s", *config.RPCConnect, fmt.Sprint(*config.TLS),
	)
	rpcC, err := chain.NewRPCClient(
		activeNet,
		*config.RPCConnect,
		*config.Username,
		*config.Password,
		certs,
		*config.TLS,
		0,
		quit,
	)
	if err != nil {
		Error(err)
		return nil, err
	}
	err = rpcC.Start()
	return rpcC, err
}
