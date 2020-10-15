package jdb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/ipfs/go-cid"
	"os"
)

func (j *JavazacDB) Write(collection string, key []byte) {
	buf := bytes.NewReader(key)
	n, err := j.peer.AddFile(j.ctx, buf, nil)
	checkError(err)
	j.index[collection] = n.Cid().String()

	var bytesBuf bytes.Buffer
	encoder := gob.NewEncoder(&bytesBuf)
	err = encoder.Encode(j.index)

	bufIndex := bytes.NewReader(bytesBuf.Bytes())
	index, err := j.peer.AddFile(j.ctx, bufIndex, nil)
	checkError(err)

	fmt.Println("cii:", n.Cid())
	fmt.Println("cii:", n.String())
	fmt.Println("cii:", n.Tree("", -1))
	fmt.Println("cii:", index.Cid())
}

func (j *JavazacDB) Read(fileName string) {
	c, _ := cid.Decode(fileName)
	rsc, err := j.peer.GetFile(j.ctx, c)
	if err != nil {
		panic(err)
	}
	//decoder := gob.NewDecoder(rsc)
	//err = decoder.Decode(key)
	fmt.Println("DroljaIZKnjazevca", rsc)
	checkError(err)
	defer rsc.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func (j *JavazacDB) ReadRoot() {

	//c, _ := cid.Decode("")
	//rsc := j.peer.DAGService.Session(j.ctx)
	//if err != nil {
	//	panic(err)
	//}
	//var key interface{}
	//decoder := gob.NewDecoder(rsc.Links())
	//err = decoder.Decode(key)
	//checkError(err)
	fmt.Println("tete ")
	//defer rsc.Close()
}
