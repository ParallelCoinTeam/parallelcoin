package bdb

import (
	"fmt"
	
	"github.com/p9c/matrjoska/pkg/walletdb"
)

const (
	dbType = "bdb"
)

// parseArgs parses the arguments from the walletdb Open/Create methods.
func parseArgs(funcName string, args ...interface{}) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf(
			"invalid arguments to %s.%s -- "+
				"expected database path", dbType, funcName,
		)
	}
	dbPath, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf(
			"first argument to %s.%s is invalid -- "+
				"expected database path string", dbType, funcName,
		)
	}
	return dbPath, nil
}

// openDBDriver is the callback provided during driver registration that opens
// an existing database for use.
func openDBDriver(args ...interface{}) (d walletdb.DB, e error) {
	var dbPath string
	if dbPath, e = parseArgs("Open", args...); E.Chk(e) {
		return
	}
	return openDB(dbPath, false)
}

// createDBDriver is the callback provided during driver registration that
// creates, initializes, and opens a database for use.
func createDBDriver(args ...interface{}) (d walletdb.DB, e error) {
	var dbPath string
	if dbPath, e = parseArgs("Create", args...); E.Chk(e) {
		return
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
	var e error
	if e = walletdb.RegisterDriver(driver); E.Chk(e) {
		panic(
			fmt.Sprintf(
				"Failed to regiser database driver '%s': %v",
				dbType, e,
			),
		)
	}
}
