package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	cwd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	ioutil.WriteFile(filepath.Join(home, ".cwd"), []byte(cwd), 0600)
}
