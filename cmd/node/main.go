package node

import (
	"net"
	"net/http"
	// // This enables pprof
	// _ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"

	"github.com/p9c/pod/cmd/node/blockdb"
	"github.com/p9c/pod/pkg/kopachctrl"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/node/path"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/cmd/node/version"
	indexers "github.com/p9c/pod/pkg/chain/index"
	"github.com/p9c/pod/pkg/conte"
	database "github.com/p9c/pod/pkg/db"
	"github.com/p9c/pod/pkg/util/interrupt"
)

// var StateCfg = new(state.Config)
// var cfg *pod.Config

// winServiceMain is only invoked on Windows.
// It detects when pod is running as a service and reacts accordingly.
// nolint
var winServiceMain func() (bool, error)

// Main is the real main function for pod.
// It is necessary to work around the fact that deferred functions do not run
// when os.Exit() is called.
// The optional serverChan parameter is mainly used by the service code to be
// notified with the server once it is setup so it can gracefully stop it
// when requested from the service control manager.
//  - shutdownchan can be used to wait for the node to shut down
//  - killswitch can be closed to shut the node down
func Main(cx *conte.Xt, shutdownChan chan struct{}) (err error) {
	L.Trace("starting up node main")
	cx.WaitGroup.Add(1)

	// show version at startup
	L.Info("version", version.Version())
	// enable http profiling server if requested
	if *cx.Config.Profile != "" {
		L.Debug("profiling requested")
		go func() {
			listenAddr := net.JoinHostPort("",
				*cx.Config.Profile)
			L.Info("profile server listening on", listenAddr)
			profileRedirect := http.RedirectHandler(
				"/debug/pprof", http.StatusSeeOther)
			http.Handle("/", profileRedirect)
			L.Error("profile server", http.ListenAndServe(listenAddr, nil))
		}()
	}
	// write cpu profile if requested
	if *cx.Config.CPUProfile != "" {
		L.Warn("cpu profiling enabled")
		var f *os.File
		f, err = os.Create(*cx.Config.CPUProfile)
		if err != nil {
			L.Error("unable to create cpu profile:", err)
			return
		}
		e := pprof.StartCPUProfile(f)
		if e != nil {
			L.Warn("failed to start up cpu profiler:", e)
		} else {
			// go func() {
			//	L.Error(http.ListenAndServe(":6060", nil))
			// }()
			interrupt.AddHandler(func() {
				L.Warn("stopping CPU profiler")
				err := f.Close()
				if err != nil {
					L.Error(err)
				}
				pprof.StopCPUProfile()
				L.Warn("finished cpu profiling", *cx.Config.CPUProfile)
			})
		}
	}
	// perform upgrades to pod as new versions require it
	if err = doUpgrades(cx); err != nil {
		L.Error(err)
		return
	}
	// return now if an interrupt signal was triggered
	if interrupt.Requested() {
		return nil
	}
	// load the block database
	var db database.DB
	db, err = loadBlockDB(cx)
	if err != nil {
		L.Error(err)
		return
	}
	defer func() {
		// ensure the database is sync'd and closed on shutdown
		L.Trace("gracefully shutting down the database")
		db.Close()
		time.Sleep(time.Second / 4)
	}()
	// return now if an interrupt signal was triggered
	if interrupt.Requested() {
		return nil
	}
	// drop indexes and exit if requested.
	// NOTE: The order is important here because dropping the
	// tx index also drops the address index since it relies on it
	if cx.StateCfg.DropAddrIndex {
		L.Warn("dropping address index")
		if err = indexers.DropAddrIndex(db,
			interrupt.ShutdownRequestChan); err != nil {
			L.Error(err)
			return
		}
	}
	if cx.StateCfg.DropTxIndex {
		L.Warn("dropping transaction index")
		if err = indexers.DropTxIndex(db,
			interrupt.ShutdownRequestChan); err != nil {
			L.Error(err)
			return
		}
	}
	if cx.StateCfg.DropCfIndex {
		L.Warn("dropping cfilter index")
		if err = indexers.DropCfIndex(db,
			interrupt.ShutdownRequestChan); err != nil {
			L.Error(err)
			if err != nil {
				L.Error(err)
				return
			}
		}
	}
	// return now if an interrupt signal was triggered
	if interrupt.Requested() {
		return nil
	}
	// create server and start it
	server, err := rpc.NewNode(*cx.Config.Listeners, db,
		interrupt.ShutdownRequestChan, conte.GetContext(cx))
	if err != nil {
		L.Errorf("unable to start server on %v: %v",
			*cx.Config.Listeners, err)
		return err
	}

	server.Start()
	cx.RealNode = server
	if len(server.RPCServers) > 0 {
		rpc.RunAPI(server.RPCServers[0], cx.NodeKill)
		L.Trace("propagating rpc server handle (node has started)")
		cx.RPCServer = server.RPCServers[0]
		if cx.NodeChan != nil {
			L.Trace("sending back node")
			cx.NodeChan <- server.RPCServers[0]
		}
	}
	// set up interrupt shutdown handlers to stop servers
	stopController := kopachctrl.Run(cx)
	cx.Controller.Store(true)
	gracefulShutdown := func() {
		L.Info("gracefully shutting down the server...")
		// server.CPUMiner.Stop()
		L.Debug("stopping controller")
		e := server.Stop()
		if e != nil {
			L.Warn("failed to stop server", e)
		}
		if stopController != nil {
			close(stopController)
		}
		server.WaitForShutdown()
		L.Info("server shutdown complete")
		cx.WaitGroup.Done()
	}
	if shutdownChan != nil {
		interrupt.AddHandler(func() {
			L.Debug("node.Main interrupt")
			if cx.Controller.Load() {
				gracefulShutdown()
			}
			close(shutdownChan)
		})
	}
	// interrupt.AddHandler(gracefulShutdown)

	// Wait until the interrupt signal is received from an OS signal or
	// shutdown is requested through one of the subsystems such as the
	// RPC server.
	select {
	case <-cx.NodeKill:
		gracefulShutdown()
		// case <-interrupt.HandlersDone:
		//	wg.Done()
	}
	return nil
}

