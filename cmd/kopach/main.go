package kopach

import (
	"encoding/binary"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/ipc"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
)

// Main the main thread of the kopach miner
func Main(cx *conte.Xt, quit chan struct{}) {
	ctrl, err := ipc.NewController()
	if err != nil {
		log.ERROR(err)
	}
	err = ctrl.Start()
	if err != nil {
		log.ERROR(err)
	}
		hash, err := chainhash.NewHash(make([]byte, 32))
	if err != nil {
		log.ERROR(err)
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
	ctrl.In.Write(b)

	ctrl.In.Write(ipc.QuitCommand)

	ctrl.Wait()
	<-quit
	log.DEBUG("stopping kopach miner")
}
