#!/bin/bash
cd $HOME
wget -c https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
sudo rm -rf go
tar xvf go1.16.2.linux-amd64.tar.gz
cat >> $HOME/.bashrc <<- EOM
export GOPATH=\$HOME
export GOROOT=\$GOPATH/go
export GOBIN=\$GOPATH/bin
export PATH=\$GOBIN:\$GOROOT/bin:\$PATH
EOM
source $HOME/.bashrc
mkdir -p $GOPATH/src/github.com/p9c
cd $GOPATH/src/github.com/p9c
git clone https://github.com/p9c/pod.git
cd pod
make -B builder
