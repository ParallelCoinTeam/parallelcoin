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

kopach:
	go install -v
	pod -D test0 -n testnet -l debug -g -G 1 --lan --solo kopach

nodegui:
	go install -v
	pod -D test0 -n testnet nodegui

gui:
	go install -v
	pod -D test0 -n testnet --lan

guis:
	go install -v
	pod -D test1

guihttpprof:
	go install -v
	pod -D test0 -n testnet --lan --solo --kopachgui --profile 6969

guiprof:
	go install -v
	pod -D test0 -n testnet --lan --solo --kopachgui

testnode:
	go install -v
	pod -D test0 -n testnet -l debug node

testwallet:
	go install -v
	pod -D test0 -n testnet -l debug wallet