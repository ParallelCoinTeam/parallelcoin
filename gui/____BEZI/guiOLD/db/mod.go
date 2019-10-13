//+build !nogui
// +build !headless

package db

import (
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	mod2 "github.com/p9c/pod/gui/____BEZI/guiOLD/mod"
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
	DbReadAll(folder string) mod2.DuOSitems
	DbWrite(folder, name string, data interface{})
}

func (d *DuOSdb) DuoVueDbInit(dataDir string) {
	db, err := scribble.New(dataDir+"/guiOLD", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	d.DB = db
}
