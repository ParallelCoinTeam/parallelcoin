# ![Logo](https://raw.githubusercontent.com/stalker-loki/pod/master/pkg/gui/logo/logo.svg) Parallelcoin Pod 

[![github](https://img.shields.io/badge/github-page-blue.svg)](https://p9c.github.io/pod)
[![GoDoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/p9c/pod) 
[![Go Report Card](https://goreportcard.com/badge/github.com/p9c/pod)](https://goreportcard.com/report/github.com/p9c/pod)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=alert_status)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=bugs)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=ncloc)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![CodeScene general](https://codescene.io/images/analyzed-by-codescene-badge.svg)](https://codescene.io/projects/7291)

Fully integrated all-in-one cli client, full node, wallet server, miner and GUI wallet for Parallelcoin

~~~~
φ plan 9 crypto
    protocols and back end           David Vennik
    gui                              Djordje Marcetin
~~~~
## Installation

Straight to business, this is the part I am looking for, so it's here at the top.

First, you need a working [Go 1.11+ installation for the platform you are using](https://golang.org).

Clone this repository where you like, here I show the SSH URL which is recommended
for speed as well as if you want to add a branch to the repository as a member of the 
team (github account and a registered SSH public key on it is required):

```
cd /where/you/keep/your/things
git clone git@github.com:p9c/pod.git
```

Before you can build it, though, see [gioui.org install instructions](https://gioui.org/doc/install)

Several important libraries are required to build on each platform.
Linux needs some input related X libraries, wayland and their GL
libraries, and similar but different for Mac, Windows, iOS and Android.

More detailed instructions will follow as we work through each 
platform build. For now we develop on FreeBSD and Ubuntu so for now,
at this early stage with the GUI, please bear with us.

Next, go to the repo root and get Go to build it.

```
cd pod
go install -v
```

Any version of Go from 1.11 should build, this is really the current
minimal production version for Go anyway with the chaos that 
modules have unleashed on Git repository branch keeping hygiene.

**GO111MODULE should be set to "on".**

If you want to build a version without any GUI, for servers or if support is
lacking on the given platform:

`go install -v -tags headless`

## Running

### Initial configuration:

For initial configuration, use the `-D` and `-n` flags combined with
the `init` subcommand like so:

`pod -D <data directory> -n <mainnet/testnet> init`

This in one step creates a fresh new configuration file, all of the
TLS certificates and default Certificate Authority to use the
web sockets interface for especially the wallet async functionality,
and prompts you on the CLI to enter a new wallet passphrase, gives
seed you need to restore the wallet later, and fills the configuration
with a set of starting mining addresses based on the wallet seed,
for the defined network type.

~~**TODO:**s yes, we want to move these keys into the directory subfolder
so it can be done without the node running and on demand with a new
subcommand for exactly this purpose. New addresses require a wallet 
but should be kept away from a public RPC or other remote protocol
endpoint. Only nodes need them while mining to use for creating
coinbase payment outpoints.~~

### Run Modes

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

#### Windows GUI build needs files: 

Files are included in Pod's root folder
- d3dcompiler_47.dll
- libEGL.dll
- libGLESv2.dll


## Binaries for legacy (pre hardfork) now available for linux amd64

Get them from here: [https://git.parallelcoin.io/dev/parallelcoin-binaries](https://git.parallelcoin.io/dev/parallelcoin-binaries)

## Developer Notes

Goland's inbuilt terminal is very slow and has several bugs that my workflow
exposes and I find very irritating, and out of the paned terminal apps I find
Tilix the most usable, but it requires writing a regular expression to
match and so the logger is written to output relative paths to the
repository root.

The regexp that I use given my system base path is (exactly this with all 
newlines removed for dconf with using tilix at the dconf path 
`/com/gexperts/Tilix/custom-hyperlinks`)

```
[
    '[ ]((([a-zA-Z@0-9-_.]+/)+([a-zA-Z@0-9-_.]+)):([0-9]+))$,goland --line $5 $HOME/Public/pod/$2,false', 
    '([/](([a-zA-Z@0-9-_.]+/)+([a-zA-Z@0-9-_.]+)):([0-9]+)),goland --line $5 /$2,false'
]
```

These two seem to the work the best including allowing clicking on stack trace 
code location references. Change goland launcher and package root path as required.
The logger code locations start with a space and absolute paths with a forward
slash and you have to set the repository path manually.
