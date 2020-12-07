# ![](https://raw.githubusercontent.com/p9c/pod/master/pkg/gui/logo/logo.svg) parallelcoin pod

## all-in-one-everything for parallelcoin

The Plan 9 from Crypto Space hardfork is resuming development.

The hard fork includes the following new features:

- **New Proof of Work hash function** - uses very large integer long 
  division to avert the centralisation of hash power by using
  the slowest and most expensive integer mathematical operation
  already being most optimally implemented in CPUs - long division. 
  GPU dividers are slower and no ASIC could do it faster.
  
- **9 way parallel prime based block intervals** - The blocks come
  in randomly anyway so why not make the schedule semi-random?
  This block product scheduler (difficulty adjustment) runs 9
  parallel block schedules using the new hash function where each
  different block has a different but regular block time with a
  different difficulty target and proportional block reward. 
  This allows a broader scale dynamic between small and larger 
  miners, who have different needs for payment regularity.
  The average block time is 36 seconds, which is sufficient
  for many retail operation types.
  
- **Multi-platform touch-friendly wallet GUI** - A responsive and 
  simple user interface for making and viewing transactions
  as well as an inbuilt block explorer, configuration and mining
  controls, available on Windows, Mac, Linux, Android and iOS.
  
- **Simple zero configuration multicast mining cluster control 
  system** - Using a simple pre-shared key it is easy to add nodes 
  and mining worker machines to a cluster. Simple redundancy by
  the use of multicast one can just run one or more controller
  nodes and the workers just listen to whoever is also using the
  same key. A customised live linux USB will be available for 
  both functions with an easy configuration file readable by
  windows for setting the password and payment addresses.
  
## Future directions

### Coinbase Rosetta Integration

A native integrated Coinbase Rosetta server implementation will
allow easier integrations for exchanges and multi-chain wallets
and other applications.

### Finality

In order to connect Parallelcoin directly to the 
[Cosmos](https://cosmos.network/) "Internet of Blockchains" 
the Nakamoto consensus, first used in bitcoin, having
probabalistic finality has to be augmented to enable its 
effective integration.

To do this in a practical way without blowing up the size of the
data required to store the chain, for this the signature algorithm
will be changed from 
[secp256k1](https://www.cryptoglobe.com/latest/2018/07/bitcoins-schnorr-upgrade-could-be-the-most-significant-change-since-segwit/)
to Schnorr signatures using curve25519, enabling very large numbers of 
cosigners in a multiple signature. The transaction format will be
revised to account for this change.

With Schnorr multisigs each transaction can have a chain of signatures and
the signature is stored with the transaction on the chain, via epidemic
transmission the signature chains form and the longest one is recorded in the
next subsequent block in a finality section in the transaction merkle tree,
thus providing total finality in the next block.

### Integration of mempool and chain database

When transactions are propagating with lots of cosigners finalising them
the transaction data can be considered immediately current even before it
gets into a block so the chain servers can return them as results in
RPC queries.

With tight integration of mempool and chain data, and immediate
finality the most cited disadvantage of proof of work as a means
of distributing new tokens randomly and evenly is eliminated.

### Increasing precision of token

Parallelcoin has a very small number of tokens in total, and as such
the need to extend its precision will arise sooner than tokens with
tens and millions more tokens. As such, the current 64 bit long fixed
integer format for currency values will be extended to a variable 
length up to 128 bit long value, allowing plenty of room for more 
zeros.

It is common practise in business financial databases to use 128 bit 
precision, and with variable length encoding many numbers will
still end up being only 64 bits long anyway.

### Upgrading to a latency minimising network transport

Latency is a major structural limit on distributed systems as they 
differ from centralised systems by the greater number of hops that
data will make between nodes, and latency delays stack up like
cars in gridlock.

A currently in-development, latency minimising, bad-network-aware
network transport will be used to replace all inter-node internet
data transmission that strongly avoids retransmit latency and 
copes with transient changes in network paths more reliably than
standard TCP/IP sessions.

### Relay routing

In order to reduce the attack surface of the network, all transactions
should be broadcast over a chain of proxies in the same manner as
with Tor network, obscuring the physical origin of the transaction.

There has been more than a few instances already of robberies taking
based on transaction origin and without going all the way into full
privacy protection the quickly expiring value of location data is 
eliminated from the options of attackers intending to rob a user
of the network.

### Mimblewimble private transactions

Ultimately all transactions should be 100% private, and the best 
design for protecting the privacy of transactions currently
available is the Mimblewimble protocol, which uses advanced
cryptographic techniques to provide only provable confirmation
without giving either the value or the sources of the transaction
to potential attackers.

Mimblewimble is a complex addition to the protocol so obviously will
be done quite a bit later than most of the other items in this list.

### Integrated agorised secure forgery-resistant low latency message relay network

Ultimately entirely replacing the mempool in exchange for a complete
secure, censorship resistant self healing peer to peer message exchange
network system, with a means to charge for relay service without
revealing the identity of the user requesting relay service.

With such a system in place it is simple to enable the negotiation of 
any kind of collaborative workspace for users, from instant messaging
to shared creative environments from software repositories to org charts,
to audio and video conferencing and anything that users of the network
come up with.

#### Namespace

As part of the relay network there should be a market in namespaces to 
make human-friendly labels for resources on the network.
