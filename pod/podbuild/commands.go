package main

var commands = map[string][]string{
	"build": {
		"podbuild generate",
		"go build -v ./pod/.",
	},
	"install": {
		"podbuild generate",
		"go install -v ./pod/.",
	},
	"gui": {
		"podbuild generate",
		"go run -v ./pod/. gui",
	},
	"node": {
		"podbuild generate",
		"go run -v ./pod/. node",
	},
	"wallet": {
		"podbuild generate",
		"go run -v ./pod/.",
	},
	"kopach": {
		"podbuild generate",
		"go run -v ./pod/.",
	},
	"headless": {
		"podbuild generate",
		"go install -v -tags headless ./pod/.",
	},
	"docker": {
		"podbuild generate",
		"go install -v -tags headless ./pod/.",
	},
	"appstores": {
		"podbuild generate",
		"go install -v -tags nominers ./pod/.",
	},
	"tests": {
		"podbuild generate",
		"go test ./...",
	},
	"builder": {
		"go install -v ./pod/podbuild/.",
	},
	"generate":{
		"cd pkg/gel/;go generate ./...",
		"cd pkg/interrupt; go generate ./...",
		"cd pkg/log; go generate ./...",
		"cd pkg/opts; go generate ./...",
		"cd pkg/qu; go generate ./...",
		"go generate ./...",
	},
}
