package node

import (
	"io"
	"os"
	"path/filepath"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/db/blockdb"
)

// dirEmpty returns whether or not the specified directory path is empty
func dirEmpty(dirPath string) (bool, error) {
	f, e := os.Open(dirPath)
	if e != nil {
		return false, e
	}
	defer func() {
		if e := f.Close(); dbg.Chk(e) {
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
func doUpgrades(cx *conte.Xt) (e error) {
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
func upgradeDBPathNet(cx *conte.Xt, oldDbPath, netName string) (e error) {
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
		newDbRoot := filepath.Join(filepath.Dir(*cx.Config.DataDir), netName)
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
func upgradeDBPaths(cx *conte.Xt) (e error) {
	// Prior to version 0.2.0 the databases were in the "db" directory and their names were suffixed by "testnet" and
	// "regtest" for their respective networks. Chk for the old database and update it to the new path introduced with
	// version 0.2.0 accordingly.
	oldDbRoot := filepath.Join(oldPodHomeDir(), "db")
	e = upgradeDBPathNet(cx, filepath.Join(oldDbRoot, "pod.db"), "mainnet")
	if e != nil {
		dbg.Ln(err)
	}
	e = upgradeDBPathNet(
		cx, filepath.Join(oldDbRoot, "pod_testnet.db"),
		"testnet",
	)
	if e != nil {
		dbg.Ln(err)
	}
	e = upgradeDBPathNet(
		cx, filepath.Join(oldDbRoot, "pod_regtest.db"),
		"regtest",
	)
	if e != nil {
		dbg.Ln(err)
	}
	// Remove the old db directory
	return os.RemoveAll(oldDbRoot)
}

// upgradeDataPaths moves the application data from its location prior to pod version 0.3.3 to its new location.
func upgradeDataPaths() (e error) {
	// No need to migrate if the old and new home paths are the same.
	oldHomePath := oldPodHomeDir()
	newHomePath := defaultHomeDir
	if oldHomePath == newHomePath {
		return nil
	}
	// Only migrate if the old path exists and the new one doesn't
	if apputil.FileExists(oldHomePath) && !apputil.FileExists(newHomePath) {
		// Create the new path
		inf.F(
			"migrating application home path from '%s' to '%s'",
			oldHomePath, newHomePath,
		)
		e := os.MkdirAll(newHomePath, 0700)
		if e != nil {
			return e
		}
		// Move old pod.conf into new location if needed
		oldConfPath := filepath.Join(oldHomePath, defaultConfigFilename)
		newConfPath := filepath.Join(newHomePath, defaultConfigFilename)
		if apputil.FileExists(oldConfPath) && !apputil.FileExists(newConfPath) {
			e := os.Rename(oldConfPath, newConfPath)
			if e != nil {
				return e
			}
		}
		// Move old data directory into new location if needed
		oldDataPath := filepath.Join(oldHomePath, defaultDataDirname)
		newDataPath := filepath.Join(newHomePath, defaultDataDirname)
		if apputil.FileExists(oldDataPath) && !apputil.FileExists(newDataPath) {
			e := os.Rename(oldDataPath, newDataPath)
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
			wrn.F(
				"not removing '%s' since it contains files not created by"+
					" this application you may want to manually move them or"+
					" delete them.", oldHomePath,
			)
		}
	}
	return nil
}
