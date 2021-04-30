package main

var commands = map[string][]string{
	"build": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go build -v ./pod/.",
	},
	"install": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go install -v ./pod/.",
	},
	"gui": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go run -v ./pod/. gui",
	},
	"node": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go run -v ./pod/. node",
	},
	"wallet": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go run -v ./pod/.",
	},
	"kopach": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go run -v ./pod/.",
	},
	"headless": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go install -v -tags headless ./pod/.",
	},
	"docker": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go install -v -tags headless ./pod/.",
	},
	"appstores": {
		"go install -v ./pod/podbuild/.",
		"podbuild generate",
		"go install -v -tags nominers ./pod/.",
	},
	"tests": {
		"go install -v ./pod/podbuild/.",
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
