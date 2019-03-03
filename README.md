<!-- [![Release](https://img.shields.io/gitlab/release/parallelcoin/parallelcoin.svg?style=flat-square)](https://gitlab.com/parallelcointeam/parallelcoin/releases/latest) -->

> ### Important update
>
> Herein lies the most up to date version of the original Parallelcoin server, mainly updated to more easily build on current versions of Linux. 
> 
> The replacement, based on conformal's [btcd](https://github.com/btcsuite/btcd), is found in this org in the [pod](https://github.com/parallelcointeam/pod) repository.
>
> The new version is significantly improved in performance with the channel based logging system, tweaks in network parameters that reduce latency delays with the peer manager, and combines in fact three of the btcsuite applications, the CLI controller, the full node, and wallet server, which have been merged into one application, mirroring the design of the original, but brought into the present.
> 
> The first draft of the combined single binary merges the configuration systems and has a mode to launch wallet and full node together within one runtime to enable the more familiar combined wallet/node as used in Bitcoin and its many close derivatives, for both web wallet/paymennt processing systems and the like, as well as forming the base of the GUI front end.
> 
> Current status is the configuration system is currently being rewritten ([see here](https://github.com/parallelcointeam/tri)) as the complexity of the configuration was excessive and several configuration bugs and a mass of quite difficult to maintain very long and repetitive code, prevents the current application to properly launch in the new combined 'shell'. The hard fork logic to enable a mining protocol is drafted but not fully unwrinkled, and the GUI has been put on hold until the hard fork and configuration issues are resolved.
>
> The new configuration library dramatically simplifies and builds up from the first draft, and will be completed soon and the whole configuration/CLI configuration system will be replaced with something far simpler and more powerful.
>
> Leaving the best to last, the fairly well tested new mining protocol uses 9 distinct, though composed of many common elements, is extremely complex and utilises the new Cryptonight7v2 hash function, the Argon2i KDF function (similar to Scrypt but far more memory- *and* processing heavy) as well as 9 different algorithms including the original sha256d and scrypt hashes. It is likely that the new scheme, requiring pipelining between steps, should be very difficult to optimise and perform equally poorly on CPU, GPU and ASIC/FPGA, as some parts stress cache, some stress memory, and some stress raw processing, no one device will perform well at all three types of work, thus equalising more between different platforms, and impeding especially ASIC development.
> 
> Mining protocols are only half about the actual hash functions, the actual work, but harder to see and more important is the adjustment and block times, in the current alpha runs a 9 second block time, and uses 4 different averaging algorithms, all-time simple average, trailing simple average, exponentially decaying all algorithm block time and a per-algorithm averager that aims to keep algorithms holding their correct average, and lastly, exponential-limited difficulty reduction that reacts less to lengthening of the average to produce a significant reduction in block times under 3 seconds.
>
> The final piece is, along with the short average block time at 9 seconds, dramatically disadvantages efficiency improvements that increase startup-latency in a trade-off that does not work with much shorter blocks, the initial CPU miner and custom built reliable UDP, with forward error correction, running a mining work delivery system that aims to further reduce the response time to new blocks.
>
> As a whole, this arrangement of 
> 
> Initial tests show that this complex unified and partially separating difficulty adjustment achieves a divergence under 20% from recent average times from target and should produce a reduced long-range divergence - as this old server here in this repository has, over the time since genesis in early 2014, has slowly stretched out towards 12.5 minute blocks.
> 
> ***The initial public beta is expected to commence in early april and when the remaining baby problems are ironed out, a hard fork block height will be announced.*** Stay tuned!


Parallelcoin integration/staging tree
====================================

Copyright (c) 2009-2014 Bitcoin Developers,
Copyright (c) 2013-2015 Parallelcoin Developers
Copyleft ~~(c)~~ 2018 Parallelcoin Team

Parallelcoin v1.2.1.2

## What is Parallelcoin?

Parallelcoin is an experimental new digital currency that enables instant payments to
anyone, anywhere in the world. Parallelcoin uses peer-to-peer technology to operate
with no central authority: managing transactions and issuing money are carried
out collectively by the network. Parallelcoin is also the name of the open source
software which enables the use of this currency.

Parallelcoin was one of the first coins to allow more than one algorithm for Proof of Work, namely Scrypt and SHA256, as from Bitcoin and Litecoin. This repository contains the best current working and updated version that builds for now at least on Linux and a working parallelcoin installation will be required to import the wallet initially.

## License

Parallelcoin is released under the terms of the MIT license. See `COPYING` for more
information or see http://opensource.org/licenses/MIT.

Further changes visible in the repository history are covered by the Unlicence (pure public domain) though they are small in number and mainly necessary to enable building with current environments and toolchains.

## Building


### Linux

A full GNU C/C++ (gcc) toolchain, Qt5 headers, boost and BerkeleyDB 4.8 are required to build parallelcoind and parallelcoin-qt. Miniupnpc is optional but recommended, it is enabled by default in the build and it won't work without it.

On Arch Linux this means you will need these packages:

    base-devel
    boost
    boost-libs
    miniupnpc (we see no reason why this should not always be in the build as it is small and helpful to many users)
    db4.8
    qt5 (select all when asked which packages to install)

This branch of the repository only the linux build has been tested and there is a script in the `src/` directory called `linux-build.sh` which builds everything, installs the desktop files for the GUI wallet and then cleans up after itself.

The installer also places a systemd service file to run the parallelcoin client as a service. Note that if it is running this way the Qt client will refuse to run without further configuration (maybe?)
