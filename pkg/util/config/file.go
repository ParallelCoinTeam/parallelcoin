package cfgutil

import (
	"github.com/p9c/pod/pkg/log"
	"os"
)

// FileExists reports whether the named file or directory exists.
func FileExists(	filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		log.ERROR(err)
if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
