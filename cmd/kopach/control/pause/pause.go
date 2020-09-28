package pause

import (
	"net"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/pkg/coding/simplebuffer"
	"github.com/p9c/pod/pkg/coding/simplebuffer/IPs"
	"github.com/p9c/pod/pkg/coding/simplebuffer/Uint16"
)

var Magic = []byte{'p', 'a', 'u', 's'}

type Container struct {
	simplebuffer.Container
}

func GetPauseContainer(cx *conte.Xt) *Container {
	mB := p2padvt.Get(cx).CreateContainer(Magic)
	return &Container{*mB}
}

func LoadPauseContainer(b []byte) (out *Container) {
	out = &Container{}
	out.Data = b
	return
}

func (j *Container) GetIPs() []*net.IP {
	return IPs.New().DecodeOne(j.Get(0)).Get()
}

func (j *Container) GetP2PListenersPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(1)).Get()
}

func (j *Container) GetRPCListenersPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(2)).Get()
}

func (j *Container) GetControllerListenerPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(3)).Get()
}
