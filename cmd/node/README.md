![](https://gitlab.com/parallelcoin/node/raw/master/assets/logo.png)

# The Parallelcoin Node [![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/git.parallelcoin.io/dev/pod/node)

Next generation full node for Parallelcoin, forked from [btcd](https://github.com/btcsuite/btcd)

## Hard Fork 1: Plan 9 from Crypto Space

9 algorithms can be used when mining:

- Blake14lr (decred)
- ~~Skein (myriadcoin)~~ Cryptonote7v2
- Lyra2REv2 (sia)
- Keccac (maxcoin, smartcash)
- Scrypt (litecoin)
- SHA256D (bitcoin)
- GOST Stribog \*
- Skein
- X11 (dash)

### Stochastic Binomial Filter Difficulty Adjustment

After the upcoming hardfork, Parallelcoin will have the following features in its difficulty adjustment regime:

- Exponential curve with power of 3 to respond gently the natural drift while moving the difficulty fast in below 10% of target and 10x target, to deal with recovering after a large increase in network hashpower

- 293 second blocks (7 seconds less than 5 minutes), 1439 block averaging window (about 4.8 days) that is varied by interpreting byte 0 of the sha256d hash of newest block hash as a signed 8 bit integer to further disturb any inherent rhythm (like dithering).

- Difficulty adjustments are based on a window ending at the previous block of each algorithm, meaning sharp rises from one algorithm do not immediately affect the other algorithms, allowing a smoother recovery from a sudden drop in hashrate, soaking up energetic movements more robustly and resiliently, and reducing vulnerability to time distortion attacks.

- Deterministic noise is added to the difficulty adjustment in a similar way as is done with digital audio and images to improve the effective resolution of the signal by reducing unwanted artifacts caused by the sampling process. Miners are random generators, and a block time is like a tuning filter, so the same principles apply.

- Rewards will be computed according to a much smoother, satoshi-precision exponential decay curve that will produce a flat annual 5% supply expansion. Increasing the precision of the denomination is planned for the next release cycle, at 0.00000001 as the minimum denomination, there may be issues as userbase increases.

- Fair Hardfork - Rewards will slowly rise from the initial hard fork at an inverse exponential rate to bring the block reward from 0.02 up to 2 in 2000 blocks, as the adjustment to network capacity takes time, so rewards will closely match the time interval they relate to until it starts faster from the minimum target stabilises in response to what miners create.

### Wallet

The pod has no RPC wallet functionality, only core chain functions. It is fully compliant with the original [parallelcoind](https://github.com/marcetin/parallelcoin) for these functions. For the wallet server, [mod](https://git.parallelcoin.io/dev/pod/walletmain), which works also with the CLI controller [podctl](https://git.parallelcoin.io/dev/pod/), it will be possible to send commands to both mod (wallet) and pod full node using the command line.

A Webview/Golang based GUI wallet will come a little later, following the release, and will be able to run on all platforms with with browser or supported built-in web application platforms, Blink and Webkit engines.

\* At the time of release there will not be any GPU nor ASIC miners for the GOST Stribog (just stribog 256 bit hash, not combined) and Highwayhash.

## Installation

For the main full node server:

```bash
go get git.parallelcoin.io/dev/pod
```

You probably will also want CLI client (can also speak to other bitcoin protocol RPC endpoints also):

```bash
go get git.parallelcoin.io/dev/pod/cmd/podctl
```

## Requirements

[Go](http://golang.org) 1.11 or newer.

## Installation

#### Windows not available yet

When it is, it will be available here:

https://git.parallelcoin.io/dev/pod/releases

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Install Go according to the [installation instructions](http://golang.org/doc/install)
- Ensure Go was installed properly and is a supported version:

```bash
$ go version
$ go env GOROOT GOPATH
```

NOTE: The `GOROOT` and `GOPATH` above must not be the same path. It is recommended that `GOPATH` is set to a directory in your home directory such as `~/goprojects` to avoid write permission issues. It is also recommended to add `$GOPATH/bin` to your `PATH` at this point.

- Run the following commands to obtain pod, all dependencies, and install it:

```bash
$ go get git.parallelcoin.io/dev/pod
```

- pod (and utilities) will now be installed in `$GOPATH/bin`. If you did
  not already add the bin directory to your system path during Go installation,
  we recommend you do so now.

## Updating

#### Windows

Install a newer MSI

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Run the following commands to update pod, all dependencies, and install it:

```bash
$ cd $GOPATH/src/git.parallelcoin.io/dev/pod
$ git pull && glide install
$ go install . ./cmd/...
```

## Getting Started

pod has several configuration options available to tweak how it runs, but all of the basic operations described in the intro section work with zero configuration.

#### Windows (Installed from MSI)

Launch pod from your Start menu.

#### Linux/BSD/POSIX/Source

```bash
$ ./pod
```

## Discord

Come and chat at our (discord server](https://discord.gg/nJKts94)

## Issue Tracker

The [integrated github issue tracker](https://git.parallelcoin.io/dev/pod/issues)
is used for this project.

## Documentation

The documentation is a work-in-progress. It is located in the [docs](https://git.parallelcoin.io/dev/pod/tree/master/docs) folder.

## License

pod is licensed under the [copyfree](http://copyfree.org) ISC License.
