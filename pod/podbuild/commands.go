package main

var commands = map[string][]string{
	"build": {
		"go build -v ./pod/.",
	},
	"install": {
		"go install -v ./pod/.",
	},
	"gui": {
		"go run -v ./pod/. gui",
	},
	"node": {
		"go run -v ./pod/. node",
	},
	"wallet": {
		"go run -v ./pod/.",
	},
	"kopach": {
		"go run -v ./pod/.",
	},
	"headless": {
		"go install -v -tags headless ./pod/.",
	},
	"docker": {
		"go install -v -tags headless ./pod/.",
	},
	"appstores": {
		"go install -v -tags nominers ./pod/.",
	},
	"tests": {
		"go test ./...",
	},
	"builder": {
		"go install -v ./pod/podbuild/.",
	},
}
