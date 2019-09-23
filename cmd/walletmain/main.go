package walletmain

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	// This enables pprof
	_ "net/http/pprof"
	"sync"

	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining/addresses"
	"github.com/p9c/pod/pkg/log"
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
func Main(config *pod.Config, stateCfg *state.Config,
	activeNet *netparams.Params,
	walletChan chan *wallet.Wallet, killswitch chan struct{},
	wg *sync.WaitGroup) error {
	log.INFO("starting wallet")
	wg.Add(1)
	if activeNet.Name == "testnet" {
		fork.IsTestnet = true
	}
	if *config.Profile != "" {
		go func() {
			listenAddr := net.JoinHostPort("127.0.0.1", *config.Profile)
			log.INFO("profile server listening on", listenAddr)
			profileRedirect := http.RedirectHandler("/debug/pprof",
				http.StatusSeeOther)
			http.Handle("/", profileRedirect)
			fmt.Println(http.ListenAndServe(listenAddr, nil))
		}()
	}
	dbPath := *config.DataDir + slash + activeNet.Params.Name
	loader := wallet.NewLoader(activeNet, dbPath, 250)
	// Create and start HTTP server to serve wallet client connections.
	// This will be updated with the wallet and chain server RPC client
	// created below after each is created.
	log.TRACE("starting RPC servers")
	rpcS, legacyServer, err := startRPCServers(config, stateCfg, activeNet,
		loader)
	if err != nil {
		log.ERROR("unable to create RPC servers:", err)
		return err
	}
	loader.RunAfterLoad(func(w *wallet.Wallet) {
		log.WARN("starting wallet RPC services", w != nil)
		startWalletRPCServices(w, rpcS, legacyServer)
		addresses.RefillMiningAddresses(w, config, stateCfg)
	})
	if !*config.NoInitialLoad {
		log.TRACE("starting rpc client connection handler")
		go rpcClientConnectLoop(config, activeNet, legacyServer, loader)
		// Create and start chain RPC client so it's ready to connect to
		// the wallet when loaded later.
		log.WARN("loading database")
		// Load the wallet database.  It must have been created already
		// or this will return an appropriate error.
		var w *wallet.Wallet
		w, err = loader.OpenExistingWallet([]byte(*config.WalletPass),
			true)
		log.TRACE("wallet", w)
		if err != nil {
			log.ERROR(err)
			return err
		}
		loader.Wallet = w
		log.TRACE("sending back wallet")
		walletChan <- w
	}
	log.TRACE("adding interrupt handler to unload wallet")
	// Add interrupt handlers to shutdown the various process components
	// before exiting.  Interrupt handlers run in LIFO order, so the wallet
	// (which should be closed last) is added first.
	interrupt.AddHandler(func() {
		err := loader.UnloadWallet()
		if err != nil && err != wallet.ErrNotLoaded {
			log.ERROR("failed to close wallet:", err)
		}
	})
	if rpcS != nil {
		interrupt.AddHandler(func() {
			// TODO: Does this need to wait for the grpc server to
			// finish up any requests?
			log.WARN("stopping RPC server")
			rpcS.Stop()
			stateCfg.DiscoveryUpdate("experimentalrpc", "")
			log.INFO("RPC server shutdown")
		})
	}
	if legacyServer != nil {
		interrupt.AddHandler(func() {
			log.TRACE("stopping wallet RPC server")
			stateCfg.DiscoveryUpdate("walletrpc", "")
			legacyServer.Stop()
			log.TRACE("wallet RPC server shutdown")
		})
		go func() {
			<-legacyServer.RequestProcessShutdownChan()
			interrupt.Request()
		}()
	}
	select {
	case <-killswitch:
		log.WARN("wallet killswitch activated")
		if legacyServer != nil {
			log.WARN("stopping wallet RPC server")
			stateCfg.DiscoveryUpdate("walletrpc", "")
			legacyServer.Stop()
			log.INFO("stopped wallet RPC server")
		}
		if rpcS != nil {
			log.WARN("stopping RPC server")
			stateCfg.DiscoveryUpdate("experimentalrpc", "")
			rpcS.Stop()
			log.INFO("RPC server shutdown")
			log.INFO("unloading wallet")
			err := loader.UnloadWallet()
			if err != nil && err != wallet.ErrNotLoaded {
				log.ERROR("failed to close wallet:", err)
			}
		}
		log.INFO("wallet shutdown from killswitch complete")
		wg.Done()
		return nil
		// <-legacyServer.RequestProcessShutdownChan()
	case <-interrupt.HandlersDone:
	}
	log.INFO("wallet shutdown complete")
	wg.Done()
	return nil
}

func ReadCAFile(config *pod.Config) []byte {
	// Read certificate file if TLS is not disabled.
	var certs []byte
	if *config.TLS {
		var err error
		certs, err = ioutil.ReadFile(*config.CAFile)
		if err != nil {
			log.WARN("cannot open CA file:", err)
			// If there's an error reading the CA file, continue
			// with nil certs and without the client connection.
			certs = nil
		}
	} else {
		log.INFO("chain server RPC TLS is disabled")
	}
	return certs
}

// rpcClientConnectLoop continuously attempts a connection to the consensus
// RPC server.
// When a connection is established,
// the client is used to sync the loaded wallet,
// either immediately or when loaded at a later time.
//
// The legacy RPC is optional. If set,
// the connected RPC client will be associated with the server for RPC
// pass-through and to enable additional methods.
func rpcClientConnectLoop(config *pod.Config, activeNet *netparams.Params,
	legacyServer *legacy.Server, loader *wallet.Loader) {
	// var certs []byte
	// if !cx.PodConfig.UseSPV {
	certs := ReadCAFile(config)
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
		chainClient, err = startChainRPC(config, activeNet, certs)
		if err != nil {
			log.ERROR(
				"unable to open connection to consensus RPC server:", err)
			continue
		}
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
	log.TRACEF(
		"attempting RPC client connection to %v, TLS: %s, %s",
		*config.RPCConnect, fmt.Sprint(*config.TLS),
	)
	rpcC, err := chain.NewRPCClient(activeNet, *config.RPCConnect,
		*config.Username, *config.Password, certs, !*config.TLS, 0)
	if err != nil {
		return nil, err
	}
	err = rpcC.Start()
	return rpcC, err
}
