package sanitizers

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	NetAddress = "netaddress"
	Password   = "password"
	FilePath   = "filepath"
	Directory  = "directory"
)

func StringType(typ, input string, defaultPort int) (cleaned string, e error) {
	switch typ {
	case NetAddress:
		var h, p string
		if h, p, e = net.SplitHostPort(input); E.Chk(e) {
			e = fmt.Errorf("address value '%s' not a valid address", input)
			return
		}
		if p == "" {
			cleaned = net.JoinHostPort(h, fmt.Sprint(defaultPort))
		}
	case Password:
		// password type is mainly here for the input method of the app using this config library
	case FilePath:
		if strings.HasPrefix(input, "~") {
			var homeDir string
			var usr *user.User
			var e error
			if usr, e = user.Current(); e == nil {
				homeDir = usr.HomeDir
			}
			// Fall back to standard HOME environment variable that works for most POSIX OSes if the directory from the Go
			// standard lib failed.
			if e != nil || homeDir == "" {
				homeDir = os.Getenv("HOME")
			}
			
			input = strings.Replace(input, "~", homeDir, 1)
		}
		if cleaned, e = filepath.Abs(filepath.Clean(input)); E.Chk(e) {
		}
	case Directory:
		if strings.HasPrefix(input, "~") {
			var homeDir string
			var usr *user.User
			var e error
			if usr, e = user.Current(); e == nil {
				homeDir = usr.HomeDir
			}
			// Fall back to standard HOME environment variable that works for most POSIX OSes if the directory from the Go
			// standard lib failed.
			if e != nil || homeDir == "" {
				homeDir = os.Getenv("HOME")
			}
			
			input = strings.Replace(input, "~", homeDir, 1)
		}
		if cleaned, e = filepath.Abs(filepath.Clean(input)); E.Chk(e) {
		}
	default:
		cleaned = input
	}
	return
}
