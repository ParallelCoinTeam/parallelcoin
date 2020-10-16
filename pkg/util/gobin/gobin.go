package gobin

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Get() (gobinary string, err error) {
	// search the environment variables for a GOROOT, if it exists we know we can run Go
	env := os.Environ()
	envMap := make(map[string]string)
	for i := range env {
		split := strings.Split(env[i], "=")
		if len(split) < 2 {
			continue
		}
		envMap[split[0]] = split[1]
	}
	goroot, ok := envMap["GOROOT"]
	if ok {
		gobinary = filepath.Join(goroot, "bin", "go")
	} else {
		err = errors.New("no GOROOT found, no Go binary available")
	}
	return
}
