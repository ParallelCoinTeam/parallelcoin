package pause

import (
	"net"

	"github.com/p9c/simplebuffer"
	"github.com/p9c/simplebuffer/IPs"
	"github.com/p9c/simplebuffer/Uint16"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/kopachctrl/p2padvt"
)

var PauseMagic = []byte{'p', 'a', 'u', 's'}

type PauseContainer struct {
	simplebuffer.Container
}

func GetPauseContainer(cx *conte.Xt) *PauseContainer {
	mB := p2padvt.Get(cx).CreateContainer(PauseMagic)
	return &PauseContainer{*mB}
}

func LoadPauseContainer(b []byte) (out *PauseContainer) {
	out = &PauseContainer{}
	out.Data = b
	return
}

func (mC *PauseContainer) GetIPs() []*net.IP {
	return IPs.New().DecodeOne(mC.Get(0)).Get()
}

func (mC *PauseContainer) GetP2PListenersPort() uint16 {
	return Uint16.New().DecodeOne(mC.Get(1)).Get()
}

func (mC *PauseContainer) GetP2PListeners() (out []string) {
	p := Uint16.New().DecodeOne(mC.Get(1)).String()
	ips := mC.GetIPs()
	for i := range ips {
		out = append(out, net.JoinHostPort(ips[i].String(), p))
	}
	return
}

func (mC *PauseContainer) GetRPCListenersPort() uint16 {
	return Uint16.New().DecodeOne(mC.Get(2)).Get()
}

func (mC *PauseContainer) GetRPCListeners() (out []string) {
	p := Uint16.New().DecodeOne(mC.Get(2)).String()
	ips := mC.GetIPs()
	for i := range ips {
		out = append(out, net.JoinHostPort(ips[i].String(), p))
	}
	return
}
func (mC *PauseContainer) GetControllerListenerPort() uint16 {
	return Uint16.New().DecodeOne(mC.Get(3)).Get()
}

func (mC *PauseContainer) GetControllerListener() (out []string) {
	p := Uint16.New().DecodeOne(mC.Get(3)).String()
	ips := mC.GetIPs()
	for i := range ips {
		out = append(out, net.JoinHostPort(ips[i].String(), p))
	}
	return
}
