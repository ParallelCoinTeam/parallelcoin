package bdb

import (
	"fmt"
	"github.com/stalker-loki/app/slog"

	"github.com/p9c/pod/pkg/db/walletdb"
)

const (
	dbType = "bdb"
)

// parseArgs parses the arguments from the walletdb Open/Create methods.
func parseArgs(funcName string, args ...interface{}) (dbPath string, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("invalid arguments to %s.%s -- expected database path", dbType, funcName)
		slog.Debug(err)
		return
	}
	var ok bool
	if dbPath, ok = args[0].(string); !ok {
		err = fmt.Errorf("first argument to %s.%s is invalid -- expected database path string", dbType, funcName)
		slog.Debug(err)
	}
	return
}

// openDBDriver is the callback provided during driver registration that opens
// an existing database for use.
func openDBDriver(args ...interface{}) (d walletdb.DB, err error) {
	var dbPath string
	if dbPath, err = parseArgs("Open", args...); slog.Check(err) {
		return
	}
	return openDB(dbPath, false)
}

// createDBDriver is the callback provided during driver registration that
// creates, initializes, and opens a database for use.
func createDBDriver(args ...interface{}) (d walletdb.DB, err error) {
	dbPath, err := parseArgs("Create", args...)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	return openDB(dbPath, true)
}
func init() {
	// Register the driver.
	driver := walletdb.Driver{
		DbType: dbType,
		Create: createDBDriver,
		Open:   openDBDriver,
	}
	if err := walletdb.RegisterDriver(driver); err != nil {
		panic(fmt.Sprintf("Failed to regiser database driver '%s': %v",
			dbType, err))
	}
}
