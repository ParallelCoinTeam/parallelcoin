# ec

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://godoc.org/git.parallelcoin.io/dev/pod/btcec?status.png)](http://godoc.org/git.parallelcoin.io/dev/pod/btcec)

Package ec implements elliptic curve cryptography needed for working with

Bitcoin (secp256k1 only for now). It is designed so that it may be used with the standard crypto/ecdsa packages provided with go. A comprehensive suite of test is provided to ensure proper functionality. Package btcec was originally based on work from ThePiachu which is licensed under the same terms as Go, but it has signficantly diverged since then. The btcsuite developers original is licensedunder the liberal ISC license.
Although this package was primarily written for pod, it has intentionally been designed so it can be used as a standalone package for any projects needing to use secp256k1 elliptic curve cryptography.

## Installation and Updating

```bash
$ go get -u git.parallelcoin.io/dev/pod/btcec
```

## Examples

- [Sign Message](http://godoc.org/git.parallelcoin.io/dev/pod/btcec#example-package--SignMessage)  
  Demonstrates signing a message with a secp256k1 private key that is first
  parsed form raw bytes and serializing the generated signature.

- [Verify Signature](http://godoc.org/git.parallelcoin.io/dev/pod/btcec#example-package--VerifySignature)  
  Demonstrates verifying a secp256k1 signature against a public key that is
  first parsed from raw bytes. The signature is also parsed from raw bytes.

- [Encryption](http://godoc.org/git.parallelcoin.io/dev/pod/btcec#example-package--EncryptMessage)
  Demonstrates encrypting a message for a public key that is first parsed from
  raw bytes, then decrypting it using the corresponding private key.

- [Decryption](http://godoc.org/git.parallelcoin.io/dev/pod/btcec#example-package--DecryptMessage)
  Demonstrates decrypting a message using a private key that is first parsed
  from raw bytes.

## License

Package btcec is licensed under the [copyfree](http://copyfree.org) ISC License except for btcec.go and btcec_test.go which is under the same license as Go.
