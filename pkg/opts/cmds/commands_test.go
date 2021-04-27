package cmds

import (
	"testing"
)

func TestCommands_GetAllCommands(t *testing.T) {
	cm := GetCommands()
	I.S(cm.GetAllCommands())
}


// GetCommands returns available subcommands in Parallelcoin Pod
func GetCommands() (c Commands) {
	c = Commands{
		{Name: "gui", Description:
		"ParallelCoin GUI Wallet/Miner/Explorer",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "version", Description:
		"print version and exit",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "ctl", Description:
		"command line wallet and chain RPC client",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "node", Description:
		"ParallelCoin blockchain node",
			Entrypoint: func(c interface{}) error { return nil },
			Commands: []Command{
				{Name: "dropaddrindex", Description:
				"drop the address database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "droptxindex", Description:
				"drop the transaction database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "dropcfindex", Description:
				"drop the cfilter database index",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "dropindexes", Description:
				"drop all of the indexes",
					Entrypoint: func(c interface{}) error { return nil },
				},
				{Name: "resetchain", Description:
				"deletes the current blockchain cache to force redownload",
					Entrypoint: func(c interface{}) error { return nil },
				},
			},
		},
		{Name: "wallet", Description:
		"run the wallet server (requires a chain node to function)",
			Entrypoint: func(c interface{}) error { return nil },
			Commands: []Command{
				{Name: "drophistory", Description:
				"reset the wallet transaction history",
					Entrypoint: func(c interface{}) error { return nil },
				},
			},
		},
		{Name: "kopach", Description:
		"standalone multicast miner for easy mining farm deployment",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "worker", Description:
		"single thread worker process, normally started by kopach",
			Entrypoint: func(c interface{}) error { return nil },
		},
	}
	return
}

