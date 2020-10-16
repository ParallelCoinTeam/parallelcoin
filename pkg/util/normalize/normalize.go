package normalize

import (
	"net"

	"github.com/urfave/cli"
)

// Address returns addr with the passed default port appended if there is not already a port specified.
func Address(addr, defaultPort string) string {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		Error(err)
		return net.JoinHostPort(addr, defaultPort)
	}
	return addr
}

// Addresses returns a new slice with all the passed peer addresses normalized with the given default port, and all
// duplicates removed.
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

// StringSliceAddresses normalizes a slice of addresses
func StringSliceAddresses(a *cli.StringSlice, port string) {
	variable := []string(*a)
	*a = Addresses(variable, port)
}
