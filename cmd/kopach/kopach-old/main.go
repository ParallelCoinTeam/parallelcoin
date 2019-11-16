package kopach_old

import (
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/ipc"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"os"
)

// Main the main thread of the kopach miner
func Main(cx *conte.Xt, quit chan struct{}) {
	args := append(os.Args[:len(os.Args)-1], "worker")
	ctrl, err := ipc.NewController(args)
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
	log.DEBUG("writing a message")
	ctrl.Write(b)
	log.DEBUG("reading")
	ctrl.Read(b)
	log.DEBUG("writing another message")
	ctrl.Write([]byte("testing"))
	log.DEBUG("reading")
	ctrl.Read(b)
	log.DEBUG("sending close signal to worker")
	ctrl.Close()
	log.DEBUG("reading")
	ctrl.Read(b)
	log.DEBUG("now what?")
	ctrl.Cmd.Process.Kill()
	log.SPEW(b)
	//ctrl.Wait()
	//<-quit
	log.DEBUG("stopping kopach miner")
}
