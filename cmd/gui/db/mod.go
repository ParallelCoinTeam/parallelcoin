//+build !nogui
// +build !headless

package db

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/pkg/log"

	scribble "github.com/nanobox-io/golang-scribble"
)

type DuOSdb struct {
	DB     *scribble.Driver
	Folder string      `json:"folder"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
}

type Ddb interface {
	DbReadAllTypes()
	DbRead(folder, name string)
	DbReadAll(folder string) mod.DuOSitems
	DbWrite(folder, name string, data interface{})
}

func (d *DuOSdb) DuoVueDbInit(dataDir string) {
	db, err := scribble.New(dataDir+"/gui", nil)
	if err != nil {
		log.ERROR(err)
fmt.Println("Error", err)
	}
	d.DB = db
}
