package ffldb

import (
	"fmt"
	
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/wire"
)

const (
	dbType = "ffldb"
)

// parseArgs parses the arguments from the database Open/Create methods.
func parseArgs(funcName string, args ...interface{}) (string, wire.BitcoinNet, error) {
	if len(args) != 2 {
		return "", 0, fmt.Errorf(
			"invalid arguments to %s.%s -- "+
				"expected database path and block network", dbType,
			funcName,
		)
	}
	dbPath, ok := args[0].(string)
	if !ok {
		return "", 0, fmt.Errorf(
			"first argument to %s.%s is invalid -- "+
				"expected database path string", dbType, funcName,
		)
	}
	network, ok := args[1].(wire.BitcoinNet)
	if !ok {
		return "", 0, fmt.Errorf(
			"second argument to %s.%s is invalid -- "+
				"expected block network", dbType, funcName,
		)
	}
	return dbPath, network, nil
}

// openDBDriver is the callback provided during driver registration that opens an existing database for use.
func openDBDriver(args ...interface{}) (database.DB, error) {
	dbPath, network, e := parseArgs("Open", args...)
	if e != nil {
		return nil, e
	}
	return openDB(dbPath, network, false)
}

// createDBDriver is the callback provided during driver registration that creates, initializes, and opens a database
// for use.
func createDBDriver(args ...interface{}) (database.DB, error) {
	dbPath, network, e := parseArgs("Create", args...)
	if e != nil {
		return nil, e
	}
	return openDB(dbPath, network, true)
}

func init() {
	T.Ln("registering ffldb database driver")
	// Register the driver.
	driver := database.Driver{
		DbType: dbType,
		Create: createDBDriver,
		Open:   openDBDriver,
	}
	if e := database.RegisterDriver(driver); E.Chk(e) {
		panic(
			fmt.Sprintf(
				"Failed to regiser database driver '%s': %v",
				dbType, e,
			),
		)
	}
}
