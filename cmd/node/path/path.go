package path

import (
	"path/filepath"
	
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
)

// BlockDbNamePrefix is the prefix for the block database name.
// The database type is appended to this value to form the full block
// database name.
const BlockDbNamePrefix = "blocks"

// BlockDb returns the path to the block database given a database type.
func BlockDb(cx *conte.Xt, dbType string) string {
	// The database name is based on the database type.
	dbName := BlockDbNamePrefix + "_" + dbType
	if dbType == "sqlite" {
		dbName += ".db"
	}
	dbPath := filepath.Join(filepath.Join(*cx.Config.DataDir,
		cx.ActiveNet.Name), dbName)
	return dbPath
}
