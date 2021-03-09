package ctl

import "github.com/p9c/pod/pkg/rpc/btcjson"

// unusableFlags are the command usage flags which this utility are not able to use. In particular it doesn't support
// websockets and consequently notifications.
const unusableFlags = btcjson.UFWebsocketOnly | btcjson.UFNotification

// ListCommands categorizes and lists all of the usable commands along with their one-line usage.
func ListCommands() (s string) {
	const (
		categoryChain uint8 = iota
		categoryWallet
		numCategories
	)
	// Get a list of registered commands and categorize and filter them.
	cmdMethods := btcjson.RegisteredCmdMethods()
	categorized := make([][]string, numCategories)
	for _, method := range cmdMethods {
		var e error
		var flags btcjson.UsageFlag
		if flags, e = btcjson.MethodUsageFlags(method); err.Chk(e) {
			continue
		}
		// Skip the commands that aren't usable from this utility.
		if flags&unusableFlags != 0 {
			continue
		}
		var usage string
		if usage, e = btcjson.MethodUsageText(method); err.Chk(e) {
			continue
		}
		// Categorize the command based on the usage flags.
		category := categoryChain
		if flags&btcjson.UFWalletOnly != 0 {
			category = categoryWallet
		}
		categorized[category] = append(categorized[category], usage)
	}
	// Display the command according to their categories.
	categoryTitles := make([]string, numCategories)
	categoryTitles[categoryChain] = "Chain Server Commands:"
	categoryTitles[categoryWallet] = "Wallet Server Commands (--wallet):"
	for category := uint8(0); category < numCategories; category++ {
		s += categoryTitles[category]
		s += "\n"
		for _, usage := range categorized[category] {
			s += "\t" + usage + "\n"
		}
		s += "\n"
	}
	return
}
