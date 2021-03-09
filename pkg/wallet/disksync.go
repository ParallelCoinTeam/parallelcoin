package wallet

import (
	"fmt"
	"os"
)

// checkCreateDir checks that the path exists and is a directory. If path does not exist, it is created.
func checkCreateDir(path string) (e error) {
	var fi os.FileInfo
	if fi, e = os.Stat(path); err.Chk(e) {
		if os.IsNotExist(e) {
			// Attempt data directory creation
			if e = os.MkdirAll(path, 0700); err.Chk(e) {
				return fmt.Errorf("cannot create directory: %s", err)
			}
		} else {
			return fmt.Errorf("error checking directory: %s", err)
		}
	} else {
		if !fi.IsDir() {
			return fmt.Errorf("path '%s' is not a directory", path)
		}
	}
	return nil
}
