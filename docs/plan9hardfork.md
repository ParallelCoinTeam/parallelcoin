# ParallelCoin

## Plan 9 Hard Fork Specifications

### 1. Hash Function

A key goal of the Plan 9 Hard Fork is to make a chain that has a very large peer to peer network composed of mostly medium and small miners, in order to reduce the loss of security from mining being under the control of a small number who make up 80% or more of blocks.

The reason why this aggregation is bad for security is that in any distributed system, the more centralised a section becomes, the more vulnerable the network becomes to extensive cascading failures, or a single lever for an attacker either upon the chain itself or against some set of users. There is an obvious economic issue as well - dominant miners on a network usually also have a lot of reserve tokens and can create the conditions for higher volatility.

The hash function consists of slicing up the block header in two pieces, shuffling and rearranging the bytes composing the data, squaring the jumbled halves and then multiplying them together, and the most important step, dividing the result with a sliced up version of the original block.

The process is done twice and then the final data...
