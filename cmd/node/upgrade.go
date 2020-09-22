package node

import (
	"github.com/stalker-loki/app/slog"
	"io"
	"os"
	"path/filepath"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/db/blockdb"
)

// dirEmpty returns whether or not the specified directory path is empty
func dirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		slog.Error(err)
		return false, err
	}
	defer f.Close()
	// Read the names of a max of one entry from the directory.
	// When the directory is empty, an io.EOF error will be returned,
	// so allow it.
	names, err := f.Readdirnames(1)
	if err != nil && err != io.EOF {
		slog.Error(err)
		return false, err
	}
	return len(names) == 0, nil
}

// doUpgrades performs upgrades to pod as new versions require it
func doUpgrades(cx *conte.Xt) (err error) {
	if err = upgradeDBPaths(cx); slog.Check(err) {
		return
	}
	return upgradeDataPaths()
}

// oldPodHomeDir returns the OS specific home directory pod used prior to
// version 0.3.3.
// This has since been replaced with util.AppDataDir but this function is still
// provided for the automatic upgrade path.
func oldPodHomeDir() string {
	// Search for Windows APPDATA first.  This won't exist on POSIX OSes
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

// upgradeDBPathNet moves the database for a specific network from its
// location prior to pod version 0.2.0 and uses heuristics to ascertain the old
// database type to rename to the new format.
func upgradeDBPathNet(cx *conte.Xt, oldDbPath, netName string) (err error) {
	// Prior to version 0.2.0,
	// the database was named the same thing for both sqlite and leveldb.
	// Use heuristics to figure out the type of the database and move it to
	// the new path and name introduced with version 0.2.0 accordingly.
	var fi os.FileInfo
	if fi, err = os.Stat(oldDbPath); slog.Check(err) {
		oldDbType := "sqlite"
		if fi.IsDir() {
			oldDbType = "leveldb"
		}
		// The new database name is based on the database type and resides in
		// a directory named after the network type.
		newDbRoot := filepath.Join(filepath.Dir(*cx.Config.DataDir), netName)
		newDbName := blockdb.NamePrefix + "_" + oldDbType
		if oldDbType == "sqlite" {
			newDbName = newDbName + ".db"
		}
		newDbPath := filepath.Join(newDbRoot, newDbName)
		// Create the new path if needed
		//
		if err = os.MkdirAll(newDbRoot, 0700); slog.Check(err) {
			return err
		}
		// Move and rename the old database
		//
		if err = os.Rename(oldDbPath, newDbPath); slog.Check(err) {
			return
		}
	}
	return
}

// upgradeDBPaths moves the databases from their locations prior to pod
// version 0.2.0 to their new locations
//
func upgradeDBPaths(cx *conte.Xt) (err error) {
	// Prior to version 0.2.0 the databases were in the "db" directory and
	// their names were suffixed by "testnet" and "regtest" for their
	// respective networks.  Check for the old database and update it
	// to the new path introduced with version 0.2.0 accordingly.
	oldDbRoot := filepath.Join(oldPodHomeDir(), "db")
	if err = upgradeDBPathNet(cx, filepath.Join(oldDbRoot, "pod.db"), "mainnet"); slog.Check(err) {
	}
	if err = upgradeDBPathNet(cx, filepath.Join(oldDbRoot, "pod_testnet.db"), "testnet"); slog.Check(err) {
	}
	if err = upgradeDBPathNet(cx, filepath.Join(oldDbRoot, "pod_regtest.db"), "regtest"); slog.Check(err) {
	}
	// Remove the old db directory
	//
	return os.RemoveAll(oldDbRoot)
}

// upgradeDataPaths moves the application data from its location prior to pod
// version 0.3.3 to its new location.
func upgradeDataPaths() (err error) {
	// No need to migrate if the old and new home paths are the same.
	oldHomePath := oldPodHomeDir()
	newHomePath := DefaultHomeDir
	if oldHomePath == newHomePath {
		return
	}
	// Only migrate if the old path exists and the new one doesn't
	if apputil.FileExists(oldHomePath) && !apputil.FileExists(newHomePath) {
		// Create the new path
		slog.Infof("migrating application home path from '%s' to '%s'",
			oldHomePath, newHomePath)
		if err = os.MkdirAll(newHomePath, 0700); slog.Check(err) {
			return
		}
		// Move old pod.conf into new location if needed
		oldConfPath := filepath.Join(oldHomePath, DefaultConfigFilename)
		newConfPath := filepath.Join(newHomePath, DefaultConfigFilename)
		if apputil.FileExists(oldConfPath) && !apputil.FileExists(newConfPath) {
			if err = os.Rename(oldConfPath, newConfPath); slog.Check(err) {
				return
			}
		}
		// Move old data directory into new location if needed
		oldDataPath := filepath.Join(oldHomePath, DefaultDataDirname)
		newDataPath := filepath.Join(newHomePath, DefaultDataDirname)
		if apputil.FileExists(oldDataPath) && !apputil.FileExists(newDataPath) {
			if err = os.Rename(oldDataPath, newDataPath); slog.Check(err) {
				return
			}
		}
		var ohpEmpty bool
		// Remove the old home if it is empty or show a warning if not
		if ohpEmpty, err = dirEmpty(oldHomePath); slog.Check(err) {
			return
		}
		if ohpEmpty {
			if err = os.Remove(oldHomePath); slog.Check(err) {
				return
			}
		} else {
			slog.Warnf("not removing '%s' since it contains files not created by"+
				" this application you may want to manually move them or"+
				" delete them.", oldHomePath)
		}
	}
	return
}
