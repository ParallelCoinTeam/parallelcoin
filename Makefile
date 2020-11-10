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
<<<<<<< HEAD
	pod -D test0 -n testnet --lan --solo

guihttpprof:
	go install -v
	pod -D test0 -n testnet --lan --solo --kopachgui --profile 6969 testnet
=======
	pod -D test0 -n testnet -l debug --lan --solo

guihttpprof:
	go install -v
	pod -D test0 -n testnet --lan --solo --kopachgui --profile 6969
>>>>>>> refgui

guiprof:
	go install -v
	pod -D test0 -n testnet --lan --solo --kopachgui
