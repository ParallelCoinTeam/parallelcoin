package normalize

import (
	"github.com/p9c/pod/pkg/log"
	"net"

	"github.com/urfave/cli"
)

// Address returns addr with the passed default port appended if there is not already a port specified.
func Address(addr, defaultPort string) string {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.ERROR(err)
		return net.JoinHostPort(addr, defaultPort)
	}
	return addr
}

// Addresses returns a new slice with all the passed peer addresses normalized with the given default port,
// and all duplicates removed.
func Addresses(addrs []string, defaultPort string) []string {
	for i := range addrs {
		addrs[i] = Address(addrs[i], defaultPort)
	}
	return removeDuplicateAddresses(addrs)
}

// RemoveDuplicateAddresses returns a new slice with all duplicate entries in addrs removed.
func removeDuplicateAddresses(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	seen := map[string]struct{}{}
	for _, val := range addrs {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = struct{}{}
		}
	}
	return result
}

// // CleanAndExpandPath expands environment variables and leading ~ in the passed
// // path, cleans the result, and returns it.
// func CleanAndExpandPath(datadir, path string) string {
// 	// Expand initial ~ to OS specific home directory.
// 	if strings.HasPrefix(path, "~") {
// 		homeDir := filepath.Dir(datadir)
// 		path = strings.Replace(path, "~", homeDir, 1)
// 	}
// 	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%, but
// 	// they variables can still be expanded via POSIX-style $VARIABLE.
// 	return filepath.Clean(os.ExpandEnv(path))
// }

// StringSliceAddresses normalizes a slice of addresses
func StringSliceAddresses(a *cli.StringSlice, port string) {
	variable := []string(*a)
	*a = Addresses(variable, port)
}
