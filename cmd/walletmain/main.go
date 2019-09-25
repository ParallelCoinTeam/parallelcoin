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
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util/cl"
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
	log <- cl.Info{"starting wallet", cl.Ine()}
	log <- cl.Trace{"wg+1", cl.Ine()}
	wg.Add(1)
	if activeNet.Name == "testnet" {
		fork.IsTestnet = true
	}
	if *config.Profile != "" {
		go func() {
			listenAddr := net.JoinHostPort("127.0.0.1", *config.Profile)
			log <- cl.Info{
				"profile server listening on", listenAddr,cl.Ine(),
			}
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
	log <- cl.Trace{"starting RPC servers", cl.Ine()}
	rpcS, legacyServer, err := startRPCServers(config, stateCfg, activeNet,
		loader)
	if err != nil {
		log <- cl.Error{
			"unable to create RPC servers:", err, cl.Ine()}
		return err
	}
	loader.RunAfterLoad(func(w *wallet.Wallet) {
		// log <- cl.Warn{"starting wallet RPC services", w != nil, cl.Ine()}
		// startWalletRPCServices(w, rpcS, legacyServer)
		addresses.RefillMiningAddresses(w, config, stateCfg)
	})
	if !*config.NoInitialLoad {
		log <- cl.Trace{"starting rpc client connection handler",cl.Ine()}
		go rpcClientConnectLoop(config, activeNet, legacyServer, loader)
		// Create and start chain RPC client so it's ready to connect to
		// the wallet when loaded later.
		log <- cl.Warn{"loading database", cl.Ine()}
		// Load the wallet database.  It must have been created already
		// or this will return an appropriate error.
		var w *wallet.Wallet
		w, err = loader.OpenExistingWallet([]byte(*config.WalletPass),
			true)
		log <- cl.Trace{"wallet", w, cl.Ine()}
		if err != nil {
			log <- cl.Error{err, cl.Ine()}
			return err
		}
		loader.Wallet = w
		log <- cl.Trace{"sending back wallet", cl.Ine()}
		walletChan <- w
	}
	log <- cl.Trace{"adding interrupt handler to unload wallet", cl.Ine()}
	// Add interrupt handlers to shutdown the various process components
	// before exiting.  Interrupt handlers run in LIFO order, so the wallet
	// (which should be closed last) is added first.
	interrupt.AddHandler(func() {
		err := loader.UnloadWallet()
		if err != nil && err != wallet.ErrNotLoaded {
			log <- cl.Error{
				"failed to close wallet:", err, cl.Ine()}
		}
	})
	if rpcS != nil {
		interrupt.AddHandler(func() {
			// TODO: Does this need to wait for the grpc server to
			// finish up any requests?
			log <- cl.Warn{"stopping RPC server...", cl.Ine()}
			rpcS.Stop()
			stateCfg.DiscoveryUpdate("experimentalrpc", "")
			log <- cl.Info{"RPC server shutdown", cl.Ine()}
		})
	}
	if legacyServer != nil {
		interrupt.AddHandler(func() {
			log <- cl.Trace{"stopping wallet RPC server...", cl.Ine()}
			stateCfg.DiscoveryUpdate("walletrpc", "")
			legacyServer.Stop()
			log <- cl.Trace{"wallet RPC server shutdown", cl.Ine()}
		})
		go func() {
			<-legacyServer.RequestProcessShutdownChan()
			interrupt.Request()
		}()
	}
	select {
	case <-killswitch:
		log <- cl.Warn{"wallet killswitch activated", cl.Ine()}
		if legacyServer != nil {
			log <- cl.Warn{"stopping wallet RPC server...", cl.Ine()}
			stateCfg.DiscoveryUpdate("walletrpc", "")
			legacyServer.Stop()
			log <- cl.Info{"stopped wallet RPC server", cl.Ine()}
		}
		if rpcS != nil {
			log <- cl.Warn{"stopping RPC server...", cl.Ine()}
			stateCfg.DiscoveryUpdate("experimentalrpc", "")
			rpcS.Stop()
			log <- cl.Info{"RPC server shutdown", cl.Ine()}
		}
		log <- cl.Info{"unloading wallet"}
		err := loader.UnloadWallet()
		if err != nil && err != wallet.ErrNotLoaded {
			log <- cl.Error{
				"failed to close wallet:", err, cl.Ine()}
		}
		log <- cl.Info{"wallet shutdown from killswitch complete", cl.Ine()}
		log <- cl.Error{"wg-1 3", cl.Ine()}
		wg.Done()
		return nil
		// <-legacyServer.RequestProcessShutdownChan()
	case <-interrupt.HandlersDone:
	}
	log <- cl.Info{"wallet shutdown complete", cl.Ine()}
	log <- cl.Trace{"wg-1 4", cl.Ine()}
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
			log <- cl.Warn{
				"cannot open CA file:", err,
			}
			// If there's an error reading the CA file, continue
			// with nil certs and without the client connection.
			certs = nil
		}
	} else {
		log <- cl.Info{"chain server RPC TLS is disabled", cl.Ine()}
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
func rpcClientConnectLoop(config *pod.Config, activenet *netparams.Params,
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
		chainClient, err = startChainRPC(config, activenet, certs)
		if err != nil {
			log <- cl.Error{
				"unable to open connection to consensus RPC server:", err, cl.Ine()}
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
	log <- cl.Tracef{
		"attempting RPC client connection to %v, TLS: %s, %s",
		*config.RPCConnect, fmt.Sprint(*config.TLS), cl.Ine(),
	}
	rpcc, err := chain.NewRPCClient(activeNet, *config.RPCConnect,
		*config.Username, *config.Password, certs, !*config.TLS, 0)
	if err != nil {
		return nil, err
	}
	err = rpcc.Start()
	return rpcc, err
}
