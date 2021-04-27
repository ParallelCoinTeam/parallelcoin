package cfgutil

import (
	"os"
)

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) (bool, error) {
	_, e := os.Stat(filePath)
	if e != nil {
		E.Ln(e)
		if os.IsNotExist(e) {
			return false, nil
		}
		return false, e
	}
	return true, nil
}