// loadBlockDB loads (or creates when needed) the block database taking into
// account the selected database backend and returns a handle to it.
// It also additional logic such warning the user if there are multiple
// databases which consume space on the file system and ensuring the
// regression test database is clean when in regression test mode.
func loadBlockDB(cx *conte.Xt) (database.DB, error) {
	// The memdb backend does not have a file path associated with it,
	// so handle it uniquely.
	// We also don't want to worry about the multiple database type
	// warnings when running with the memory database.
	if *cx.Config.DbType == "memdb" {
		L.Info("creating block database in memory")
		db, err := database.Create(*cx.Config.DbType)
		if err != nil {
			L.Error(err)
			return nil, err
		}
		return db, nil
	}
	warnMultipleDBs(cx)
	// The database name is based on the database type.
	dbPath := path.BlockDb(cx, *cx.Config.DbType, blockdb.NamePrefix)
	// The regression test is special in that it needs a clean database
	// for each run, so remove it now if it already exists.
	e := removeRegressionDB(cx, dbPath)
	if e != nil {
		L.Debug("failed to remove regression db:", e)
	}
	L.Infof("loading block database from '%s'", dbPath)
	db, err := database.Open(*cx.Config.DbType, dbPath, cx.ActiveNet.Net)
	if err != nil {
		L.Trace(err) // return the error if it's not because the database doesn't exist
		if dbErr, ok := err.(database.Error); !ok || dbErr.ErrorCode !=
			database.ErrDbDoesNotExist {
			return nil, err
		}
		// create the db if it does not exist
		err = os.MkdirAll(*cx.Config.DataDir, 0700)
		if err != nil {
			L.Error(err)
			return nil, err
		}
		db, err = database.Create(*cx.Config.DbType, dbPath, cx.ActiveNet.Net)
		if err != nil {
			L.Error(err)
			return nil, err
		}
	}
	L.Trace("block database loaded")
	return db, nil
}

// removeRegressionDB removes the existing regression test database if
// running in regression test mode and it already exists.
func removeRegressionDB(cx *conte.Xt, dbPath string) error {
	// don't do anything if not in regression test mode
	if !((*cx.Config.Network)[0] == 'r') {
		return nil
	}
	// remove the old regression test database if it already exists
	fi, err := os.Stat(dbPath)
	if err == nil {
		L.Infof("removing regression test database from '%s' %s", dbPath)
		if fi.IsDir() {
			if err = os.RemoveAll(dbPath); err != nil {
				return err
			}
		} else {
			if err = os.Remove(dbPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// warnMultipleDBs shows a warning if multiple block database types are
// detected. This is not a situation most users want.
// It is handy for development however to support multiple side-by-side databases.
func warnMultipleDBs(cx *conte.Xt) {
	// This is intentionally not using the known db types which depend on the
	// database types compiled into the binary since we want to detect legacy
	// db types as well.
	dbTypes := []string{"ffldb", "leveldb", "sqlite"}
	duplicateDbPaths := make([]string, 0, len(dbTypes)-1)
	for _, dbType := range dbTypes {
		if dbType == *cx.Config.DbType {
			continue
		}
		// store db path as a duplicate db if it exists
		dbPath := path.BlockDb(cx, dbType, blockdb.NamePrefix)
		if apputil.FileExists(dbPath) {
			duplicateDbPaths = append(duplicateDbPaths, dbPath)
		}
	}
	// warn if there are extra databases
	if len(duplicateDbPaths) > 0 {
		selectedDbPath := path.BlockDb(cx, *cx.Config.DbType, blockdb.NamePrefix)
		L.Warnf(
			"\nThere are multiple block chain databases using different"+
				" database types.\nYou probably don't want to waste disk"+
				" space by having more than one."+
				"\nYour current database is located at [%v]."+
				"\nThe additional database is located at %v",
			selectedDbPath,
			duplicateDbPaths)
	}
}
