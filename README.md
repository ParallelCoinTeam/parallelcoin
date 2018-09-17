<!-- [![Release](https://img.shields.io/gitlab/release/parallelcoin/parallelcoin.svg?style=flat-square)](https://gitlab.com/parallelcointeam/parallelcoin/releases/latest) -->

Parallelcoin integration/staging tree
====================================

Copyright (c) 2009-2014 Bitcoin Developers,
Copyright (c) 2013-2015 Parallelcoin Developers

Parallelcoin v1.2.0.0

What is Parallelcoin?
--------------------

Parallelcoin is an experimental new digital currency that enables instant payments to
anyone, anywhere in the world. Parallelcoin uses peer-to-peer technology to operate
with no central authority: managing transactions and issuing money are carried
out collectively by the network. Parallelcoin is also the name of the open source
software which enables the use of this currency.

License
-------

Parallelcoin is released under the terms of the MIT license. See `COPYING` for more
information or see http://opensource.org/licenses/MIT.

Building
--------

A full GNU C/C++ toolchain, Qt5 headers and BerkeleyDB 4.8 are required to build parallelcoind and parallelcoin-qt. Miniupnpc is optional but recommended, it is enabled by default in the build.

This branch of the repository only the linux build has been tested and there is a script in the `src/` directory called `linux-build.sh` which builds everything, installs the desktop files for the GUI wallet and then cleans up after itself.