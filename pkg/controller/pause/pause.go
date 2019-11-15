package pause

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/advertisment"
	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/simplebuffer/IPs"
	"github.com/p9c/pod/pkg/simplebuffer/Uint16"
	"net"
)

var PauseMagic = [4]byte{'p', 'a', 'u', 's'}

type PauseContainer struct {
	simplebuffer.Container
}

func GetPauseContainer(cx *conte.Xt) *PauseContainer {
	mB := advertisment.Get(cx).CreateContainer(PauseMagic)
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

func (mC *PauseContainer) GetP2PListenerPort() uint16 {
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

func (mC *PauseContainer) GetRPCListenerPort() uint16 {
	return Uint16.New().DecodeOne(mC.Get(2)).Get()
}

func (mC *PauseContainer) GetRPCListener() (out []string) {
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
