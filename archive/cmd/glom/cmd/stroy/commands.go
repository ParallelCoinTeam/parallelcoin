package main

var commands = map[string][]string{
	"build": {
		"go build -v %ldflags",
	},
	"install": {
		"go install -v %ldflags",
	},
	"builder": {
		"go install -v %ldflags ./cmd/stroy/.",
	},
}
