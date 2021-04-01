![Logo](https://raw.githubusercontent.com/p9c/pod/l0k1/pkg/gui/logo/logo_small.svg) 
## ParallelCoin Pod

Fully integrated all-in-one cli client, full node, wallet server, miner and GUI wallet for ParallelCoin

###### [Official Website](https://parallelcoin.info) 
wiki: [![](docs/wiki.svg)](https://github.com/p9c/pod/wiki) social: [![](docs/telegram.svg)](https://t.me/ParallelCoinPlan9) , [![](docs/discord.svg)](https://discord.gg/yB9sYmm3cZ) , [![](docs/fb.svg)](https://www.facebook.com/parallelcoin) & [![](docs/twitter.svg)](https://twitter.com/parallelcoinduo) 

[![github](https://img.shields.io/badge/github-page-blue.svg)](https://p9c.github.io/pod)
[![GoDoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/p9c/pod)
[![Go Report Card](https://goreportcard.com/badge/github.com/p9c/pod)](https://goreportcard.com/report/github.com/p9c/pod)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=alert_status)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=bugs)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=ncloc)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=p9c_pod&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=p9c_pod)
[![CodeScene general](https://codescene.io/images/analyzed-by-codescene-badge.svg)](https://codescene.io/projects/7291)

##### *If you are looking for the legacy 1.2.0 version, go here [https://github.com/p9c/pod/tree/master/legacy](https://github.com/p9c/pod/tree/master/legacy)* - see [https://github.com/p9c/pod/blob/master/legacy/src/linux-install.sh](https://github.com/p9c/pod/blob/master/legacy/src/linux-install.sh) for how to build it from the source on Linux

##### For the legacy version in general, binary releases are [here](https://github.com/p9c/pod/releases/tag/1.2.0)

# WARNING: work in progress, will probably break

## Installation

Straight to business, this is the part I am looking for, so it's here at the top.

### Ubuntu

These instructions expect that you understand basic usage of the bash shell and terminal,
and git usage. You should be able to edit a text file on the command line interface
though you could use `gedit` or similar if you prefer.

First, you need a working [Go 1.14+ installation for the platform you are 
using](https://golang.org).

```bash
cd ~
sudo apt get install -y wget git build-essential
wget https://golang.org/dl/go1.14.13.linux-amd64.tar.gz
tar zxvf go1.14.13.linux-amd64.tar.gz
```

The build tools expect a working GOBIN which is also in PATH, so, open `~/.bashrc` in 
a text editor and put in it:

```bash
export GOBIN=$HOME/bin
export GOPATH=$HOME
export GOROOT=$HOME/go
export PATH=$GOBIN:$PATH
```

then run

```bash
source ~/.bashrc
```

Clone this repository where you like:

```
cd /where/you/keep/your/things
git clone https://github.com/p9c/pod.git
cd pod
```

Before you can build it, though, see [gioui.org install instructions](https://gioui.org/doc/install) - note that on 
linux xsel or other clipboard apps are required for clipboard functionality. The instructions there cover Ubuntu.

Several important libraries are required to build on each platform.
Linux needs some input related X libraries, wayland and their GL
libraries, and similar but different for Mac, Windows, iOS and Android.

More detailed instructions will follow as we work through each 
platform build. For now we develop on FreeBSD and Ubuntu so for now,
at this early stage with the GUI, please bear with us.

Next, go to the repo root and build it.

```
make -B stroy
```

`stroy` is a pure go replacement for `make`. Run it without any options to get
a list of currently enabled build options. 

If you don't have access to make (it can be tricky on windows), you can 
build `stroy` with this command:

```bash
go install -v github.com/p9c/pod/stroy
stroy stroy
```

`stroy` replaces make as it is written in Go, and is customised to put values
into the binaries that show what the binary was built from.

The second invocation rebuilds stroy with stroy itself, which plants the 
build information into the main package, which records the details of the source
code that it was built with, eg:

```
app information: repo: github.com/p9c/pod branch: refs/heads/l0k1 commit: 8b98e47b8a61a7e68a945ce65129f4dd49c6c086 built: 2021-01-12T06:34:34+01:00 tag: v0.4.25-testnet+...
```

The purpose of this is to improve debugging as this info is printed at the 
start of all components' startup, and allows the logs to inform developers 
what the user's build came from, when dealing with a bug report, to 
facilitate *exact* replication of the fault.

## Running

As you will find when you run `stroy` without a build command name, there
is several commands that will launch wallet, node, the GUI, reset the 
wallet, and so on.

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

wiki: [![](docs/wiki.svg)](https://github.com/p9c/pod/wiki) social: [![](docs/telegram.svg)](https://t.me/ParallelCoinPlan9) , [![](docs/discord.svg)](https://discord.gg/yB9sYmm3cZ) , [![](docs/fb.svg)](https://www.facebook.com/parallelcoin) & [![](docs/twitter.svg)](https://twitter.com/parallelcoinduo) 

### Contribution Guidelines

Loki, the primary author of this code has somewhat unconventional ideas about everything including Go, developed out of the process of building this application and the use of compact logging syntax with visual recognisability as a key debugging technique for mostly multi-threaded code, and so, tools have been developed and conventions are found throughout the code and increasingly as all the dusty parts get looked at again.

So, there's a few things that you should do, which Loki thinks you will anyway realise that should be the default anyway. He will stop speaking in the third person now.

1. e is error

being a vowel, `err` is a stolen name. But it also is 3 times as long. It is nearly universal among programmers that i, usually also j and k are iterator variables, and commonly also x, y, z, most cogently relevant when they are coordinates. So use e, not err

2. Use the logger, E.Chk for serious conditions or while debugging, non-blocking failures use D.Chk 

It's not difficult, literally just copy any log.go file out of this repo, change the package name and then you can use the logger, which includes these handy check functions that some certain Go Authors use in several of their lectures and texts, and actually probably the origin of the whole idea of making the error type a pointer to string, which led to the Wrap/Unwrap interface which I don't use. 

4. Use if statements with ALL calls that return errors or booleans

And declare their variables with var statements, unless there is no more uses of the `e` again before all paths of execution return.

5. Prefer to return early

The `.Chk` functions return true if not nil, and so without the `!` in front, the most common case, the content of the following `{}` are the error handling code. 

In general, and especially when the process is not idempotent (changed order breaks the process), which will be most of the time with several processes in a sequence, especially in Gio code which is naturally a bit wide if you write it to make it easily sliced, you want to keep the success line of execution as far left as possible.

Sometimes the negative condition is ignored, as there is a retry or it is not critical, and for these cases use `!E.Chk` and put the success path inside its if block.

6. In the Gio code of the wallet, take advantage of the fact that you can break a line after the dot `.` operator and as you will see amply throughout the code, as it allows items to be repositioned and added with minimal fuss.

Deep levels of closures and stacked variables of any kind tend to lead quickly to a lot of nasty red lines in one's IDE. The fluent method chaining pattern is used because it is far more concise than the raw type definition followed by closure attached by a dot operator, but since it would be valid and pass `gofmt` unchanged to put the whole chain in there (assuming it somehow had no anonymous functions in it).

7. Use the logger

This is being partly repeated as it's very important. Regardless of programmer opinions about whether a debugger is a better tool than a log viewer, note that while it is not fully implemented, `pod` already contains a protocol to aggregate child process logs invisibly from the terminal through an IPC, and logging is one means to enabling auditability of code. 

So long as logs concern themselves primarily with metadata information and only expose data in `trace` level (with the `T` error type) and put the really heavy stuff like printing tree walks over thousands of nodes or other similarly very short operations, put them inside closures with `T.C(func()string{})`

