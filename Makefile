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

node:
	go install -v
	pod -D test0 -n testnet -l debug -g=false -G 0 --lan --solo --kopachgui node