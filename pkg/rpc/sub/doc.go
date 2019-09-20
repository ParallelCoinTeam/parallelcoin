// Package sub is a short message publication/subscription library that uses UDP transport, Reed Solomon erasure coding, ed25519 EC signatures for tamper-resistance, for allowing clients to subscribe to updates from a server for time-sensitive messaging, written to implement a low latency work delivery system for Parallelcoin miners.
//
// To prevent retransmits for messages up to 3kb in size, data sent in a burst as 9 packets containing a 9/3 Reed Solomon encoding such that any 3 packets received guarantee retransmit-less delivery, covering the worst case for packet loss and corruption over a network
//
// Payload can be encrypted via AES-256 encryption using a pre-shared key known by both ends to function as both access control and security against eavesdropping and spoofing attacks.
//
// Authentication of data is done using an ED25119 EC key for which each known endpoint has shared the public key as part of the subscription request
package sub
