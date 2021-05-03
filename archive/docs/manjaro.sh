#!/bin/bash
cd ~
yay -Syy wget git base-devel 
wget -c https://golang.org/dl/go1.14.13.linux-amd64.tar.gz
sudo rm -rf ~/go
tar zxvf go1.14.13.linux-amd64.tar.gz
cat <<EOF >> ~/.bashrc 
export GOBIN=$HOME/.local/bin
export GOPATH=$HOME
export GOROOT=$HOME/go
export PATH=$HOME/go/bin:$PATH
EOF
mkdir -p src/github.com/p9c
cd src/github.com/p9c
git clone https://github.com/p9c/pod.git
cd pod
source ~/.bashrc
go mod tidy
make -B stroy
stroy guass
