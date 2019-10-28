# ![Logo](https://git.parallelcoin.io/dev/legacy/raw/commit/f709194e16960103834b0d0e25aec06c3d84f85b/logo/logo48x48.png) Parallelcoin Pod 

[![GoDoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/p9c/pod) 
[![master branch](https://img.shields.io/badge/branch-master-gray.svg)](https://github.com/p9c/pod) 
[![discord chat](https://img.shields.io/badge/discord-chat-purple.svg)](https://discord.gg/YgBWNgK)

Fully integrated all-in-one cli client, full node, wallet server, miner and GUI wallet for Parallelcoin

#### Binaries for legacy now available for linux amd64

Get them from here: [https://git.parallelcoin.io/dev/parallelcoin-binaries](https://git.parallelcoin.io/dev/parallelcoin-binaries)

Pod is a multi-application with multiple submodules for different functions. 
It is self-configuring and configurations can be changed from the commandline
 as well as editing the json files directly, so the binary itself is the
  complete distribution for the suite.

It consists of 6 main modules:

1. pod/ctl - command line interface to send queries to a node or wallet and 
    prints the results to the stdout
2. pod/node - full node for Parallelcoin network, including optional indexes for 
    address and transaction search, low latency miner UDP broadcast based controller
3. pod/wallet - wallet server that runs separately from the full node but 
    depends on a full node RPC for much of its functionality. Currently does not
    have a full accounts implementation (TODO: fixme!)
4. pod/shell - combined full node and wallet server of 2. and 3. running 
    concurrently
5. pod/gui - webview based desktop wallet GUI
6. pod/kopach - standalone miner with LAN UDP broadcast work delivery system

## Building

You can just `go install` in the root directory and `pod` will be placed in your `GOBIN` directory.

## Installation

TODO: Initial release will include Linux, Mac and Windows binaries including the GUI, 
binaries for all platform targets of Go 1.12.9+ without the GUI and standalone kopach
miner also for all targets of Go v1.12.9+.

## Tilix custom hyperlinks

The documentation of Tilix is not the best neither is the hyperlinks
 interface, but after much frustration I was able to find both the regexp and
  the right command to use with it to allow Tilix (my preferred terminal
   because of its sweet gtk+-3 interface and fast VTE backend) to open
    relative path links that are printed in the logs:
    
    #### regexp:
    
```
(([a-zA-Z0-9-_.]+/)+([a-zA-Z0-9-_.]+)):([0-9]+)
```

#### goland launch command:

```
goland --line $4 $GOPATH/src/github.com/p9c/pod/$1
```

Change the $GOPATH as required for the absolute path of your copy of this repo.
