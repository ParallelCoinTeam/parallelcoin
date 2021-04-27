/*Package node is a full-node Parallelcoin implementation written in Go.

The default options are sane for most users. This means pod will work 'out of the box' for most users. However, there
are also a wide variety of flags that can be used to control it.

The following section provides a usage overview which enumerates the flags. An interesting point to note is that the
long form of all of these options ( except -C/--configfile and -D --datadir) can be specified in a configuration file
that is automatically parsed when pod starts up. By default, the configuration file is located at ~/.pod/pod. conf on
POSIX-style operating systems and %LOCALAPPDATA%\pod\pod. conf on Windows. The -D (--datadir) flag, can be used to
override this location.

NAME:
   pod node - start parallelcoin full node

USAGE:
   pod node [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

COMMANDS:
     dropaddrindex  drop the address search index
     droptxindex    drop the address search index
     dropcfindex    drop the address search index

GLOBAL OPTIONS:
   --help, -h  show help
*/
package node

import (
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/p9c/qu"

	"github.com/p9c/interrupt"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/cmd/ctrl"
	"github.com/p9c/matrjoska/pkg/apputil"
	"github.com/p9c/matrjoska/pkg/chainrpc"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/database/blockdb"
	"github.com/p9c/matrjoska/pkg/indexers"
	"github.com/p9c/matrjoska/pod/state"
)

// // This enables pprof
// _ "net/http/pprof"

// winServiceMain is only invoked on Windows. It detects when pod is running as a service and reacts accordingly.
var winServiceMain func() (bool, error)

