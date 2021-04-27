package normalize

import (
	"net"
)

// address returns addr with the passed default port appended if there is not
// already a port specified.
func address(addr, defaultPort string) string {
	var e error
	if _, _, e = net.SplitHostPort(addr); E.Chk(e) {
		return net.JoinHostPort(addr, defaultPort)
	}
	return addr
}

// Addresses returns a new slice with all the passed peer addresses normalized
// with the given default port, and all duplicates removed.
func Addresses(addrs []string, defaultPort string) []string {
	for i := range addrs {
		addrs[i] = address(addrs[i], defaultPort)
	}
	return RemoveDuplicateAddresses(addrs)
}

// RemoveDuplicateAddresses returns a new slice with all duplicate entries in
// addrs removed.
func RemoveDuplicateAddresses(addrs []string) (result []string) {
	result = make([]string, 0, len(addrs))
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
func StringSliceAddresses(a []string, port string) {
	variable := a
	a = Addresses(variable, port)
}
