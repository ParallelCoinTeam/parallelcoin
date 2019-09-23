//+build !nogui
// +build !headless

package db

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/vue/mod"

	scribble "github.com/nanobox-io/golang-scribble"
)

type DuoVUEdb struct {
	DB     *scribble.Driver
	Folder string      `json:"folder"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
}

type DVdb interface {
	DbReadAllTypes()
	DbRead(folder, name string)
	DbReadAll(folder string) mod.DuoGuiItems
	DbWrite(folder, name string, data interface{})
}

func (d *DuoVUEdb) DuoVueDbInit(dataDir string) {
	db, err := scribble.New(dataDir+"/gui", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	d.DB = db
}