// NodeMain is the real main function for pod.
//
// The optional serverChan parameter is mainly used by the service code to be notified with the server once it is setup
// so it can gracefully stop it when requested from the service control manager.
func NodeMain(cx *state.State) (e error) {
	T.Ln("starting up node main")
	// cx.WaitGroup.Add(1)
	cx.WaitAdd()
	// enable http profiling server if requested
	if cx.Config.Profile.V() != "" {
		D.Ln("profiling requested")
		go func() {
			listenAddr := net.JoinHostPort("", cx.Config.Profile.V())
			I.Ln("profile server listening on", listenAddr)
			profileRedirect := http.RedirectHandler("/debug/pprof", http.StatusSeeOther)
			http.Handle("/", profileRedirect)
			D.Ln("profile server", http.ListenAndServe(listenAddr, nil))
		}()
	}
	// write cpu profile if requested
	if cx.Config.CPUProfile.V() != "" && os.Getenv("POD_TRACE") != "on" {
		D.Ln("cpu profiling enabled")
		var f *os.File
		f, e = os.Create(cx.Config.CPUProfile.V())
		if e != nil {
			E.Ln("unable to create cpu profile:", e)
			return
		}
		e = pprof.StartCPUProfile(f)
		if e != nil {
			D.Ln("failed to start up cpu profiler:", e)
		} else {
			defer func() {
				if e = f.Close(); E.Chk(e) {
				}
			}()
			defer pprof.StopCPUProfile()
			interrupt.AddHandler(
				func() {
					D.Ln("stopping CPU profiler")
					e = f.Close()
					if e != nil {
					}
					pprof.StopCPUProfile()
					D.Ln("finished cpu profiling", *cx.Config.CPUProfile)
				},
			)
		}
	}
	// perform upgrades to pod as new versions require it
	if e = doUpgrades(cx); E.Chk(e) {
		return
	}
	// return now if an interrupt signal was triggered
	if interrupt.Requested() {
		return nil
	}
	// load the block database
	var db database.DB
	db, e = loadBlockDB(cx)
	if e != nil {
		return
	}
	closeDb := func() {
		// ensure the database is synced and closed on shutdown
		T.Ln("gracefully shutting down the database")
		func() {
			if e = db.Close(); E.Chk(e) {
			}
		}()
	}
	defer closeDb()
	interrupt.AddHandler(closeDb)
	// return now if an interrupt signal was triggered
	if interrupt.Requested() {
		return nil
	}
	// drop indexes and exit if requested.
	//
	// NOTE: The order is important here because dropping the tx index also drops the address index since it relies on
	// it
	if cx.StateCfg.DropAddrIndex {
		W.Ln("dropping address index")
		if e = indexers.DropAddrIndex(db, interrupt.ShutdownRequestChan); E.Chk(e) {
			return
		}
	}
	if cx.StateCfg.DropTxIndex {
		W.Ln("dropping transaction index")
		if e = indexers.DropTxIndex(db, interrupt.ShutdownRequestChan); E.Chk(e) {
			return
		}
	}
	if cx.StateCfg.DropCfIndex {
		W.Ln("dropping cfilter index")
		if e = indexers.DropCfIndex(db, interrupt.ShutdownRequestChan); E.Chk(e) {
			return
		}
	}
	// return now if an interrupt signal was triggered
	if interrupt.Requested() {
		return nil
	}
	mempoolUpdateChan := qu.Ts(1)
	mempoolUpdateHook := func() {
		mempoolUpdateChan.Signal()
	}
	// create server and start it
	var server *chainrpc.Node
	server, e = chainrpc.NewNode(
		cx.Config.P2PListeners.S(),
		db,
		interrupt.ShutdownRequestChan,
		state.GetContext(cx),
		mempoolUpdateHook,
	)
	if e != nil {
		E.F("unable to start server on %v: %v", cx.Config.P2PListeners.S(), e)
		return e
	}
	server.Start()
	cx.RealNode = server
	// if len(server.RPCServers) > 0 && *cx.Config.CAPI {
	// 	D.Ln("starting cAPI.....")
	// 	// chainrpc.RunAPI(server.RPCServers[0], cx.NodeKill)
	// 	// D.Ln("propagating rpc server handle (node has started)")
	// }
	// I.S(server.RPCServers)
	if len(server.RPCServers) > 0 {
		cx.RPCServer = server.RPCServers[0]
		D.Ln("sending back node")
		cx.NodeChan <- cx.RPCServer
	}
	D.Ln("starting controller")
	cx.Controller, e = ctrl.New(
		cx.Syncing,
		cx.Config,
		cx.StateCfg,
		cx.RealNode,
		cx.RPCServer.Cfg.ConnMgr,
		mempoolUpdateChan,
		uint64(cx.Config.UUID.V()),
		cx.KillAll,
		cx.RealNode.StartController, cx.RealNode.StopController,
	)
	go cx.Controller.Run()
	cx.Controller.Start()
	D.Ln("controller started")
	once := true
	gracefulShutdown := func() {
		if !once {
			return
		}
		if once {
			once = false
		}
		D.Ln("gracefully shutting down the server...")
		D.Ln("stopping controller")
		cx.Controller.Shutdown()
		D.Ln("stopping server")
		e := server.Stop()
		if e != nil {
			W.Ln("failed to stop server", e)
		}
		server.WaitForShutdown()
		I.Ln("server shutdown complete")
		log.LogChanDisabled.Store(true)
		cx.WaitDone()
		cx.KillAll.Q()
		cx.NodeKill.Q()
	}
	D.Ln("adding interrupt handler for node")
	interrupt.AddHandler(gracefulShutdown)
	// Wait until the interrupt signal is received from an OS signal or shutdown is requested through one of the
	// subsystems such as the RPC server.
	select {
	case <-cx.NodeKill.Wait():
		D.Ln("NodeKill")
		if !interrupt.Requested() {
			interrupt.Request()
		}
		break
	case <-cx.KillAll.Wait():
		D.Ln("KillAll")
		if !interrupt.Requested() {
			interrupt.Request()
		}
		break
	}
	gracefulShutdown()
	return nil
}

// loadBlockDB loads (or creates when needed) the block database taking into account the selected database backend and
// returns a handle to it. It also additional logic such warning the user if there are multiple databases which consume
// space on the file system and ensuring the regression test database is clean when in regression test mode.
func loadBlockDB(cx *state.State) (db database.DB, e error) {
	// The memdb backend does not have a file path associated with it, so handle it uniquely. We also don't want to
	// worry about the multiple database type warnings when running with the memory database.
	if cx.Config.DbType.V() == "memdb" {
		I.Ln("creating block database in memory")
		if db, e = database.Create(cx.Config.DbType.V()); state.E.Chk(e) {
			return nil, e
		}
		return db, nil
	}
	warnMultipleDBs(cx)
	// The database name is based on the database type.
	dbPath := state.BlockDb(cx, cx.Config.DbType.V(), blockdb.NamePrefix)
	// The regression test is special in that it needs a clean database for each
	// run, so remove it now if it already exists.
	e = removeRegressionDB(cx, dbPath)
	if e != nil {
		D.Ln("failed to remove regression db:", e)
	}
	I.F("loading block database from '%s'", dbPath)
	I.Ln(database.SupportedDrivers())
	if db, e = database.Open(cx.Config.DbType.V(), dbPath, cx.ActiveNet.Net); E.Chk(e) {
		T.Ln(e) // return the error if it's not because the database doesn't exist
		if dbErr, ok := e.(database.DBError); !ok || dbErr.ErrorCode !=
			database.ErrDbDoesNotExist {
			return nil, e
		}
		// create the db if it does not exist
		e = os.MkdirAll(cx.Config.DataDir.V(), 0700)
		if e != nil {
			return nil, e
		}
		db, e = database.Create(cx.Config.DbType.V(), dbPath, cx.ActiveNet.Net)
		if e != nil {
			return nil, e
		}
	}
	T.Ln("block database loaded")
	return db, nil
}

