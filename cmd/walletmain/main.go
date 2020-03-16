package walletmain

import (
	"fmt"
	"io/ioutil"
	// This enables pprof
	// _ "net/http/pprof"
	"sync"

	"github.com/p9c/pod/pkg/chain/mining/addresses"
	"github.com/p9c/pod/pkg/conte"

	log "github.com/p9c/pod/pkg/logi"

	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	"github.com/p9c/pod/pkg/wallet/chain"
)

// Main is a work-around main function that is required since deferred
// functions (such as log flushing) are not called with calls to os.Exit.
// Instead, main runs this function and checks for a non-nil error, at point
// any defers have already run, and if the error is non-nil, the program can be
// exited with an error exit status.
func Main(cx *conte.Xt) (err error) {
	log.L.Info("starting wallet")
	cx.WaitGroup.Add(1)

	// if *config.Profile != "" {
	//	go func() {
	//		listenAddr := net.JoinHostPort("127.0.0.1", *config.Profile)
	//		log.L.Info("profile server listening on", listenAddr)
	//		profileRedirect := http.RedirectHandler("/debug/pprof",
	//			http.StatusSeeOther)
	//		http.Handle("/", profileRedirect)
	//		log.Println(http.ListenAndServe(listenAddr, nil))
	//	}()
	// }
	loader := wallet.NewLoader(cx.ActiveNet, *cx.Config.WalletFile, 250)
	// Create and start HTTP server to serve wallet client connections.
	// This will be updated with the wallet and chain server RPC client
	// created below after each is created.
	log.L.Trace("starting RPC servers")
	legacyServer, err := startRPCServers(cx.Config, cx.StateCfg, cx.ActiveNet,
		loader)
	if err != nil {
		log.L.Error("unable to create RPC servers:", err)
		return
	}
	loader.RunAfterLoad(func(w *wallet.Wallet) {
		log.L.Warn("starting wallet RPC services", w != nil)
		startWalletRPCServices(w, legacyServer)
	})
	if !*cx.Config.NoInitialLoad {
		log.L.Trace("starting rpc client connection handler")
		// Create and start chain RPC client so it's ready to connect to
		// the wallet when loaded later.
		log.L.Trace("loading database")
		// Load the wallet database.  It must have been created already
		// or this will return an appropriate error.
		var w *wallet.Wallet
		w, err = loader.OpenExistingWallet([]byte(*cx.Config.WalletPass),
			true, cx.Config)
		// log.L.Warn("wallet", w)
		if err != nil {
			log.L.Error(err)
			return
		}
		go func() {
			addresses.RefillMiningAddresses(w, cx.Config, cx.StateCfg)
		}()
		go rpcClientConnectLoop(cx, legacyServer, loader)
		loader.Wallet = w
		legacy.RunAPI(nil, w, cx.WalletKill)
		log.L.Trace("sending back wallet")
		cx.WalletChan <- w
		log.L.Trace("adding interrupt handler to unload wallet")
		// Add interrupt handlers to shutdown the various process components
		// before exiting.  Interrupt handlers run in LIFO order, so the wallet
		// (which should be closed last) is added first.
		interrupt.AddHandler(func() {
			log.L.Debug("wallet.Main interrupt")
			err := loader.UnloadWallet()
			if err != nil && err != wallet.ErrNotLoaded {
				log.L.Error("failed to close wallet:", err)
			}
		})
		if legacyServer != nil {
			interrupt.AddHandler(func() {
				log.L.Trace("stopping wallet RPC server")
				legacyServer.Stop()
				log.L.Trace("wallet RPC server shutdown")
			})
		}
		go func() {
			<-legacyServer.RequestProcessShutdownChan()
			interrupt.Request()
		}()
	}
	select {
	case <-cx.WalletKill:
		log.L.Warn("wallet killswitch activated")
		if legacyServer != nil {
			log.L.Warn("stopping wallet RPC server")
			legacyServer.Stop()
			log.L.Info("stopped wallet RPC server")
		}
		log.L.Info("wallet shutdown from killswitch complete")
		cx.WaitGroup.Done()
		return
		// <-legacyServer.RequestProcessShutdownChan()
	case <-interrupt.HandlersDone:
	}
	log.L.Info("wallet shutdown complete")
	cx.WaitGroup.Done()
	return
}

func ReadCAFile(config *pod.Config) []byte {
	// Read certificate file if TLS is not disabled.
	var certs []byte
	if *config.TLS {
		var err error
		certs, err = ioutil.ReadFile(*config.CAFile)
		if err != nil {
			log.L.Error("cannot open CA file:", err)
			// If there's an error reading the CA file, continue
			// with nil certs and without the client connection.
			certs = nil
		}
	} else {
		log.L.Info("chain server RPC TLS is disabled")
	}
	return certs
}

// rpcClientConnectLoop continuously attempts a connection to the consensus RPC server.
// When a connection is established, the client is used to sync the loaded wallet, either immediately or when loaded at
// a later time.
//
// The legacy RPC is optional. If set, the connected RPC client will be associated with the server for RPC pass-through
// and to enable additional methods.
func rpcClientConnectLoop(cx *conte.Xt, legacyServer *legacy.Server,
	loader *wallet.Loader) {
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
		// 	netDir := networkDir(cx.PodConfig.AppDataDir.Value, ActiveNet.Params)
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
		cc, err = startChainRPC(cx.Config, cx.ActiveNet, certs)
		if err != nil {
			log.L.Error(
				"unable to open connection to consensus RPC server:", err)
			continue
		}
		cx.ChainClient = cc
		cx.ChainClientReady.Store(true)
		chainClient = cc
		// }
		// Rather than inlining this logic directly into the loader
		// callback, a function variable is used to avoid running any of
		// this after the client disconnects by setting it to nil.  This
		// prevents the callback from associating a wallet loaded at a
		// later time with a client that has already disconnected.  A
		// mutex is used to make this concurrent safe.
		associateRPCClient := func(w *wallet.Wallet) {
			if w != nil {
				w.SynchronizeRPC(chainClient)
			}
			if legacyServer != nil {
				legacyServer.SetChainServer(chainClient)
			}
		}
		mu := new(sync.Mutex)
		loader.RunAfterLoad(func(w *wallet.Wallet) {
			mu.Lock()
			associate := associateRPCClient
			mu.Unlock()
			if associate != nil {
				associate(w)
			}
		})
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

// startChainRPC opens a RPC client connection to a pod server for blockchain
// services.  This function uses the RPC options from the global config and
// there is no recovery in case the server is not available or if there is an
// authentication error.  Instead, all requests to the client will simply error.
func startChainRPC(config *pod.Config, activeNet *netparams.Params, certs []byte) (*chain.RPCClient, error) {
	log.L.Tracef(
		"attempting RPC client connection to %v, TLS: %s",
		*config.RPCConnect, fmt.Sprint(*config.TLS),
	)
	rpcC, err := chain.NewRPCClient(activeNet, *config.RPCConnect,
		*config.Username, *config.Password, certs, !*config.TLS, 0)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	err = rpcC.Start()
	return rpcC, err
}
