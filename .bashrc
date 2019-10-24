#!/bin/bash
alias goland=$HOME/bin/goland
export GORACE="strip_path_prefix=$GOPATH/src/github.com/p9c/pod/"
export GOFLAGS=""
export GOTMPDIR="/dev/shm"
export CGO_CFLAGS="-g -O2 -w"
export GOBIN=$HOME/bin
export PATH=$HOME/go/bin:$HOME/goland/bin:$GOBIN:$HOME/android-studio/bin:$HOME/flutter/bin:$HOME/bin:$PATH
export GOPATH=$HOME
export GOROOT=/usr/lib/go
export GO111MODULE=off
