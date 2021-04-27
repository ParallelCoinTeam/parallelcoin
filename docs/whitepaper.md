![](https://raw.githubusercontent.com/p9c/logo/main/logo_nocircle128x128.svg)

# The Parallelcoin Pod White Paper

###### by

### David Vennik

#### April 2021, Novi Sad, Serbia

---

## Open, Programmable, Distributed Business Administration Systems for All

The greatest challenge for cryptocurrency and related open and public
distributed network systems is integration, but this is not unique to it, it is
the central criteria that has determined the success of software development
companies since the 1970s.

In any company with more than a handful of employees, there is likely at least
one user who can at least construct a multi-page spreadsheet with basic sums,
averages and so forth, to generate financial reports and evaluate projections.

The number one software company in the world for the majority of the last 50
years is Microsoft, whose key business strategy was in building business data
administration systems that were usable by minimally trained and skilled users.

The one thing that Microsoft understood best was that businesses need
interoperability of data and granular depth of access to programmability of the
applications, to fit the business's needs, organisational communication
structure and protocols, and the business operators' philosophies.

Cryptocurrency represents the first and most foundational element of a complete
business network administration system, a financial ledger, that can be trusted
without needing to trust any individual or group participating in the network.

## 1. Evolution of Open, Distributed Business Administration Systems

### Bitcoin

#### Permissionless, Accountless, Trustless, Reliable Distributed Financial Ledger

Bitcoin was the first distributed database system that is reliable enough to be
trusted while distributing the risk of malicious network participants such that
the cost of bad behaviour is higher than the potential gain.

It has confused pundits and experts alike for some time what it really is about,
some focus on its use of artificial scarcity, some focus on the inherently
conservative governance created by a system with purely staked participation in
its development direction, due to their ability to embargo updates to the
protocol, causing a fork and damage to the interests of all with such stakes at
risk, while still others focus on its survivability compared to more brittle
centralised systems.

> ### The one thing all businesses have in common is money

It is not necessary to deal in either borrowing or lending, and outside of this
is the laws of business conduct, such as contracts of sale and lease, which must
have human adjudication.

Money, however, before the concept of distributing some or all of the paper
records started to emerge with, first, printing presses, then typewriters, and
finally computers, evolving from punch cards through spinning disks to
microchips with billions of operations MTBF, as a distributed ledger is 
susceptible to being altered to benefit one group against the entire rest of 
the marketplace who use the currency issued by this group, and trusting the 
veracity of this group's clear privilege to author any record it wants recorded.

The greatest source of power within a distributed ledger is that of authorship.
Who has the authority to cause a particular record to become canonical and thus
immutable, and provide the service thereby of a stable base on which to do
economic calculations when deciding how to allocate working capital in a
business... because that common data affects every decision even if only in a
subtle way?

At the centre of the proposal for the protocol proposed by the pseudonymous
Satoshi Nakamoto was the idea of creating a distributed lottery for the right to
be the author of a batch of records in this financial ledger, thus distributing
the risk of false records being made canonical, and using probability to
progressively raise the cost of having an altered record become trusted.

### Ethereum

#### Programmable Distributed Database Application Layer

The programmability tends to be the attention grabbing part of the description
of this system, but it is not this that makes it both an innovation, and the
basis for further refinements in the scope of what can be done to reduce error
and automatically defend against fraud in consensus business records.

But what is more important is also more subtle - it is about allowing the
creation of arbitrary new recompositions of user input and data stored in the
financial database and other databases such as the 'gas' that evaluates code by
its processing costs and data weight to ensure that blocks are produced within a
limited span, and the network cannot be gridlocked by spammy processing, the
first and most important of the many little databases that an ethereum node
stores.

The emerging downside of the technology is that the encoding of applications
into immutable blockchains is quite expensive in space, taxes the serial
processing nature of the Nakamoto Consensus upon which it has been unable to
completely separate itself (due to the unarguably extremely strong safety of
transactions beyond a certain age), and led to several strategies that followed.

### Tendermint and Cosmos

#### Separating Consensus and Execution

The Tendermint protocol is a derivative of Practical Byzantine Fault Tolerance,
a distributed database protocol first proposed in the ACM
[Practical Byzantine fault tolerance](https://dl.acm.org/doi/10.5555/296806.296824)
paper by Miguel Castro and Barbara Liskov.

The refinement Tendermint adds is the ability to alter the database of certified
proposers, since the original design assumes that there is no malicious nodes,
which is a dangerous assumption in an open, or as the correct way to describe
it, fluid-membership consortium with both automated and human-operated methods
of preventing the altering of the data, and the finalization of the records
immediately after publication and before the next batch is processed, called '
Immediate Finality'.

Like all distributed network systems, it is impossible to solve every problem
perfectly, there is always one or more element that has a limitation that
provides the ability to give strong guarantees about the other, mutually
exclusive properties in the protocol.

#### Consensus is A Protocol Level

A key element of the Tendermint design is that it separates the validation from
the replication. This creates an agnostic base level to the protocol that can be
then used to carry any kind of data in a batched, cryptographically chained
block format.

This greatly simplified the process of building a new protocol - it was not
necessary to build a scripting engine into the server, but instead to be able to
write it in a language that compiles down to fast binary code, and the time cost
of taking a database application and turning it into a distributed database
application has been greatly reduced.

#### Finality

At the very head of a blockchain whose head is chosen by the highest cumulative
total of work along its path, there is the possibility and both malicious and
accidental, network- and machine- failure caused corruption of blocks that can
be mined on prior blocks in parallel, and it is only the high probability of the
majority of miners being honest enough to maintain the value of the network,
this is an assumption that depends on total domination of the market in miners.

But outside of Bitcoin, the number of miners per chain is far smaller and there
can be overlaps between mining capability and several chains that can allow
concentrated hashpower attacks, that have occurred several times previously.

The necessary cost of a delay before strong confirmation of a transaction is
available makes the operation of applications with a lower latency requirement
impossible. The most extreme form is realtime massive multiplayer games, which
further add the existence of multiple parallel 'spaces'. Limiting a whole 
ecosystem to the capacity of one, however high spec, server machine, greatly 
limits the possible applications where updates need to be more timely.

#### Transportability of Ledgers

Where Ethereum and other chains have considered sharding, in order to
parallelise the processing of transactions, the Cosmos project launched the
concept of instead making a chain for the application, and allowing a third
layer to the distributed system, where the ledgers are able to exchange between
multiple networks transparently, tokens are locked on one and proof of the 
lock is used to issue new, identical tokens on the other chain.

Parallelcoin is building a unified application environment that allows 
programmability at all levels of the system, from protocol to interface.

## 2. The Parallelcoin Pod

The central concept behind the Pod can be expressed in well known eastern
european Matryoshka dolls. Everything in the system has a container type that
can carry a message intended for a different processing system, with the only
common limitation of requiring a public key or address and the signature on it
to prove authenticity.

This is essentially the same as the Tendermint model, in that transactions are
generalised, and the validation is separate though the validator governs the
re-distribution of a received message or its processed product.

Where Pod goes further is not in only the base layer, but rather in a holistic
view of the application starting from the user input and ending at the network
protocol.

Most blockchain wallets have some kind of RPC console built into them, and in
this console, with varying degrees of interface sophistication, one can query
the blockchain's databases and get back results.

Hard-core, early adopters originally accessed the chains only this way to begin
with, but the sophistication of these interfaces has been pretty much halted at
a level that is far below that of even basic CLI interfaces like the `sh`
command interpreter.

#### Even young children can construct a basic script

It could even be argued that in fact, narrative, the serialisation of a series
of expressions into a coherent message that conveys some kind of information is
the primary skill and activity of the human brain, in the processing of
language.

In order to perform any task, a person must first learn the program of the task.
What things come first, which things can occur in parallel or have non-competing
delays, and the productivity of a worker exactly relates to their ability to
learn to perform the task in the way that straightens out the nebulous cloud of
requirements into a directed acyclic graph showing the various steps in the
process and includes parts where parallelisation is visible, such as waiting for
dough to rise, or a solution to settle out solids.

#### From small, open things, great and complex things grow

Not every need of a person wanting to keep business records or conduct
transactions of business records with another business is complicated. The web
has exploded in the main due to the fact that it exactly allows aggregation of
data. The number one business being search, since any data set's backbone is the
index hash that all records bear, without it, you simply cannot find it. 
Someone has to build that index and it is a utility that is universal to all 
business.

By making scripting, and eventually programming, from network protocol to GUI,
accessible to the user only limited by their technical skill, and making the
blockchain a first class citizen in the programming environment, the Pod will do
for cryptocurrency what Microsoft did for word processing, spreadsheets and
databases.

### 2.1 Phase 1: The Parallelcoin Plan 9 Hard Fork

The first step in the process is capitalisation. The founder of the Community
Takeover, Djordje Marcetin, is a significant holder of the rare Parallelcoin DUO
token, and has funded and developed a significant part of the body of work that
is going into the launch of the upcoming protocol upgrade.

Equity markets have a very high lower bound for entry, and in any case require
some kind of demonstration of the capabilities of the team to produce the
promised systems. Thus, the upgrading, securing, and promotion of the use of the
old platform and its tokens is a means to also secure development funding for
the early adopters who are working also to bring a new and compelling product to
the market that establishes a new baseline for what computing in the age of the
internet should be about.

The hard fork includes:

- #### CPU oriented proof of work
  that depends on very large integer long division as a bottleneck step, thus
  proportionally disadvantaging generic, programmable and custom silicon
  production that is not specialised for this expensive, iterative mathematical
  operation, and more scarce and thus narrowing the decentralisation, and 
  thus security, caused by this concentration, as seen in the established 
  ASIC based PoW coins.

- #### 9 Parallel Block Intervals with Network Cost Orphan Minimization

  The block schedule is composed of 9 parallel schedules based on powers of 2,
  which each provide a reward that is in proportion with their consensus target
  average, with the base interval at 18 seconds, producing a resultant 
  average of 9 seconds between blocks.

  This provides a more accurate census of the miner population, needed for
  accurate difficulty adjustment, than a single interval, as each new block on
  any of the intervals also triggers the recomputing of the targets of the other
  intervals, without further complicating the data in the on-chain consensus, it
  is derived through averaging of the history only, over several common and
  interval-specific timespans extending up to 3600 times the longest, being 2304
  seconds or 38.4 minutes, or 96 days.
  
  ##### Orphan Minimization

  Because the block time is short, many orphans are mined. To combat this, 
  in addition to the standard cumulative work total, (the lowest sum of all 
  block hashes in the chain path), the value of the block is a criteria that 
  modulates the fork evaluation.
  
  In addition, with varying block rewards, it is natural that hash power 
  attacks on the difficulty adjustment will target the longest intervals, 
  and thus be more likely to be the blocks of miners performing hashpower 
  attacks.
  
  When these lower reward blocks are added to the head of the chain, the 
  difficulty of the higher reward blocks is recomputed, ensuring that a 
  spike in numbers of these blocks is reined in quickly, providing the 
  benefit of more frequent adjustments with an optimised cost for the 
  network value (dilution of currency).
  
  Thus, rather than only determining the best block by work, it is the work 
  of the most frequent blocks that get priority position to join 
  the chain, in proportion with their target and historical frequency, by 
  expanding the work value formula to add an equal percentage to the work 
  inverse of the relative times of the blocks. 
  
  So all else being exactly equal, the lowest value block is chosen by a factor 
  the percentage out of 256, the maximum block interval's differential, up 
  to 100%. This thins out the number of orphans that can form a fork in an 
  economic way for the holders of the coins, specifically, and favours also 
  a fairer distribution amongst small miners of the most frequently occuring 
  block versions.

- #### Zero configuration multicast local area network

  Adoption is the key goal of Plan 9 Crypto, and participating in the mining 
  of the base token is something we want to see as wide as possible. More 
  blocks means more wide distribution of tokens as well as improving the 
  liveness of the network (time between updates).
  
  To facilitate this, configuring a standard mining farm is as easy as 
  deploying a wireless network, a pre-shared key acts as both signal 
  security and channel separation for the miners, and jobs are streamed 
  regularly in advertisments for miners to come online and start work on blocks.
  
- #### Multiplatform GUI wallet built in pure Go

  Go is a programming language generallly associated with scalable 
  distributed networked database systems, and has had quite wide adoption in 
  the cryptocurrency space. It has a key advantage over the other languages used 
  in that it has a cooperative, work stealing scheduler that coordinates 
  lightweight threads communicating over atomic FIFO queues called 'channels'. 
  
  This programming model, known as Concurrent Sequential Processes, allows 
  extremely low latency responses to data processing with widely varying and 
  unpredictable volume and frequency. This makes it great for dividing up 
  fluctuating loads of processing in parallel, but also makes it easy to 
  build user interfaces that represent many concurrent processes.
  
  The Pod is a single binary but contains run modes for all functions 
  it performs, that run as independent processes, sharing data with other 
  components over sockets and pipes.
  
  The GUI interface is designed to be usable and efficient when used on any 
  device, whether it be a mobile phone, tablet, laptop, or desktop computer. 
  Given the necessary limitations of size, all functions are accessible on 
  all platforms.
  
## 2.2 Phase 2: You will be assimilated!

### The Pod Console

The key technology that is planned to be developed once the new protocol is 
in place, and the hard fork activated is an open, accessible and intuitive 
system for creating scripts, and eventually full programming logic, that 
treats command line style invocations as equal to Go statements and 
expressions, and facilitating grouping these lists into functions, 
automating the specification of variables, parameters and return values.

This starts out as a simple RPC console, as will be found in the Hard Fork 
GUI wallet, and stepwise adds features to it that allow users to build 
applications and supporting libraries and servers. Between an RPC API call 
and a Go function invocation there is a 1:1 correspondence between the 
variables involved and input values and their types, and so, recognising 
this, the Pod Console will allow anyone to easily create simple queries 
combining multiple data sources and produce a report that can be updated 
from data connected to the chain.

#### Build anything

Rather than separating end user and server and peer side programming as is 
convention in the industry, with, for example with ethereum, the server 
written in C++, Go or whatever other language, while the scripts are written 
in Solidity, the Pod Console allows you to use the same language to define 
everything from the user interface to the format of protocol packets in a 
peer to peer data aggregation application that can be deployed with a simple 
command line invocation or automated through a GUI interface.

This programming environment natively deals in a log format consisting of 
the actions that can be immediately executed or deferred, grouped into 
functions, functions grouped into a package and either an application or 
supporting library created to facilitate the building of more applications. 
Not just smart contracts but smart user interfaces as well.

The source code that is generated out of the log of commands input by the 
user and structured into functions, types, and executable programs is stored 
in its original Directed Acyclic Graph form, which can then be played back 
like a script to produce the state of the application source code at a given 
point, like a commit in a Git repository.

These logs are bundled into updates on IPFS stores and by URLs including 
version specifications, once an application is developed and tested, it can 
be then immediately upgraded or installed by other nodes on the network.

#### Ad-hoc integration with external data sources

Since any code that can be integrated with Go can be used in this 
environment, it then becomes possible to create interfacing layers that 
encapsulate foreign server applications whether in Go, in other languages or 
only by their binaries.

From this basis, applications can require access to these services, which 
can then be automatically installed and once synchronised and ready for 
queries, new services can be created that depend on multiple data sources. 

#### Converting data sources into an Oracle

Unlike Tendermint/Cosmos, these data sources do not necessarily have to be 
other, deterministic types of sources, ie blockchains, by building a 
consensus around the external data source, and thus producing a trustworthy 
ledger of third party information that can be used as a data source in place 
of the third party's own outlets.

Once the data is certified by the oracle, it becomes deterministic, shared 
data for applications to connect this data feed to blockchain based 
applications.

### Bootstrapping

The first phase of developing the Pod Console aims to first quickly make the 
console itself stored in the format used by the console and editable. 
Similar to the bootstrap step in compiler construction, this step enables 
the platform to quickly advance to being changed to suit how it suits being 
changed.

The first milestone is to have created a native Console Log that builds the 
Console.

Once the system understands itself, it can be rapidly extended in any 
direction whatsoever. It will have a scheme for converting between API/CLI 
call format, and the native Go function syntax that is then compiled to be 
executed, and the user can chain lists of commands in any of the implemented 
parameter syntax formats, symbols given names, functions and types grouped 
together automatically and turned into Go packages that can then be shared 
immediately over IPFS.

From this point onwards, Plan 9 Crypto will pivot to focus on developing the 
platform with examples, integrated documentation, peer to peer hosting of 
source code and assets and one or several proof of concept projects.

### Proof of Concept 1: Inventory, Logistics and Compliance for Construction Industry

An application part way into development with a basic demonstration allows 
the creation of public records relating to materials, sources of materials, 
market current offered inventory, and automated generation of documentation 
required for compliance with government regulations to streamline compliance 
and reduce compliance administration costs to the minimum possible, leaving 
more capital aside to apply to other business purposes.

Through the application of the principle of designing protocols that 
encapsulate, composit and present data from one or many sources, the detail 
level of the data can eventually reach a point where it becomes effectively 
an authoritative resource on the information about everything related to 
governance of the industry and the facilitation of trade and discovery of 
profit opportunities and cost-savings opportunities that better connections 
can provide.

#### Implementation Overview

In order to implement this system, there is the creation of notary services, 
which certify and host the certified content on IPFS content addresses, 
built on several extensions to the ParallelCoin ledger for creating a market 
in notary services by providing a strongly immutable record of these 
notarisations.

The system contains three sections:

1. Realtime market data, offers for goods and services, offered available 
   inventory, prices and payment methods
   
2. Storage of public data required for legal compliance related to projects 
   including all of the technical specifications each different regulation 
   requires, with automated feeding of the data from inputs into the market 
   system and the settlements of these activities providing transparency and 
   compliance to users of the system operating in the industry.
   
3. A platform for the certifiable record of the complete projects from 
   inception, to the parties involved, the physical locations, the materials 
   used, the plans and engineering specifications, and to whatever degree 
   desired, additional media such as video, 3d models, photographs, text and 
   composite documents containing all of the above.

The Plan 9 Hard Fork protocol will be extended with several important 
features required for this service:

- Schnorr multisignature Public Key Infrastructure

  In one step using Schnorr signatures instead of Elliptic Curve signatures, 
  the data weight of multisignature transactions is reduced to a uniform 32 
  bytes, and facilitates the creation of a masternodes system that creates 
  block finalisation transactions from the given consensus masternode set.
  
- Mineable, exchangeable masternode tokens with a regular emission rate that 
  can be traded and deployed by masternode operators to sign the 
  certifications that allow the network to 100% finalise transactions, 
  within an average of 9 seconds of their publication.
  
- Integration with multiple naming services to facilitate easier access to 
  resources for humans while solving the sybil attack problem through 
  enforced scarcity, a consensus database of names and user-configurable 
  metadata used to access the named resources.
  
- Development of simple form GUI interface systems that combined with wallet 
  keychain management enable the creation of data entry for updating 
  inventory, service availability, and other elements of the 3 tier 
  architecture of the system, such as browsing market offers, compiling 
  historical data into charts, confidential communications between 
  businesses and secure mutually accessible related data storage and backup, 
  and so on.
  
## Proof of Concept Zero: Codename Jorm

In the preliminary work required before the development of Proof of Concept 1, as a testbed and initial implementation, the same protocol upgrades will be first deployed to create an oracle that aggregates many sources of data relating to cryptocurrencies.

This system will then become a service that can be accessed from any other application to retrieve data that can be trusted to be identical to the original source, and thus as a historical ledger of information updates, an Oracle that can be used for all distributed applications needing access to this data.

The second step is creating interfaces that bridge between the Oracle and the write-side of these data sources, subsuming the entire space into a portal that can be come a one-stop-shop for all information and eventually building interfaces that allow the arbitrary merging of services and data sources into a new application composition.
