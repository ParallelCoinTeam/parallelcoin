package ctl

import (
	"fmt"
	"path/filepath"

	"github.com/stalker-loki/pod/app/appdata"
	"github.com/stalker-loki/pod/pkg/rpc/btcjson"
)

// unusableFlags are the command usage flags which this utility are not able to
// use.  In particular it doesn't support websockets and consequently
// notifications.
const unusableFlags = btcjson.UFWebsocketOnly | btcjson.UFNotification

// nolint
var (
	// DefaultConfigFile is
	DefaultConfigFile = filepath.Join(PodCtlHomeDir, "conf.json")
	// DefaultRPCCertFile is
	DefaultRPCCertFile = filepath.Join(NodeHomeDir, "rpc.cert")
	// DefaultRPCServer is
	DefaultRPCServer = "127.0.0.1:11048"
	// DefaultWallet is
	DefaultWallet = "127.0.0.1:11046"
	// DefaultWalletCertFile is
	DefaultWalletCertFile = filepath.Join(SPVHomeDir, "rpc.cert")
	// NodeHomeDir is
	NodeHomeDir = appdata.Dir("pod", false)
	// PodCtlHomeDir is
	PodCtlHomeDir = appdata.Dir("pod/ctl", false)
	// SPVHomeDir is
	SPVHomeDir = appdata.Dir("pod/spv", false)
)

// ListCommands categorizes and lists all of the usable commands along with
// their one-line usage.
func ListCommands() {
	const (
		categoryChain uint8 = iota
		categoryWallet
		numCategories
	)
	// Get a list of registered commands and categorize and filter them.
	cmdMethods := btcjson.RegisteredCmdMethods()
	categorized := make([][]string, numCategories)
	for _, method := range cmdMethods {
		flags, err := btcjson.MethodUsageFlags(method)
		if err != nil {
			Error(err)
			// This should never happen since the method was just returned
			// from the package, but be safe.
			continue
		}
		// Skip the commands that aren't usable from this utility.
		if flags&unusableFlags != 0 {
			continue
		}
		usage, err := btcjson.MethodUsageText(method)
		if err != nil {
			Error(err)
			// This should never happen since the method was just returned
			// from the package, but be safe.
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
		fmt.Println(categoryTitles[category])
		fmt.Println()
		for _, usage := range categorized[category] {
			fmt.Println("  ", usage)
		}
		fmt.Println()
	}
}
