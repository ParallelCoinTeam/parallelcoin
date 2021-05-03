package main

var commands = map[string][]string{
	"build": {
		"go build -v",
	},
	"install": {
		"go install -v",
	},
	"tests": {
		"go test ./...",
	},
	"builder": {
		"go install -v ./logbuild/.",
	},
}
