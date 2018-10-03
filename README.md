<!-- [![Release](https://img.shields.io/gitlab/release/parallelcoin/parallelcoin.svg?style=flat-square)](https://gitlab.com/parallelcointeam/parallelcoin/releases/latest) -->

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