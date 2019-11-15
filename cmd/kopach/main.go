package kopach

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

// Main the main thread of the kopach miner
func Main(cx *conte.Xt, quit chan struct{}) {
	<-quit
	log.DEBUG("stopping kopach miner")
}
