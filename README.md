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

#### 26 November 2019 update

The full set of features aside from the GUI have now been implemented and last details before the beta are in process and the GUI will be finished within a week or two. Watch this space.

## Building

You can just `go install` in the repository root and `pod` will be placed in your `GOBIN` directory.

## Installation

TODO: Initial release will include Linux, Mac and Windows binaries including the GUI, 
binaries for all platform targets of Go 1.12.9+ without the GUI and standalone kopach
miner also for all targets of Go v1.12.9+.

## Developer Notes

Goland's inbuilt terminal is very slow and has several bugs that my workflow
exposes and I find very irritating, and out of the paned terminal apps I find
Tilix the most usable, but it requires writing a regular expression to
match and so the logger is written to output relative paths to the
repository root.

The regexp that I use given my system base path is (exactly this with all newlines removed for dconf with using tilix at the dconf path `/com/gexperts/Tilix/custom-hyperlinks`)

```
[
    ' [/]((([a-zA-Z0-9-_.]+/)+([a-zA-Z0-9-_.]+)):([0-9]+)),
        <goland executable> --line $5 /$1,false', 
    'github[.]((([a-zA-Z0-9-_.]+/)+([a-zA-Z0-9-_.]+)):([0-9]+)),
        <goland executable> --line $5 <$GOPATH>/src/github.$1,
        false', 
    '((([a-zA-Z0-9-_.]+/)+([a-zA-Z0-9-_.]+)):([0-9]+)),
        /usr/local/bin/goland --line $5 <$GOPATH>/src/github.com/p9c/pod/$1,
        false'
]
```

(the text fields in tilix's editor are very weird so it will be easier to
just paste this in and gnome dconf editor should remove the newlines
automatically)

Replace the parts inside `<` `>` with the relevant path from your environment
and enjoy quickly hopping to source code locations while you are working on
this project. Goland's terminal recognises most of them anyway but when you
get a runaway log print going on it can take some time before the terminal
decides it will listen to your control C.
  
The configuration shown above covers the most common relative to project root
paths as used in the logger, as well as `go get` style paths starting with
the domain name, as well as absolute paths (first one) that will work for
any relevant file path with line number reference.
