# chaincfg

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/p9c/pod/chaincfg)

Package chaincfg defines chain configuration parameters for the three standard Parallelcoin networks and provides the ability for callers to define their own custom networks.

Although this package was primarily written for pod, it has intentionally been designed so it can be used as a standalone package for any projects needing to use parameters for the standard Bitcoin networks or for projects needing to define their own network.

## Sample Use

```Go
package main

import (
	"flag"
	"fmt"
	"log"
	
	"github.com/p9c/pod/pkg/chain/config/netparams"
	
	"git.parallelcoin.io/util"
	"github.com/p9c/pod/pkg/chain/config"
)

var testnet = flag.Bool("testnet", false, "operate on the testnet Bitcoin network")

// By default (without -testnet), use mainnet.
var chainParams = &netparams.MainNetParams

func main() {
	flag.Parse()
	// Modify active network parameters if operating on testnet.
	if *testnet {
		chainParams = &netparams.TestNet3Params
	}
	// later...
	// Create and print new payment address, specific to the active network.
	pubKeyHash := make([]byte, 20)
	addr, e := util.NewAddressPubKeyHash(pubKeyHash, chainParams)
	if e != nil  {
		L.	}
	log.Println(addr)
}
```

## Installation and Updating

```bash
$ go get -u github.com/p9c/pod/chaincfg
```

## License

Package chaincfg is licensed under the [copyfree](http://copyfree.org) ISC
License.
