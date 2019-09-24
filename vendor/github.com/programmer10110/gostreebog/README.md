# gostreebog
Streebog (GOST 34.11-2012) hash implementation

## Overview
Standard GOST 34.11-2012 "Information technology. CRYPTOGRAPHIC DATA SECURITY. Hash-functions" came into the effect in 2013. The standard is based on previously published hash function Streebog.

This project includes implementation of the hash function in Go.

Propositions, comments, speedup or other improvements are welcome.

## Features
- 256-bit hash function
- 512-bit hash function

## Install
    go get github.com/programmer10110/gostreebog
    
## Usage

```Go
package main

import (
	"encoding/hex"
	"fmt"
	"github.com/programmer10110/gostreebog"
)

func main() {
	message := []byte("Streebog")
	fmt.Println(hex.EncodeToString(gostreebog.Hash(message, "512"))) //512-bit hash function
	fmt.Println(hex.EncodeToString(gostreebog.Hash(message, "256"))) //256-bit hash function
}
```
