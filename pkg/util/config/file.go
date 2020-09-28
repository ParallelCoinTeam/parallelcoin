package cfgutil

import (
	"github.com/p9c/pkg/app/slog"
	"os"
)

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) (bool, err error) {
	_, err := os.Stat(filePath)
	if err != nil {
		slog.Error(err)
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
