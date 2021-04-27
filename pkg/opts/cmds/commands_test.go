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
		{Name: "gui", Title:
		"ParallelCoin GUI Wallet/Miner/Explorer",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "version", Title:
		"print version and exit",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "ctl", Title:
		"command line wallet and chain RPC client",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "node", Title:
		"ParallelCoin blockchain node",
			Entrypoint: func(c interface{}) error { return nil },
			Commands: []Command{
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
			Entrypoint: func(c interface{}) error { return nil },
			Commands: []Command{
				{Name: "drophistory", Title:
				"reset the wallet transaction history",
					Entrypoint: func(c interface{}) error { return nil },
				},
			},
		},
		{Name: "kopach", Title:
		"standalone multicast miner for easy mining farm deployment",
			Entrypoint: func(c interface{}) error { return nil },
		},
		{Name: "worker", Title:
		"single thread worker process, normally started by kopach",
			Entrypoint: func(c interface{}) error { return nil },
		},
	}
	return
}

