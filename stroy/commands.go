package main

var commands = map[string][]string{
	"build": {
		"go build -v",
	},
	"headless": {
		"go build -v -tags",
	},
	"windows": {
		`go build -v -ldflags="-H windowsgui \"%ldflags"\"`,
	},
	"tests": {
		"go test ./...",
	},
	"kopachgui": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l debug --lan --solo --kopachgui kopach",
	},
	"gui": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet --lan",
	},
	"guis": {
		"go install -v %ldflags",
		"pod -D test1 --minerpass pa55word",
	},
	"guass": {
		"go install -v %ldflags",
		"pod -D %datadir -g 1 -G=false --lan --minerpass pa55word",
	},
	"guihttpprof": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet --lan --solo --kopachgui --profile 6969",
	},
	"guiprof": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet --lan --solo --kopachgui",
	},
	"mainnode": {
		"go install -v %ldflags -tags headless",
		"pod -D testmain -n mainnet -l debug --addpeer seed3.parallelcoin.io:11047 node",
	},
	"mainwallet": {
		"go install -v %ldflags",
		"pod -D testmain -n mainnet -l trace wallet",
	},
	"teststopkopach": {
		"go install -v %ldflags",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan kopach",
	},
	"teststopnode": {
		"go install -v %ldflags",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan node",
	},
	"teststopwallet": {
		"go install -v %ldflags",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan wallet",
	},
	"nodegui": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet nodegui",
	},
	"testnode": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l trace --solo --lan --norpc=false node",
	},
	"testwallet": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l trace --walletpass aoeuaoeu --solo --lan wallet",
	},
	"testkopach": {
		"go install -v %ldflags",
		"pod -D %datadir -n testnet -l trace -g -G 1 --solo --lan kopach",
	},
	"resetwallet": {
		"pod -D %datadir -l trace --walletpass aoeuaoeu wallet drophistory",
	},
	"stroy": {
		"go install -v %ldflags ./stroy/.",
	},
}
