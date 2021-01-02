#!/usr/bin/make -f

build:
	go build -v

windows:
	 go build -v -ldflags="-H windowsgui"

tests:
	go test ./...

kopachgui:
	go install -v
	pod -D test0 -n testnet -l debug --lan --solo --kopachgui kopach

testkopach:
	go install -v
	pod -D test0 -n testnet -l trace -g -G 1 --lan kopach

testnode:
	go install -v
	pod -D test0 -n testnet -l debug --solo --lan node

nodegui:
	go install -v
	pod -D test0 -n testnet nodegui

gui:
	go install -v
	pod -D test0 -n testnet --lan

guis:
	go install -v
	pod -D test1 --minerpass pa55word

guass:
	go install -v
	pod -D test0 --minerpass pa55word

resetwallet0:
	pod -D test0 -l trace --walletpass aoeuaoeu wallet drophistory

resetwallet1:
	pod -D test1 -l trace --walletpass aoeuaoeu wallet drophistory

guihttpprof:
	go install -v
	pod -D test0 -n testnet --lan --solo --kopachgui --profile 6969

guiprof:
	go install -v
	pod -D test0 -n testnet --lan --solo --kopachgui

mainnode:
	go install -v
	pod -D testmain -n mainnet -l info --connect seed3.parallelcoin.io:11047 node

testwallet:
	go install -v
	pod -D test0 -n testnet -l trace --walletpass aoeuaoeu wallet

mainwallet:
	go install -v
	pod -D testmain -n mainnet -l trace wallet

teststopkopach:
	go install -v
	go install -v ./pkg/util/logi/pipe
	pipe pod -D test0 --pipelog -l trace --walletpass aoeuaoeu -g -G 1 --solo --lan kopach

teststopnode:
	go install -v
	go install -v ./pkg/util/logi/pipe
	pipe pod -D test0 --pipelog -l trace --walletpass aoeuaoeu -g -G 1 --solo --lan node

teststopwallet:
	go install -v
	go install -v ./pkg/util/logi/pipe
	pipe pod -D test0 --pipelog -l trace --walletpass aoeuaoeu -g -G 1 --solo --lan wallet
