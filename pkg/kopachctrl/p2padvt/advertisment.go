package p2padvt

import (
	"net"

	"github.com/p9c/simplebuffer"
	"github.com/p9c/simplebuffer/IPs"
	"github.com/p9c/simplebuffer/Uint16"

	"github.com/p9c/pod/pkg/conte"
)

var Magic = []byte{'a', 'd', 'v', 't'}

type Container struct {
	simplebuffer.Container
}

// LoadContainer takes a message byte slice payload and loads it into a container
// ready to be decoded
func LoadContainer(b []byte) (out Container) {
	out.Data = b
	return
}

func Get(cx *conte.Xt) simplebuffer.Serializers {
	return simplebuffer.Serializers{
		IPs.GetListenable(),
		Uint16.GetPort((*cx.Config.Listeners)[0]),
		Uint16.GetPort((*cx.Config.RPCListeners)[0]),
		Uint16.GetPort(*cx.Config.Controller),
	}
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
