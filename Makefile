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
	pod -D test0 -n testnet -l debug -g -G 1 --lan --solo --kopachgui node

nodegui:
	go install -v
	pod -D test0 -n testnet nodegui

gui:
	go install -v
	pod -D test0 -n testnet -g -G 1 --lan --solo

guihttpprof:
	go install -v
	pod -D test0 -n testnet -g -G 1 --lan --solo --kopachgui --profile 6969 testnet

guiprof:
	go install -v
	pod -D test0 -n testnet -g -G 1 --lan --solo --kopachgui