// removeRegressionDB removes the existing regression test database if running
// in regression test mode and it already exists.
func removeRegressionDB(cx *state.State, dbPath string) (e error) {
	// don't do anything if not in regression test mode
	if !((cx.Config.Network.V())[0] == 'r') {
		return nil
	}
	// remove the old regression test database if it already exists
	fi, e := os.Stat(dbPath)
	if e == nil {
		I.F("removing regression test database from '%s' %s", dbPath)
		if fi.IsDir() {
			if e = os.RemoveAll(dbPath); E.Chk(e) {
				return e
			}
		} else {
			if e = os.Remove(dbPath); E.Chk(e) {
				return e
			}
		}
	}
	return nil
}

// warnMultipleDBs shows a warning if multiple block database types are
// detected. This is not a situation most users want. It is handy for
// development however to support multiple side-by-side databases.
func warnMultipleDBs(cx *state.State) {
	// This is intentionally not using the known db types which depend on the
	// database types compiled into the binary since we want to detect legacy db
	// types as well.
	dbTypes := []string{"ffldb", "leveldb", "sqlite"}
	duplicateDbPaths := make([]string, 0, len(dbTypes)-1)
	for _, dbType := range dbTypes {
		if dbType == cx.Config.DbType.V() {
			continue
		}
		// store db path as a duplicate db if it exists
		dbPath := state.BlockDb(cx, dbType, blockdb.NamePrefix)
		if apputil.FileExists(dbPath) {
			duplicateDbPaths = append(duplicateDbPaths, dbPath)
		}
	}
	// warn if there are extra databases
	if len(duplicateDbPaths) > 0 {
		selectedDbPath := state.BlockDb(cx, cx.Config.DbType.V(), blockdb.NamePrefix)
		W.F(
			"\nThere are multiple block chain databases using different"+
				" database types.\nYou probably don't want to waste disk"+
				" space by having more than one."+
				"\nYour current database is located at [%v]."+
				"\nThe additional database is located at %v",
			selectedDbPath,
			duplicateDbPaths,
		)
	}
}

// dirEmpty returns whether or not the specified directory path is empty
func dirEmpty(dirPath string) (bool, error) {
	f, e := os.Open(dirPath)
	if e != nil {
		return false, e
	}
	defer func() {
		if e = f.Close(); E.Chk(e) {
		}
	}()
	// Read the names of a max of one entry from the directory. When the directory is empty, an io.EOF error will be
	// returned, so allow it.
	names, e := f.Readdirnames(1)
	if e != nil && e != io.EOF {
		return false, e
	}
	return len(names) == 0, nil
}

// doUpgrades performs upgrades to pod as new versions require it
func doUpgrades(cx *state.State) (e error) {
	e = upgradeDBPaths(cx)
	if e != nil {
		return e
	}
	return upgradeDataPaths()
}

// oldPodHomeDir returns the OS specific home directory pod used prior to version 0.3.3. This has since been replaced
// with util.AppDataDir but this function is still provided for the automatic upgrade path.
func oldPodHomeDir() string {
	// Search for Windows APPDATA first. This won't exist on POSIX OSes
	appData := os.Getenv("APPDATA")
	if appData != "" {
		return filepath.Join(appData, "pod")
	}
	// Fall back to standard HOME directory that works for most POSIX OSes
	home := os.Getenv("HOME")
	if home != "" {
		return filepath.Join(home, ".pod")
	}
	// In the worst case, use the current directory
	return "."
}

