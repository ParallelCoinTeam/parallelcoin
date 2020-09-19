# Parallelcoin

## all-in-one-everything for parallelcoin

The Plan 9 from Crypto Space hardfork is resuming development.

The hard fork includes the following new features:

- **New Proof of work hash function** - uses very large integer long 
  division to avert the centralisation of hash power by using
  the slowest and most expensive integer mathematical operation
  already being most optimally implemented in CPUs. GPU dividers
  are slower and no ASIC could do it faster.
  
- **9 way parallel prime based block intervals** - The blocks come
  in randomly anyway so why not make the schedule semi-random?
  This block product scheduler (difficulty adjustment) runs 9
  parallel block schedules using the new hash function where each
  different block has a different but regular block time with a
  different difficulty target and proportional block reward. 
  This allows a broader scale dynamic between small and larger 
  miners, who have different needs for payment regularity.
  The average block time is 18 seconds, which is sufficient
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
  windows for setting the password.
  
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
