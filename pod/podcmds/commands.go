package podcmds

import (
	"fmt"

	"github.com/gookit/color"

	"github.com/p9c/matrjoska/pod/launchers"
	"github.com/p9c/matrjoska/version"
	"github.com/p9c/opts/cmds"
)

// GetCommands returns available subcommands in Parallelcoin Pod
func GetCommands() (c cmds.Commands) {
	c = cmds.Commands{
		{Name: "gui", Title:
		"ParallelCoin GUI Wallet/Miner/Explorer",
			Entrypoint: launchers.GUIHandle,
			Colorizer:  color.Bit24(128, 255, 255, false).Sprint,
			AppText:    "   gui",
		},
		{Name: "version", Title:
		"print version and exit",
			Entrypoint: func(c interface{}) error {
				fmt.Println(version.Tag)
				return nil
			},
		},
		{Name: "ctl", Title:
		"command line wallet and chain RPC client",
			Entrypoint: launchers.CtlHandle,
			Colorizer:  color.Bit24(128, 255, 128, false).Sprint,
			AppText:    "   ctl",
			Commands: []cmds.Command{
				{Name: "list", Title:
				"list available commands",
					Entrypoint: launchers.CtlHandleList,
				},
			},
		},
		{Name: "node", Title:
		"ParallelCoin Blockchain Node",
			Entrypoint: launchers.NodeHandle,
			Colorizer:  color.Bit24(128, 128, 255, false).Sprint,
			AppText:    "  node",
			Description: "The ParallelCoin Node synchronises with the chain" +
			" and can be used to search for block and transaction data as" +
			" well as sending out new transactions.\n" +
			"It can be used as the chain server for a wallet or other" +
			" application server consuming chain data",
			Commands: []cmds.Command{
				{Name: "dropaddrindex", Title:
				"drop the address database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "droptxindex", Title:
				"drop the transaction database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "dropcfindex", Title:
				"drop the cfilter database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "dropindexes", Title:
				"drop all of the indexes",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "resetchain", Title:
				"deletes the current blockchain cache to force redownload",
					Entrypoint: func(c interface{}) error { return nil },
				},
			},
		},
		{Name: "wallet", Title:
		"run the wallet server (requires a chain node to function)",
			Entrypoint: launchers.WalletHandle, // func(c interface{}) error { return nil },
			Commands: []cmds.Command{
				{Name: "drophistory", Title:
				"reset the wallet transaction history",
					Entrypoint: func(c interface{}) error { return nil },
				},
			},
			Colorizer: color.Bit24(255, 255, 128, false).Sprint,
			AppText:   "wallet",
		},
		{Name: "kopach", Title:
		"standalone multicast miner for easy mining farm deployment",
			Entrypoint: launchers.Kopach,
			Colorizer:  color.Bit24(255, 128, 128, false).Sprint,
			AppText:    "kopach",
		},
		{Name: "worker", Title:
		"single thread worker process, normally started by kopach",
			Entrypoint: launchers.Worker,
			Colorizer:  color.Bit24(255, 128, 255, false).Sprint,
			AppText:    "worker",
		},
	}
	c.PopulateParents(nil)
	// I.S(c)
	return
}