// upgradeDBPathNet moves the database for a specific network from its location prior to pod version 0.2.0 and uses
// heuristics to ascertain the old database type to rename to the new format.
func upgradeDBPathNet(cx *state.State, oldDbPath, netName string) (e error) {
	// Prior to version 0.2.0, the database was named the same thing for both sqlite and leveldb. Use heuristics to
	// figure out the type of the database and move it to the new path and name introduced with version 0.2.0
	// accordingly.
	fi, e := os.Stat(oldDbPath)
	if e == nil {
		oldDbType := "sqlite"
		if fi.IsDir() {
			oldDbType = "leveldb"
		}
		// The new database name is based on the database type and resides in a directory named after the network type.
		newDbRoot := filepath.Join(filepath.Dir(cx.Config.DataDir.V()), netName)
		newDbName := blockdb.NamePrefix + "_" + oldDbType
		if oldDbType == "sqlite" {
			newDbName = newDbName + ".db"
		}
		newDbPath := filepath.Join(newDbRoot, newDbName)
		// Create the new path if needed
		//
		e = os.MkdirAll(newDbRoot, 0700)
		if e != nil {
			return e
		}
		// Move and rename the old database
		//
		e := os.Rename(oldDbPath, newDbPath)
		if e != nil {
			return e
		}
	}
	return nil
}

// upgradeDBPaths moves the databases from their locations prior to pod version 0.2.0 to their new locations
func upgradeDBPaths(cx *state.State) (e error) {
	// Prior to version 0.2.0 the databases were in the "db" directory and their names were suffixed by "testnet" and
	// "regtest" for their respective networks. Chk for the old database and update it to the new path introduced with
	// version 0.2.0 accordingly.
	oldDbRoot := filepath.Join(oldPodHomeDir(), "db")
	e = upgradeDBPathNet(cx, filepath.Join(oldDbRoot, "pod.db"), "mainnet")
	if e != nil {
		D.Ln(e)
	}
	e = upgradeDBPathNet(
		cx, filepath.Join(oldDbRoot, "pod_testnet.db"),
		"testnet",
	)
	if e != nil {
		D.Ln(e)
	}
	e = upgradeDBPathNet(
		cx, filepath.Join(oldDbRoot, "pod_regtest.db"),
		"regtest",
	)
	if e != nil {
		D.Ln(e)
	}
	// Remove the old db directory
	return os.RemoveAll(oldDbRoot)
}

// upgradeDataPaths moves the application data from its location prior to pod version 0.3.3 to its new location.
func upgradeDataPaths() (e error) {
	// No need to migrate if the old and new home paths are the same.
	oldHomePath := oldPodHomeDir()
	newHomePath := constant.DefaultHomeDir
	if oldHomePath == newHomePath {
		return nil
	}
	// Only migrate if the old path exists and the new one doesn't
	if apputil.FileExists(oldHomePath) && !apputil.FileExists(newHomePath) {
		// Create the new path
		I.F(
			"migrating application home path from '%s' to '%s'",
			oldHomePath, newHomePath,
		)
		e := os.MkdirAll(newHomePath, 0700)
		if e != nil {
			return e
		}
		// Move old pod.conf into new location if needed
		oldConfPath := filepath.Join(oldHomePath, constant.DefaultConfigFilename)
		newConfPath := filepath.Join(newHomePath, constant.DefaultConfigFilename)
		if apputil.FileExists(oldConfPath) && !apputil.FileExists(newConfPath) {
			e = os.Rename(oldConfPath, newConfPath)
			if e != nil {
				return e
			}
		}
		// Move old data directory into new location if needed
		oldDataPath := filepath.Join(oldHomePath, constant.DefaultDataDirname)
		newDataPath := filepath.Join(newHomePath, constant.DefaultDataDirname)
		if apputil.FileExists(oldDataPath) && !apputil.FileExists(newDataPath) {
			e = os.Rename(oldDataPath, newDataPath)
			if e != nil {
				return e
			}
		}
		// Remove the old home if it is empty or show a warning if not
		ohpEmpty, e := dirEmpty(oldHomePath)
		if e != nil {
			return e
		}
		if ohpEmpty {
			e := os.Remove(oldHomePath)
			if e != nil {
				return e
			}
		} else {
			W.F(
				"not removing '%s' since it contains files not created by"+
					" this application you may want to manually move them or"+
					" delete them.", oldHomePath,
			)
		}
	}
	return nil
}
