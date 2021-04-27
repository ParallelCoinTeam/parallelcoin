package bdb_test

// This file intended to be copied into each backend driver directory. Each driver should have their own driver_test.go
// file which creates a database and invokes the testInterface function in this file to ensure the driver properly
// implements the interface. See the bdb backend driver for a working example.
//
// NOTE: When copying this file into the backend driver folder, the package name will need to be changed accordingly.
import (
	"os"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/walletdb/ci"
)

// TestInterface performs all interfaces tests for this database driver.
func TestInterface(t *testing.T) {
	
	dbPath := "interfacetest.db"
	defer func() {
		if e := os.RemoveAll(dbPath); E.Chk(e) {
		}
	}()
	ci.TestInterface(t, dbType, dbPath)
}
