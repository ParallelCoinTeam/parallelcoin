//+build ignore

package worker

import (
	"encoding/binary"
	"fmt"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/ipc"
	"github.com/p9c/pod/pkg/util"
)

func main() {
	ctrl, err := ipc.NewController()
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
	//ctrl.Out.Write(b)
	ctrl.Out.Write(b)
	ctrl.Out.Write(ipc.QuitCommand)
	err = ctrl.Wait()
}
