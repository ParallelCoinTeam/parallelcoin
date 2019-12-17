# ![Logo](https://git.parallelcoin.io/dev/legacy/raw/commit/f709194e16960103834b0d0e25aec06c3d84f85b/logo/logo48x48.png) Parallelcoin Pod 

[![GoDoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/p9c/pod) 

Fully integrated all-in-one cli client, full node, wallet server, miner and GUI wallet for Parallelcoin

## Installation

Straight to business, this is the part I am looking for so it's here at the top.

First, you need a working [Go 1.11+ installation for the platform you are using](https://golang.org).

Clone this repository where you like, here I show the SSH URL which is recommended
for speed as well as if you want to add a branch to the repository as a member of the 
team (github account and a registered SSH public key on it is required):

```
cd /where/you/keep/your/things
git clone git@github.com:p9c/pod.git
cd pod
go install -v
```

You should use modules for this project, as everyone else is and many
forgot to protect their master from version 2 on the same URL.

## Running

If you just want to use it as an RPC for only node services at localhost:11047 (no wallet)

```
pod node
```

For wallet only at localhost:11048 (a full node must be configured, by default should be found at localhost:11047)

```
pod wallet
```

For combined RPC wallet at localhost:11046

```
pod shell
```

For the standalone multicast miner worker 'Kopach':

```
pod kopach
```

The list of commands and options can be seen using the following command:

```
pod help
```

## Notable items and their short forms:

### `-D`

Set the root folder of the data directory. Default is ~/.pod or the string 'pod'
as the folder name in other systems.

### `-g`

`-g=false` disables mining

Enable mining, using inbuilt for run modes that enable a p2p blockchain node

### `-G` 

Set the number of threads to mine with. Performance with the Plan 9 hardfork
will entirely depend on the performance characteristics of the processor and 
its' long division units and how they are scheduled. The inbuilt miner
(which will be deprecated) has significantly inferior performance. Concurrency is
not parallelism, and the stand-alone miner is better. The inbuilt miner will
be entirely removed by release.

### `-n`

Set the network type, mainnet and testnet are the main important options. Note
that this is the main configuration as well as pre-shared key, to run the multi-
cast mining system, as the different networks have different start heights for
hard forks.

## Configuration

Configuration is designed to be largely automatic, however manual edits can be
made, from `<pod profile directory>`/pod.json - notably critical elements for
the cluster mining configurations is the 'MiningPass' item matches up between 
nodes you intend to communicate with each other.

### Mining Farm Setup

For the time being all that is necessary is to copy the `pod.json` file, and 
that all nodes deployed are on the same subnets as the nodes. Note that it is
possible to isolate subnets and join them using nodes via dual network (virtual)
interfaces and that worker nodes trust implicitly all nodes that use the same
pre shared key (thus the configuration file).

Before beta release there will be a FreeBSD based live image that is written
to using a utility app with the correct key and network settings and will be
basically turn-key if used as default configured. BSD is being used because it 
is lighter and ensures your hardware is doing nothing more than exactly crunching
giant numbers for the chance to get a block reward.s
 
### Configuration for adjunct services (block explorers, exchanges)

`rpc.cert` `ca.cert` and `rpc.key` files, which as they are can be used (not so
securely) for connecting nodes in one's server set up. The system can be run by
default in an 'insecure' configuration (they are wired to connect via localhost
ports). Presumably for this kind of production application one would use a complete
set of ports and custom CA file. What is provided by default is for development
purposes and on a relatively unconnected end user setup. 

Further improvements in security are planned. 

For now it is advisable to isolate wallet services strongly and the main attack
vector is covered. Easier to use GUI interface for offline transaction signing
and similar features also are planned for later implementation.

### GUI build info

The GUI subsystem can be disabled in the build using

```
go install -tags headless
```

To build it, there are some GL and X prerequisites for the
Linux build

```
sudo apt-get install libgles2-mesa-dev \
     libxkbcommon-dev \
     libxkbcommon-x11-dev
```

More info about building for other platforms to follow. 
There should be a build for Android and iOS eventually, they
have extra build environment requirements (android sdk and 
xcode/mac respectively). Specifics for Windows builds also to come.

## Binaries for legacy (pre hardfork) now available for linux amd64

Get them from here: [https://git.parallelcoin.io/dev/parallelcoin-binaries](https://git.parallelcoin.io/dev/parallelcoin-binaries)

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
        <goland executable> --line $5 <$GOPATH>/src/github.com/p9c/pod/$1,
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
