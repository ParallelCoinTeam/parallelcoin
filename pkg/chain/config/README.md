# chaincfg

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/stalker-loki/pod/chaincfg)

Package chaincfg defines chain configuration parameters for the three standard Parallelcoin networks and provides the ability for callers to define their own custom networks.

Although this package was primarily written for pod, it has intentionally been designed so it can be used as a standalone package for any projects needing to use parameters for the standard Bitcoin networks or for projects needing to define their own network.

## Sample Use

```Go
package main
import (
"flag"
"fmt"
"github.com/stalker-loki/pod/pkg/chain/config/netparams"
"log"




"github.com/stalker-loki/pod/pkg/chain/config"
"git.parallelcoin.io/util"
)
var testnet = flag.Bool("testnet", false, "operate on the testnet Bitcoin network")
// By default (without -testnet), use mainnet.
var chainParams = &netparams.MainNetParams
func main(	) {
	flag.Parse()
	// Modify active network parameters if operating on testnet.
	if *testnet {
        chainParams = &netparams.TestNet3Params 
	}
	// later...
	// Create and print new payment address, specific to the active network.
	pubKeyHash := make([]byte, 20)
	addr, err := util.NewAddressPubKeyHash(pubKeyHash, chainParams)
	if err != nil {
		L.Error(err)
	}
	log.Println(addr)
}
```

## Installation and Updating

```bash
$ go get -u github.com/stalker-loki/pod/chaincfg
```

## License

Package chaincfg is licensed under the [copyfree](http://copyfree.org) ISC
License.
