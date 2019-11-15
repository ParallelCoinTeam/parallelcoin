package main

import (
	"encoding/binary"
	"fmt"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/util"
	"io"
	"os"
	"os/exec"
)

var Q = []byte{255, 255, 255, 255}

type ControllerIPC struct {
	*exec.Cmd
	Stdin  io.Writer
	Stdout io.Reader
}

func NewControllerIPC() (out *ControllerIPC, err error) {
	out = &ControllerIPC{
		Cmd: exec.Command("go", "run", "worker.go"),
	}
	out.Stdin, err = out.StdinPipe()
	if err != nil {
		panic(err)
	}
	out.Stderr = os.Stdout
	out.Stdout, err = out.StdoutPipe()
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	ctrl, err := NewControllerIPC()
	if err != nil {
		panic(err)
	}
	err = ctrl.Start()
	if err != nil {
		fmt.Println(err)
	}
	hash, err := chainhash.NewHash(make([]byte, 32))
	if err != nil {
		fmt.Println(err)
	}
	blk := util.NewBlock(wire.NewMsgBlock(wire.NewBlockHeader(
		100,
		hash,
		hash,
		4242,
		4242,
	)))
	b, err := blk.Bytes()
	prefix := make([]byte, 4)
	binary.BigEndian.PutUint32(prefix, uint32(len(b)))
	b = append(prefix, b...)
	//ctrl.Stdout.Write(b)
	ctrl.Stdin.Write(b)
	ctrl.Stdin.Write(Q)
	err = ctrl.Wait()
}
