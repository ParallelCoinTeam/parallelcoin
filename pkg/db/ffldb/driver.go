package ffldb

import (
	"fmt"
	"github.com/stalker-loki/app/slog"

	"github.com/p9c/pod/pkg/chain/wire"
	database "github.com/p9c/pod/pkg/db"
)

const (
	dbType = "ffldb"
)

// parseArgs parses the arguments from the database Open/Create methods.
func parseArgs(funcName string, args ...interface{}) (dbPath string, network wire.BitcoinNet, err error) {
	if len(args) != 2 {
		err = fmt.Errorf("invalid arguments to %s.%s -- expected database path and block network", dbType, funcName)
		slog.Debug(err)
		return
	}
	var ok bool
	if dbPath, ok = args[0].(string); !ok {
		err = fmt.Errorf("first argument to %s.%s is invalid -- expected database path string", dbType, funcName)
		slog.Debug(err)
		return
	}
	if network, ok = args[1].(wire.BitcoinNet); !ok {
		err = fmt.Errorf("second argument to %s.%s is invalid -- expected block network", dbType, funcName)
		slog.Debug(err)
		return
	}
	return
}

// openDBDriver is the callback provided during driver registration that opens an existing database for use.
func openDBDriver(args ...interface{}) (d database.DB, err error) {
	var dbPath string
	var network wire.BitcoinNet
	if dbPath, network, err = parseArgs("Open", args...); slog.Check(err) {
		return
	}
	return openDB(dbPath, network, false)
}

// createDBDriver is the callback provided during driver registration that creates, initializes, and opens a database for use.
func createDBDriver(args ...interface{}) (d database.DB, err error) {
	var dbPath string
	var network wire.BitcoinNet
	if dbPath, network, err = parseArgs("Create", args...); slog.Check(err) {
		return
	}
	return openDB(dbPath, network, true)
}
func init() {
	// Register the driver.
	driver := database.Driver{
		DbType: dbType,
		Create: createDBDriver,
		Open:   openDBDriver,
	}
	if err := database.RegisterDriver(driver); err != nil {
		panic(fmt.Sprintf("Failed to regiser database driver '%s': %v",
			dbType, err))
	}
}
