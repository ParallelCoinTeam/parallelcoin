package path

import (
	"path/filepath"

	"github.com/p9c/pod/pkg/conte"
)

// BlockDb returns the path to the block database given a database type.
func BlockDb(cx *conte.Xt, dbType string, namePrefix string) string {
	// The database name is based on the database type.
	dbName := namePrefix + "_" + dbType
	if dbType == "sqlite" {
		dbName += ".db"
	}
	dbPath := filepath.Join(filepath.Join(*cx.Config.DataDir,
		cx.ActiveNet.Name), dbName)
	return dbPath
}
