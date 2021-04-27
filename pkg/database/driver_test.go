package database_test

import (
	"fmt"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/database"
	_ "github.com/p9c/matrjoska/pkg/database/ffldb"
)

var (
	// ignoreDbTypes are types which should be ignored when running tests that iterate all supported DB types. This allows
	// some tests to add bogus drivers for testing purposes while still allowing other tests to easily iterate all
	// supported drivers.
	ignoreDbTypes = map[string]bool{"createopenfail": true}
)

// checkDbError ensures the passed error is a database.DBError with an error code that matches the passed  error code.
func checkDbError(t *testing.T, testName string, gotErr error, wantErrCode database.ErrorCode) bool {
	dbErr, ok := gotErr.(database.DBError)
	if !ok {
		t.Errorf("%s: unexpected error type - got %T, want %T",
			testName, gotErr, database.DBError{},
		)
		return false
	}
	if dbErr.ErrorCode != wantErrCode {
		t.Errorf("%s: unexpected error code - got %s (%s), want %s",
			testName, dbErr.ErrorCode, dbErr.Description,
			wantErrCode,
		)
		return false
	}
	return true
}

// TestAddDuplicateDriver ensures that adding a duplicate driver does not overwrite an existing one.
func TestAddDuplicateDriver(t *testing.T) {
	supportedDrivers := database.SupportedDrivers()
	if len(supportedDrivers) == 0 {
		t.Errorf("no backends to test")
		return
	}
	dbType := supportedDrivers[0]
	// bogusCreateDB is a function which acts as a bogus create and open driver function and intentionally returns a
	// failure that can be detected if the interface allows a duplicate driver to overwrite an existing one.
	bogusCreateDB := func(args ...interface{}) (database.DB, error) {
		return nil, fmt.Errorf("duplicate driver allowed for database "+
			"type [%v]", dbType,
		)
	}
	// Create a driver that tries to replace an existing one. Set its create and open functions to a function that
	// causes a test failure if they are invoked.
	driver := database.Driver{
		DbType: dbType,
		Create: bogusCreateDB,
		Open:   bogusCreateDB,
	}
	testName := "duplicate driver registration"
	e := database.RegisterDriver(driver)
	if !checkDbError(t, testName, e, database.ErrDbTypeRegistered) {
		return
	}
}

// TestCreateOpenFail ensures that errors which occur while opening or closing a database are handled properly.
func TestCreateOpenFail(t *testing.T) {
	// bogusCreateDB is a function which acts as a bogus create and open driver function that intentionally returns a
	// failure which can be detected.
	dbType := "createopenfail"
	openError := fmt.Errorf("failed to create or open database for "+
		"database type [%v]", dbType,
	)
	bogusCreateDB := func(args ...interface{}) (database.DB, error) {
		return nil, openError
	}
	// Create and add driver that intentionally fails when created or opened to ensure errors on database open and
	// create are handled properly.
	driver := database.Driver{
		DbType: dbType,
		Create: bogusCreateDB,
		Open:   bogusCreateDB,
	}
	e := database.RegisterDriver(driver)
	if e != nil {
		t.Log(e)
	}
	// Ensure creating a database with the new type fails with the expected error.
	_, e = database.Create(dbType)
	if e != openError {
		t.Errorf("expected error not received - got: %v, want %v", e,
			openError,
		)
		return
	}
	// Ensure opening a database with the new type fails with the expected error.
	_, e = database.Open(dbType)
	if e != openError {
		t.Errorf("expected error not received - got: %v, want %v", e,
			openError,
		)
		return
	}
}

// TestCreateOpenUnsupported ensures that attempting to create or open an unsupported database type is handled properly.
func TestCreateOpenUnsupported(t *testing.T) {
	// Ensure creating a database with an unsupported type fails with the expected error.
	testName := "create with unsupported database type"
	dbType := "unsupported"
	_, e := database.Create(dbType)
	if !checkDbError(t, testName, e, database.ErrDbUnknownType) {
		return
	}
	// Ensure opening a database with the an unsupported type fails with the expected error.
	testName = "open with unsupported database type"
	_, e = database.Open(dbType)
	if !checkDbError(t, testName, e, database.ErrDbUnknownType) {
		return
	}
}
