package main

var commands = map[string][]string{
	"build": {
		"go generate ./version/.",
		"go build -v ./pod/.",
	},
	"install": {
		"go generate ./version/.",
		"go install -v ./pod/.",
	},
	"gui": {
		"go generate ./version/.",
		"go run -v ./pod/. gui",
	},
	"node": {
		"go generate ./version/.",
		"go run -v ./pod/. node",
	},
	"wallet": {
		"go generate ./version/.",
		"go run -v ./pod/.",
	},
	"kopach": {
		"go generate ./version/.",
		"go run -v ./pod/.",
	},
	"headless": {
		"go generate ./version/.",
		"go install -v -tags headless ./pod/.",
	},
	"docker": {
		"go generate ./version/.",
		"go install -v -tags headless ./pod/.",
	},
	"appstores": {
		"go generate ./version/.",
		"go install -v -tags nominers ./pod/.",
	},
	"tests": {
		"go generate ./version/.",
		"go test ./...",
	},
	"builder": {
		"go generate ./version/.",
		"go install -v ./pod/podbuild/.",
	},
	"generate": {
		"cd pkg/gel/;go generate ./...",
		"cd pkg/interrupt; go generate ./...",
		"cd pkg/log; go generate ./...",
		"cd pkg/opts; go generate ./...",
		"cd pkg/qu; go generate ./...",
		"go generate ./...",
	},
}
